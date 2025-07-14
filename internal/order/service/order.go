package service

import (
	"context"
	"github.com/teakingwang/ginmicro/api/user"
	"github.com/teakingwang/ginmicro/internal/order/repository"
	"github.com/teakingwang/ginmicro/pkg/datastore"
	"github.com/teakingwang/ginmicro/pkg/logger"
	kafka "github.com/teakingwang/ginmicro/pkg/mq"
)

type OrderDTO struct {
	OrderID  int64
	OrderSN  string
	UserID   int64
	Username string
}

type OrderService interface {
	GetOrder(ctx context.Context, id int64) (*OrderDTO, error)
}

type orderService struct {
	orderRepo   repository.OrderRepo
	redis       datastore.Store
	userClient  user.UserServiceClient
	kafkaClient *kafka.KafkaClient
}

func NewOrderService(orderRepo repository.OrderRepo, redis datastore.Store, kafkaClient *kafka.KafkaClient, userClient user.UserServiceClient) OrderService {
	return &orderService{orderRepo: orderRepo, redis: redis, kafkaClient: kafkaClient, userClient: userClient}
}

func (s *orderService) GetOrder(ctx context.Context, id int64) (*OrderDTO, error) {
	logger.Info("GetOrder called with ID:", id)
	or, err := s.orderRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if or == nil {
		return nil, nil
	}

	// 调用user
	logger.Info("Fetching user for order userID:", or.UserID)
	userRes, err := s.userClient.GetUser(ctx, &user.GetUserRequest{Id: or.UserID})
	if err != nil {
		return nil, err
	}

	dto := &OrderDTO{
		OrderID: or.OrderID,
		OrderSN: or.OrderSN,
	}

	if userRes != nil {
		dto.UserID = userRes.Id
		dto.Username = userRes.Username
	}

	// 演示kafka发消息
	err = s.kafkaClient.Produce(ctx, []byte("order_event"), []byte("Order retrieved: "+dto.OrderSN))
	if err != nil {
		logger.Error("Failed to produce Kafka message:", err)
	}

	return dto, nil
}
