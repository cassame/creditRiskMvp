package domain

import "context"

type Repository interface {
	SaveApplication(ctx context.Context, app Application) error
	GetByID(ctx context.Context, id string) (Application, error)
}
