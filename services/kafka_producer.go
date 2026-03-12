package services

import (
	"log"

	"github.com/IBM/sarama"
)

var producer sarama.SyncProducer

func InitProducer(brokers []string) {
	var err error
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	producer, err = sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.Fatalf("error creating producer Kafka: %v", err)
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
