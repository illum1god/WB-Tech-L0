package kafka

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"WB-L0/internal/configs"
	"WB-L0/internal/usecase"

	"github.com/IBM/sarama"
)

type Consumer interface {
	Consume(wg *sync.WaitGroup, ctx context.Context) error
	Stop()
}

type consumer struct {
	ready   chan bool
	client  sarama.ConsumerGroup
	service usecase.Service
	config  configs.KafkaConfig
	groupID string
	topic   string
}

func NewConsumer(cfg configs.KafkaConfig, service usecase.Service) Consumer {

	saramaConfig := sarama.NewConfig()
	saramaConfig.Version = sarama.V2_8_0_0
	saramaConfig.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	saramaConfig.Consumer.Offsets.Initial = sarama.OffsetNewest
	saramaConfig.Consumer.Return.Errors = true

	client, err := sarama.NewConsumerGroup(cfg.Brokers, cfg.GroupID, saramaConfig)
	if err != nil {
		log.Fatalf("Ошибка создания клиента группы потребителей: %v", err)
	}

	return &consumer{
		ready:   make(chan bool),
		client:  client,
		service: service,
		config:  cfg,
		groupID: cfg.GroupID,
		topic:   cfg.Topic,
	}
}

func (c *consumer) Consume(wg *sync.WaitGroup, ctx context.Context) error {
	handler := c

	for {
		if err := c.client.Consume(ctx, []string{c.topic}, handler); err != nil {
			log.Printf("Ошибка из потребителя: %v", err)
			time.Sleep(2 * time.Second)
			continue
		}

		if ctx.Err() != nil {
			return ctx.Err()
		}

		c.ready = make(chan bool)
	}
}

func (c *consumer) Stop() {
	if err := c.client.Close(); err != nil {
		log.Printf("Ошибка при закрытии клиента: %v", err)
	}
}

func (c *consumer) Setup(sarama.ConsumerGroupSession) error {
	close(c.ready)
	return nil
}

func (c *consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (c *consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		log.Printf("Сообщение получено: значение = %s, временная метка = %v, топик = %s",
			string(message.Value), message.Timestamp, message.Topic)

		var order usecase.OrderInput
		if err := json.Unmarshal(message.Value, &order); err != nil {
			log.Printf("Ошибка маршалинга сообщения: %v", err)
			continue
		}

		if err := c.service.SaveOrder(context.Background(), order.Order); err != nil {
			log.Printf("Ошибка сохранения заказа: %v", err)
			continue
		}

		session.MarkMessage(message, "")
	}

	return nil
}
