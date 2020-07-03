package redis

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"go-api/db"
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

func (Redis *Redis) Get(key string) (res string, err error) {
	redisData, err := Redis.Client.Get(Redis.Client.Context(), key).Result()

	if err == redis.Nil {
		return "", nil
	} else if err != nil {
		return "", err
	}
	return redisData, err
}

func (Redis *Redis) Set(key string, dest interface{}) (res interface{}, err error) {
	jsonData, err := json.Marshal(dest)
	if err != nil {
		return nil, err
	}
	Redis.Client.Set(Redis.Client.Context(), key, jsonData, 0)
	return dest, err
}

func (Redis *Redis) GetOrders() (res []db.Order, err error) {
	redisOrders, err := Redis.Get(ORDERS)
	if redisOrders == "" {
		return []db.Order{}, nil
	}

	err = json.Unmarshal([]byte(redisOrders), &res)
	if err != nil {
		return nil, err
	}

	return res, err
}

func (Redis *Redis) GetUsers() (res []db.User, err error) {
	redisUsers, err := Redis.Get(USERS)

	err = json.Unmarshal([]byte(redisUsers), &res)
	if err != nil {
		return nil, err
	}
	return res, err
}
