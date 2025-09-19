package consumer

import (
	"context"
	"demo-service/internal/models"
	"encoding/json"
	"errors"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func (k *Kafka) Start() error {
	if !k.runing {
		k.runing = true
		go k.kafkaLoop()
		return nil
	}
	return errors.New("kafka is runing")
}

func (k *Kafka) subsribe(topic string) error {
	return k.consumer.SubscribeTopics([]string{topic}, nil)
}

func (k *Kafka) Close() error {
	select {
	case _, ex := <-k.signalCh:
		if !ex {
			return errors.New("TODO")
		}
		close(k.signalCh)
	default:
		close(k.signalCh)
	}

	<-k.statusCh
	return k.consumer.Close()
}

func (k *Kafka) kafkaLoop() {
	defer func() {
		k.statusCh <- struct{}{}
	}()

	for {
		select {
		case <-k.signalCh:
			return
		default:
			ev := k.consumer.Poll(100)
			if ev == nil {
				continue
			}
			switch e := ev.(type) {
			case *kafka.Message:
				if k.handler != nil {
					order, err := unmarshal(e.Value)
					if err != nil {
						continue
					}
					if err := k.handler(context.TODO(), order); err != nil {
						continue
					}
				}
				k.consumer.CommitMessage(e)
			case kafka.Error:
				if e.Code() == kafka.ErrAllBrokersDown {
					return
				}
			}
		}
	}
}

func unmarshal(data []byte) (*models.Order, error) {
	order := &models.Order{}
	if err := json.Unmarshal(data, order); err != nil {
		return nil, err
	}
	return order, nil
}
