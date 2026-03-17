package task

import (
	"context"
	"log"
	"sync"

	"github.com/teakingwang/ginmicro/internal/order/service"
	"github.com/teakingwang/ginmicro/pkg/mq"
)

type TaskManager struct {
	orderConsumer *OrderConsumer
}

func NewTaskManager(
	kafkaClientOrder *mq.KafkaClient,
	orderSrv service.OrderService,
) *TaskManager {
	return &TaskManager{
		orderConsumer: NewOrderConsumer(kafkaClientOrder, orderSrv),
	}
}

func (tm *TaskManager) Start(ctx context.Context) error {
	var wg sync.WaitGroup

	// 启动订单消费者
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := tm.orderConsumer.Run(ctx); err != nil {
			log.Printf("Order consumer error: %v", err)
		}
	}()

	// 等待消费者完成
	wg.Wait()
	return nil
}
