package services

import (
	"log"

	"github.com/IBM/sarama"
)

func StartConsumer(topic string) {

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer([]string{"localhost:9092"}, config)
	if err != nil {
		log.Fatalf("error connecting to Kafka: %v", err)
	}
	defer consumer.Close()

	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("error starting partition consumer: %v", err)
	}
	defer partitionConsumer.Close()
	for msg := range partitionConsumer.Messages() {
		log.Printf("new message: %s", string(msg.Value))
	}
}
