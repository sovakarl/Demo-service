package consumer

import (
	"context"
	"demo-service/internal/models"
	"fmt"

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
	statusCh chan struct{}
	signalCh chan struct{}
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

	broker := Kafka{
		statusCh: make(chan struct{}, 1),
		signalCh: make(chan struct{}),
		runing:   false,
		consumer: consumer,
		handler:  handler,
	}

	broker.subsribe(cnf.Topic)

	return &broker, nil
}
