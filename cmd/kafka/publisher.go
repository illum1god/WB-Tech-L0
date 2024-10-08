package main

import (
	"WB-L0/internal/domain"
	"WB-L0/internal/usecase"
	"encoding/json"
	"log"
	"os"
	"strings"
	"time"

	"github.com/IBM/sarama"
)

func main() {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	brokersEnv := os.Getenv("KAFKA_BROKERS")
	if brokersEnv == "" {
		brokersEnv = "kafka:9092"
	}
	brokers := strings.Split(brokersEnv, ",")

	var producer sarama.SyncProducer
	var err error

	for i := 0; i < 10; i++ {
		producer, err = sarama.NewSyncProducer(brokers, config)
		if err == nil {
			break
		}
		log.Printf("Не удалось подключиться к Kafka, попытка %d/5: %v", i+1, err)
		time.Sleep(5 * time.Second)
	}

	if err != nil {
		log.Fatalf("Ошибка создания продюсера после 5 попыток: %v", err)
	}
	defer func(producer sarama.SyncProducer) {
		_ = producer.Close()
	}(producer)

	order := usecase.OrderInput{
		Order: domain.Order{
			OrderUID:          "b563feb7b2b84b6new",
			TrackNumber:       "WBILMTESTTRACK",
			Entry:             "WBIL",
			Delivery:          domain.Delivery{},
			Payment:           domain.Payment{},
			Items:             []domain.Item{},
			Locale:            "en",
			InternalSignature: "",
			CustomerID:        "test",
			DeliveryService:   "meest",
			ShardKey:          "9",
			SmID:              99,
			DateCreated:       time.Now(),
			OofShard:          "1",
		},
	}

	data, err := json.Marshal(order)
	if err != nil {
		log.Fatalf("Ошибка маршалинга: %v", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: "orders",
		Value: sarama.ByteEncoder(data),
	}

	time.Sleep(20 * time.Second)
	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		log.Fatalf("Ошибка отправки сообщения: %v", err)
	}

	log.Printf("Сообщение отправлено в раздел %d с оффсетом %d\n", partition, offset)

	time.Sleep(4 * time.Second)
	order2 := usecase.OrderInput{
		Order: domain.Order{
			OrderUID:          "b563feb7b2b84b6new2",
			TrackNumber:       "WBILMTESTTRACK",
			Entry:             "WBIL",
			Delivery:          domain.Delivery{},
			Payment:           domain.Payment{},
			Items:             []domain.Item{},
			Locale:            "en",
			InternalSignature: "",
			CustomerID:        "test",
			DeliveryService:   "meest",
			ShardKey:          "9",
			SmID:              99,
			DateCreated:       time.Now(),
			OofShard:          "1",
		},
	}

	data2, err := json.Marshal(order2)
	if err != nil {
		log.Fatalf("Ошибка маршалинга: %v", err)
	}

	msg2 := &sarama.ProducerMessage{
		Topic: "orders",
		Value: sarama.ByteEncoder(data2),
	}

	partition2, offset2, err := producer.SendMessage(msg2)
	if err != nil {
		log.Fatalf("Ошибка отправки сообщения: %v", err)
	}

	log.Printf("Сообщение отправлено в раздел %d с оффсетом %d\n", partition2, offset2)
}
