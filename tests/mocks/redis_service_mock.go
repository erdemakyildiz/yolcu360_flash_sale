package mocks

import (
	"github.com/stretchr/testify/mock"
)

type RedisService struct {
	mock.Mock
}

func (rs *RedisService) Set(key string, value interface{}) error {
	return nil
}

func (rs *RedisService) Get(key string) (string, error) {
	args := rs.Called(key)
	var err error
	if args.Get(1) != nil {
		err = args.Get(1).(error)
		return "", err
	}
	return args.Get(0).(string), err
}

func (rs *RedisService) Delete(key string) error {
	return nil
}
