package services

import (
	"credit-risk-mvp/internal/logger"
	"os"

	"github.com/IBM/sarama"
)

var producer sarama.SyncProducer

func InitProducer(brokers []string) {
	var err error
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	producer, err = sarama.NewSyncProducer(brokers, config)
	if err != nil {
		logger.Lg.Error("error creating producer Kafka", "error", err)
		os.Exit(1)
	}
}

func CloseProducer() {
	if producer != nil {
		if err := producer.Close(); err != nil {
			logger.Lg.Error("Error closing Kafka producer", "error", err)
		}
	}
}

func SendMessage(topic string, key string, message []byte) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.ByteEncoder(message),
	}
	_, _, err := producer.SendMessage(msg)
	return err
}
