package consumer

import (
	"context"
	"demo-service/internal/models"
	"fmt"
	"sync"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type MessageHandler func(ctx context.Context, order *models.Order) error

type Config struct {
	Topic   string
	Host    string
	Port    uint16
	GroupID string
}

func (c Config) hostPort() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

type Kafka struct {
	ctx      context.Context
	cancel   context.CancelFunc
	wg       sync.WaitGroup
	runing   bool
	consumer *kafka.Consumer
	handler  MessageHandler
}

func NewKafka(cnf Config, handler MessageHandler) (*Kafka, error) {
	kafkaConf := &kafka.ConfigMap{
		"bootstrap.servers":  cnf.hostPort(),
		"group.id":           cnf.GroupID,
		"auto.offset.reset":  "earliest",
		"enable.auto.commit": false,
	}

	consumer, err := kafka.NewConsumer(kafkaConf)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	broker := Kafka{
		ctx:      ctx,
		cancel:   cancel,
		runing:   false,
		consumer: consumer,
		handler:  handler,
	}

	broker.subsribe(cnf.Topic)

	return &broker, nil
}
