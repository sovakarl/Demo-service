// main.go
package main

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"producer/models"
	"strconv"
	"time"

	"github.com/segmentio/kafka-go"
)

// Генерация случайного заказа на основе ваших моделей
func generateRandomOrder(uid string) models.Order {
	now := time.Now().UTC()
	trackNumber := "WBILMTESTTRACK-" + randomString(8)

	return models.Order{
		OrderUID:          uid,
		TrackNumber:       trackNumber,
		Entry:             "WBIL",
		Locale:            "en",
		InternalSignature: "",
		CustomerID:        "test",
		DeliveryService:   "meest",
		ShardKey:          "9",
		SmID:              99,
		DateCreated:       now,
		OffShard:          "1",
		Delivery: models.Delivery{
			Name:    "Test Testov " + randomString(5),
			Phone:   "+972" + randomDigits(7),
			Zip:     randomDigits(7),
			City:    "Kiryat Mozkin",
			Address: "Ploshad Mira " + randomDigits(2),
			Region:  "Kraiot",
			Email:   "test" + randomDigits(3) + "@gmail.com",
		},
		Payment: models.Payment{
			Transaction:  uid,
			RequestID:    "",
			Currency:     "USD",
			Provider:     "wbpay",
			Amount:       rand.Intn(5000) + 1000,
			PaymentDT:    int(now.Unix()), // ← int, как в вашей модели
			Bank:         "alpha",
			DeliveryCost: 1500,
			GoodsTotal:   rand.Intn(1000) + 100,
			CustomFee:    0,
		},
		Items: []models.Item{
			{
				ChrtID:      rand.Int63n(10000000) + 1000000,
				TrackNumber: trackNumber,
				Price:       rand.Intn(800) + 200,
				Rid:         "ab4219087a764ae0b" + randomString(4),
				Name:        "Mascaras " + randomString(4),
				Sale:        rand.Intn(50),
				Size:        "0",
				TotalPrice:  rand.Intn(500) + 100,
				NmID:        rand.Intn(5000000) + 1000000,
				Brand:       "Vivienne Sabo",
				Status:      202,
			},
		},
	}
}

func randomString(n int) string {
	letters := "abcdefghijklmnopqrstuvwxyz"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func randomDigits(n int) string {
	digits := "0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = digits[rand.Intn(len(digits))]
	}
	return string(b)
}

func main() {
	// rand.Seed(time.Now().UnixNano())

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "order",
	})

	defer writer.Close()
	log.Println("🚀 Продюсер запущен. Отправка в топик 'order'...")

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	UID := 1
	for range ticker.C {
		order := generateRandomOrder(strconv.Itoa(UID))

		// Сериализуем в JSON
		jsonData, err := json.Marshal(order)
		if err != nil {
			log.Printf("❌ Ошибка marshal: %v", err)
			continue
		}

		// Отправляем в Kafka
		err = writer.WriteMessages(context.Background(), kafka.Message{
			Key:   []byte(order.OrderUID), // опционально: ключ = order_uid
			Value: jsonData,
		})
		if err != nil {
			log.Printf("❌ Ошибка Kafka: %v", err)
		} else {
			log.Printf("✅ Отправлен заказ: %s", order.OrderUID)
			UID++
		}

	}
}
