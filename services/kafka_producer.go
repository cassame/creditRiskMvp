package services

import (
	"log"

	"github.com/IBM/sarama"
)

var producer sarama.SyncProducer

func InitProducer() {
	var err error
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	producer, err = sarama.NewSyncProducer([]string{"localhost:9092"}, config)
	if err != nil {
		log.Fatalf("error creating producer Kafka: %v", err)
	}
}

func SendMessage(topic string, message []byte) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(message),
	}
	_, _, err := producer.SendMessage(msg)
	return err
}
