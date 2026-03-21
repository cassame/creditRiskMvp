package services

import (
	"context"
	"credit-risk-mvp/internal/logger"

	"github.com/IBM/sarama"
)

func StartConsumer(ctx context.Context, brokers []string, topic string) {

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		logger.Lg.Error("error connecting to Kafka", "error", err)
		return
	}
	defer func() { _ = consumer.Close() }()
	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		logger.Lg.Error("error starting partition consumer", "error", err)
		return
	}
	defer func() { _ = partitionConsumer.Close() }()

	logger.Lg.Info("Kafka consumer started", "topic", topic)

	for {
		select {
		case <-ctx.Done():
			logger.Lg.Info("Stopping Kafka consumer by signal...")
			return
		case msg := <-partitionConsumer.Messages():
			logger.Lg.Info("received new message",
				"topic", msg.Topic,
				"partition", msg.Partition,
				"offset", msg.Offset,
				"content", string(msg.Value),
			)
		case err := <-partitionConsumer.Errors():
			logger.Lg.Error("Kafka consume error", "error", err)
		}
	}
}
