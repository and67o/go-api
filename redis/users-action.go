package redis

import (
	"encoding/json"
	"go-api/db"
)

type UserOperations interface {
	GetUsers() (res []db.User, err error)
	DeleteUsers()
}

func (Redis *Manager) GetUsers() (res []db.User, err error) {
	redisUsers, err := Redis.Get(USERS)

	err = json.Unmarshal([]byte(redisUsers), &res)
	if err != nil {
		return nil, err
	}
	return res, err
}

func (Redis *Manager) DeleteUsers() {
	Redis.Client.Del(Redis.Client.Context(), USERS)
}
