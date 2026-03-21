package domain

import "context"

type TerroristStore interface {
	IsTerrorist(ctx context.Context, passport string) (bool, error)
	UpdateList(ctx context.Context, passports []string) error
}
