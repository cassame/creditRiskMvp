package storage

import (
	"context"
	"credit-risk-mvp/internal/domain"

	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) GetByID(ctx context.Context, id string) (domain.Application, error) {
	args := m.Called(ctx, id)
	app := args.Get(0).(domain.Application)
	return app, args.Error(1)
}

func (m *MockRepository) SaveApplication(ctx context.Context, app domain.Application) error {
	args := m.Called(ctx, app)
	return args.Error(0)
}
