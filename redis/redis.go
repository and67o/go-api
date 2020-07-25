package redis

import (
	"github.com/go-redis/redis"
	"go-api/config"
	"log"
)

type Operations interface {
	OrdersOperations
	UserOperations
	BaseOperations
}

type Manager struct {
	Client *redis.Client
}

var Redis Operations

func init() {
	conf := config.New()
	redisClient := redis.NewClient(&redis.Options{
		Addr:     conf.Redis.Addr,
		Password: conf.Redis.Pass,
		DB:       conf.Redis.Db,
	})

	_, err := redisClient.Ping(redisClient.Context()).Result()
	if err != nil {
		log.Fatal(err)
	}

	log.Print("Redis connect")
	Redis = &Manager{redisClient}
}
