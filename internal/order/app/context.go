package app

import (
	"github.com/teakingwang/ginmicro/api/user"
	"github.com/teakingwang/ginmicro/config"
	"github.com/teakingwang/ginmicro/internal/order/model"
	"github.com/teakingwang/ginmicro/internal/order/repository"
	"github.com/teakingwang/ginmicro/internal/order/service"
	"github.com/teakingwang/ginmicro/pkg/consul"
	"github.com/teakingwang/ginmicro/pkg/datastore"
	"github.com/teakingwang/ginmicro/pkg/db"
	"github.com/teakingwang/ginmicro/pkg/logger"
	kafka "github.com/teakingwang/ginmicro/pkg/mq"
	"github.com/teakingwang/ginmicro/pkg/utils/uuid"
	"google.golang.org/grpc"
	"gorm.io/gorm"
	"time"
)

type AppContext struct {
	DB           *gorm.DB
	Redis        datastore.Store
	OrderService service.OrderService
	UserClient   user.UserServiceClient
	UserConn     *grpc.ClientConn
	KafkaClient  *kafka.KafkaClient
}

func NewAppContext() (*AppContext, error) {
	// 初始化 DB
	gormDB, err := db.NewDB()
	if err != nil {
		logger.Error("Failed to initialize database:", err)
		return nil, err
	}

	// 初始化 Redis
	redisStore := datastore.NewRedisClient()
	if err != nil {
		logger.Error("Failed to initialize Redis:", err)
		return nil, err
	}

	// 初始化 KafkaClient
	kafkaClient := kafka.NewKafkaClient(config.Config.Kafka.Brokers, config.Config.Kafka.Topic, config.Config.Kafka.GroupID)

	// 初始化 OrderService
	orderRepo := repository.NewOrderRepo(gormDB)
	if err := orderRepo.Migrate(); err != nil {
		logger.Error("Failed to migrate order repository:", err)
		return nil, err
	}

	// ✅ 插入初始订单记录
	if err := seedInitialOrder(gormDB); err != nil {
		logger.Warn("Failed to seed initial order:", err)
	}

	// 连接 user 服务 gRPC
	dClient, err := consul.NewConsulDiscovery("consul:8500")
	if err != nil {
		logger.Error("Failed to connect to Consul:", err)
		return nil, err
	}

	conn, err := dClient.GetGRPCConn("user")
	if err != nil {
		logger.Error("Failed to connect to user service:", err)
		return nil, err
	}

	userClient := user.NewUserServiceClient(conn)

	orderSrv := service.NewOrderService(orderRepo, redisStore, kafkaClient, userClient)

	return &AppContext{
		DB:           gormDB,
		Redis:        redisStore,
		OrderService: orderSrv,
		UserClient:   userClient,
		UserConn:     conn,
		KafkaClient:  kafkaClient,
	}, nil
}

func seedInitialOrder(db *gorm.DB) error {
	// 检查是否已存在记录，避免重复插入
	var count int64
	db.Model(&model.Order{}).Count(&count)
	if count > 0 {
		return nil
	}

	now := time.Now()
	order := &model.Order{
		OrderID:     100000000000000001,
		OrderSN:     uuid.NewUUID(),
		UserID:      100000000000000001,
		Status:      model.OrderStatusPaid,
		TotalAmount: 9999,
		PayAmount:   9999,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	return db.Create(order).Error
}
