package main

import (
	"context"
	"fmt"
	"github.com/teakingwang/ginmicro/config"
	"github.com/teakingwang/ginmicro/internal/order/app"
	"github.com/teakingwang/ginmicro/internal/order/controller"
	"github.com/teakingwang/ginmicro/internal/task"
	"github.com/teakingwang/ginmicro/pkg/consul"
	"github.com/teakingwang/ginmicro/pkg/logger"
	"github.com/teakingwang/ginmicro/pkg/mq"
	"github.com/teakingwang/ginmicro/pkg/utils/idgen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"strconv"
	"syscall"
	"time"

	"github.com/teakingwang/ginmicro/api/order"
)

func main() {
	if err := logger.Init(true); err != nil {
		panic("logger init failed: " + err.Error())
	}

	defer func() {
		if r := recover(); r != nil {
			logger.Errorf("panic occurred: %v", r)
			logger.Errorf("stack trace:\n%s", string(debug.Stack()))
		}
		logger.Sync()
	}()

	if err := config.LoadConfigFromConsul("config/order"); err != nil {
		logger.Warn("load from consul failed: %v, falling back to local config", err)
		if err := config.LoadConfig(); err != nil {
			panic(fmt.Errorf("failed to load config: %v", err))
		}
	}

	// 初始化 ID 生成器
	if err := idgen.Init(); err != nil {
		panic(fmt.Errorf("failed to initialize idgen: %v", err))
	}

	if err := run(); err != nil {
		logger.Errorf("service exited with error: %v", err)
	}
}

func run() error {

	// 注入依赖
	ctx, err := app.NewAppContext()
	if err != nil {
		panic(fmt.Errorf("new appcontext err:%v", err))
	}

	// 启动 HTTP 服务
	go func() {
		httpAddr := ":" + config.Config.Server.Order.HTTPPort
		router := controller.NewHTTPRouter(ctx.OrderService)
		logger.Infof("HTTP server listening on %s", httpAddr)
		if err := http.ListenAndServe(httpAddr, router); err != nil {
			logger.Errorf("HTTP server error: %v", err)
		}
	}()

	// 启动 gRPC 服务
	go func() {
		grpcAddr := "0.0.0.0:" + config.Config.Server.Order.GRPCPort
		lis, err := net.Listen("tcp", grpcAddr)
		if err != nil {
			logger.Errorf("failed to listen grpc: %v", err)
			return
		}

		s := grpc.NewServer()
		registerHealthCheck(s)
		order.RegisterOrderServiceServer(s, controller.NewOrderController(ctx.OrderService))

		// 注册到 Consul
		consulClient, err := consul.NewConsulClient(config.Config.Consul.Address)
		if err != nil {
			logger.Errorf("failed to create consul client: %v", err)
			return
		}
		serviceID := config.Config.Server.Order.Name + "-" + config.Config.Server.Order.GRPCPort
		serviceName := config.Config.Server.Order.Name
		serviceAddress := config.Config.Server.Order.Host
		servicePort, err := strconv.Atoi(config.Config.Server.Order.GRPCPort)
		if err != nil {
			logger.Errorf("invalid service port: %v", err)
			return
		}

		logger.Infof("Registering order service to Consul: %s", serviceID)
		if err := consulClient.RegisterService(serviceID, serviceName, serviceAddress, servicePort, []string{"grpc", "order"}); err != nil {
			logger.Errorf("consul register error: %v", err)
			return
		}
		defer consulClient.DeregisterService(serviceID)

		logger.Infof("gRPC server listening on %s", grpcAddr)
		if err := s.Serve(lis); err != nil {
			logger.Errorf("gRPC server failed: %v", err)
		}
	}()

	// 启动 Kafka 消费者
	go func() {
		kafkaClientOrder := mq.NewKafkaClient(
			config.Config.Kafka.Brokers,
			config.Config.Kafka.Topic,
			config.Config.Kafka.GroupID,
		)

		// 初始化任务管理器
		taskManager := task.NewTaskManager(
			kafkaClientOrder,
			ctx.OrderService, // 假设 appContext 中包含 OrderService
		)

		// 启动消费者
		consumerCtx, cancel := context.WithCancel(context.Background())
		defer cancel()

		logger.Info("Starting Kafka consumers...")
		if err := taskManager.Start(consumerCtx); err != nil {
			logger.Errorf("Task manager error: %v", err)
		}
	}()

	// 等待退出信号
	waitForShutdown()

	return nil
}

func waitForShutdown() {
	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, syscall.SIGINT, syscall.SIGTERM)
	<-stopCh
	logger.Info("Shutdown signal received, cleaning up...")
	logger.Sync()
	time.Sleep(time.Second)
	os.Exit(0)
}

func registerHealthCheck(s *grpc.Server) {
	hs := health.NewServer()
	// 设置空服务名和 order 服务名的健康状态
	hs.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)
	hs.SetServingStatus("order", grpc_health_v1.HealthCheckResponse_SERVING)
	grpc_health_v1.RegisterHealthServer(s, hs)
}
