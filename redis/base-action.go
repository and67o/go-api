package redis

import (
	"encoding/json"
	"github.com/go-redis/redis"
)

type BaseOperations interface {
	Set(key string, dest interface{}) (res interface{}, err error)
	Get(key string) (res string, err error)
}

func (Redis *Manager) Get(key string) (res string, err error) {
	redisData, err := Redis.Client.Get(Redis.Client.Context(), key).Result()

	if err == redis.Nil {
		return "", nil
	} else if err != nil {
		return "", err
	}
	return redisData, err
}

func (Redis *Manager) Set(key string, dest interface{}) (res interface{}, err error) {
	jsonData, err := json.Marshal(dest)
	if err != nil {
		return nil, err
	}
	Redis.Client.Set(Redis.Client.Context(), key, jsonData, 0)
	return dest, err
}


