package consumer

import (
	"context"
	"demo-service/internal/models"
	"encoding/json"
	"errors"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func (k *Kafka) Start() error {
	if k.runing {
		return errors.New("kafka is runing")
	}
	defer func() {
		k.runing = true
	}()
	k.wg.Add(1)
	go k.kafkaStart(k.ctx)
	return nil
}

func (k *Kafka) kafkaStart(ctx context.Context) {
	defer k.wg.Done()
	for {
		select {
		case <-k.ctx.Done():
			return
		default:
			k.kafkaWork()
		}
	}
}

func (k *Kafka) kafkaWork() error {
	ev := k.consumer.Poll(100)
	if ev == nil {
		return nil
	}
	switch e := ev.(type) {
	case *kafka.Message:
		if k.handler != nil {
			order, err := unmarshal(e.Value)
			if err != nil {
				return err
			}
			if err := k.handler(context.TODO(), order); err != nil {
				return err
			}
		}
		k.consumer.CommitMessage(e)
	case kafka.Error:
		if e.Code() == kafka.ErrAllBrokersDown {
			return errors.New(kafka.ErrAllBrokersDown.String())
		}
	}
	return nil
}

func (k *Kafka) subsribe(topic string) error {
	return k.consumer.SubscribeTopics([]string{topic}, nil)
}

func (k *Kafka) cancelKafkaContext() {
	defer k.cancel()
}

func (k *Kafka) Close() error {
	defer func() {
		k.runing = false
	}()

	if !k.runing {
		return errors.New("kafka is closed")
	}
	k.cancelKafkaContext()
	k.wg.Wait()
	return k.consumer.Close()
}

func unmarshal(data []byte) (*models.Order, error) {
	order := &models.Order{}
	if err := json.Unmarshal(data, order); err != nil {
		return nil, err
	}
	return order, nil
}
