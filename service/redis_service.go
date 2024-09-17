package service

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
)

type RedisServiceInterface interface {
	Set(key string, value interface{}) error
	Get(key string) (string, error)
	Delete(key string) error
}

type RedisService struct {
	client *redis.Client
}

func NewRedisService(client *redis.Client) RedisService {
	return RedisService{client: client}
}

func (rs *RedisService) Set(key string, value interface{}) error {
	p, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return rs.client.Set(context.Background(), key, p, 0).Err()
}

func (rs *RedisService) Get(key string) (string, error) {
	p, err := rs.client.Get(context.Background(), key).Result()
	if err != nil {
		return "", err
	}

	return p, err
}

func (rs *RedisService) Delete(key string) error {
	return rs.client.Del(context.Background(), key).Err()
}
