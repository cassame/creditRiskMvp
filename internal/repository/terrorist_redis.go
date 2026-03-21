package repository

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisTerroristStore struct {
	client *redis.Client
	key    string
}

func NewRedisTerroristStore(client *redis.Client) *RedisTerroristStore {
	return &RedisTerroristStore{
		client: client,
		key:    "terrorist_list",
	}
}

func (r *RedisTerroristStore) IsTerrorist(ctx context.Context, passport string) (bool, error) {
	return r.client.SIsMember(ctx, r.key, passport).Result()
}

func (r *RedisTerroristStore) UpdateList(ctx context.Context, passports []string) error {
	tempKey := r.key + "_temp"

	r.client.Del(ctx, tempKey)
	if len(passports) > 0 {
		err := r.client.SAdd(ctx, tempKey, passports).Err()
		if err != nil {
			return err
		}
	}
	return r.client.Rename(ctx, tempKey, r.key).Err()
}
