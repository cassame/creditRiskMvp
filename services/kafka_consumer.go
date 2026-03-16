package services

import (
	"credit-risk-mvp/internal/logger"
	"os"

	"github.com/IBM/sarama"
)

func StartConsumer(brokers []string, topic string) {

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		logger.Lg.Error("error connecting to Kafka", "error", err)
		os.Exit(1)
	}
	defer consumer.Close()

	logger.Lg.Info("Kafka consumer started", "topic", topic)

	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		logger.Lg.Error("error starting partition consumer", "error", err)
		os.Exit(1)
	}
	defer partitionConsumer.Close()
	for msg := range partitionConsumer.Messages() {
		logger.Lg.Info("received new message",
			"topic", msg.Topic,
			"partition", msg.Partition,
			"offset", msg.Offset,
			"content", string(msg.Value),
		)
	}
}
