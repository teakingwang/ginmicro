package mq

import (
	"context"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaClient struct {
	Brokers []string
	Topic   string
	GroupID string
}

func NewKafkaClient(brokers []string, topic, groupID string) *KafkaClient {
	return &KafkaClient{
		Brokers: brokers,
		Topic:   topic,
		GroupID: groupID,
	}
}

// 发送消息
func (k *KafkaClient) Produce(ctx context.Context, key, value []byte) error {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  k.Brokers,
		Topic:    k.Topic,
		Balancer: &kafka.Hash{},
	})
	defer writer.Close()

	msg := kafka.Message{
		Key:   key,
		Value: value,
		Time:  time.Now(),
	}

	return writer.WriteMessages(ctx, msg)
}

// 消费消息，handler是业务处理函数
func (k *KafkaClient) Consume(ctx context.Context, handler func(msg kafka.Message) error) error {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  k.Brokers,
		GroupID:  k.GroupID,
		Topic:    k.Topic,
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})
	defer reader.Close()

	for {
		m, err := reader.ReadMessage(ctx)
		if err != nil {
			if err == context.Canceled {
				return nil
			}
			log.Printf("kafka read error: %v", err)
			continue
		}

		if err := handler(m); err != nil {
			log.Printf("kafka handle message error: %v", err)
		}
	}
}
