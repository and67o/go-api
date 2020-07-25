package redis

import (
	"encoding/json"
	"go-api/db"
)

type OrdersOperations interface {
	GetOrders() (res []db.Order, err error)
}

func (Redis *Manager) GetOrders() (res []db.Order, err error) {
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
