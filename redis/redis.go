package redis

import (
	"encoding/json"
	"hello/db"
	"log"
	"strconv"

	"github.com/go-redis/redis"
)

type Redis struct {
	client *redis.Client
}

func NewClient() (*Redis, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	Redis := &Redis{redisClient}
	_, err := Redis.client.Ping(Redis.client.Context()).Result()
	if err != nil {
		log.Fatal(err)

		return nil, err
	}

	return Redis, nil
}

func (redis *Redis) GetUser(tgId int) User {
	redisKey := "tg_id_" + strconv.Itoa(tgId)
	user, err := redis.client.Get(redis.client.Context(), redisKey).Result()
	if err == redis.Nil {
		userDb := db.Db.GetUser(tgId)
		jsonUser, err := json.Marshal(userDb)
		if err != nil {
			panic(err)
		}
		redis.client.Set(redis.client.Context(), redisKey, jsonUser, 0).Err()
		return User{userDb.Id, userDb.TgId, userDb.Name}
	} else {
		usr := User{}
		err := json.Unmarshal([]byte(user), &usr)
		if err != nil {
			panic(err)
		}
		return usr
	}
}
