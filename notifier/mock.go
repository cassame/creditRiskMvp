package notifier

import (
	"credit-risk-mvp/internal/domain"

	"github.com/stretchr/testify/mock"
)

type MockNotifier struct {
	mock.Mock
}

func (m *MockNotifier) Notify(app domain.Application, status string) error {
	args := m.Called(app, status)
	return args.Error(0)
}
