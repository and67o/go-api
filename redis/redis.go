package redis

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"log"
)

type Redis struct {
	Client *redis.Client
}

func NewClient() (*Redis, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	Redis := &Redis{redisClient}
	_, err := Redis.Client.Ping(Redis.Client.Context()).Result()
	if err != nil {
		log.Fatal(err)

		return nil, err
	}
	log.Print("Redis connect")
	return Redis, nil
}

func (Redis *Redis) Get(key string, dest interface{}) (res interface{}, err error)  {
	redisData, err := Redis.Client.Get(Redis.Client.Context(), key).Result()
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal([]byte(redisData), &dest)
	if err != nil {
		panic(err)
	}
	return dest, err
}

func (Redis *Redis) Set(key string, dest interface{}) (res interface{}, err error)  {
	jsonData, err := json.Marshal(dest)
	if err != nil {
		return nil, err
	}
	Redis.Client.Set(Redis.Client.Context(), key, jsonData, 0)
	return dest, err
}