package persistence

import (
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisPersistence struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisPersistence(addr, password string, db int) *RedisPersistence {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	return &RedisPersistence{
		client: client,
		ctx:    context.Background(),
	}
}

func (rp *RedisPersistence) GetCount(key string) (int, error) {
	value, err := rp.client.Get(rp.ctx, key).Result()
	if err == redis.Nil {
		return 0, nil
	} else if err != nil {
		return 0, err
	}

	count, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (rp *RedisPersistence) Increment(key string) error {
	result, err := rp.client.Incr(rp.ctx, key).Result()
	if err != nil {
		return err
	}

	if result == 1 {
		_, err = rp.client.Expire(rp.ctx, key, time.Minute).Result()
		if err != nil {
			return err
		}
	}

	return nil
}

func (rp *RedisPersistence) Block(key string, duration time.Duration) error {
	return rp.client.Set(rp.ctx, key+":blocked", "1", duration).Err()
}

func (rp *RedisPersistence) IsBlocked(key string) bool {
	_, err := rp.client.Get(rp.ctx, key+":blocked").Result()
	return err == nil
}
