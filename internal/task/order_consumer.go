package task

import (
	"context"
	"encoding/json"
	"github.com/teakingwang/ginmicro/pkg/mq"
	"log"

	"github.com/segmentio/kafka-go"
	"github.com/teakingwang/ginmicro/internal/order/service"
)

type OrderEvent struct {
	OrderID string `json:"orderID"`
	Status  string `json:"status"`
}

type OrderConsumer struct {
	KafkaClient *mq.KafkaClient
	OrderSrv    service.OrderService
}

func NewOrderConsumer(kafkaClient *mq.KafkaClient, orderSrv service.OrderService) *OrderConsumer {
	return &OrderConsumer{
		KafkaClient: kafkaClient,
		OrderSrv:    orderSrv,
	}
}

func (c *OrderConsumer) Run(ctx context.Context) error {
	return c.KafkaClient.Consume(ctx, c.handleMessage)
}

func (c *OrderConsumer) handleMessage(msg kafka.Message) error {
	log.Printf("Kafka message: topic=%s partition=%d offset=%d key=%s value=%s",
		msg.Topic, msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))

	var event OrderEvent
	if err := json.Unmarshal(msg.Value, &event); err != nil {
		log.Printf("invalid message format: %v", err)
		return err
	}

	// TODO: 处理业务逻辑
	log.Printf("Process order event: %+v", event)

	// 这里可以调用 c.OrderSvc 相关方法更新状态等

	return nil
}
