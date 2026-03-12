package services

import "github.com/stretchr/testify/mock"

type MessageQueue interface {
	SendMessage(topic string, key string, value []byte) error
}

type KafkaProducer struct{}

func (k KafkaProducer) SendMessage(topic string, key string, value []byte) error {
	return SendMessage(topic, key, value)
}

type MockQueue struct {
	mock.Mock
}

func (m *MockQueue) SendMessage(topic string, key string, value []byte) error {
	args := m.Called(topic, key, value)
	return args.Error(0)
}
