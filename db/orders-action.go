package db

import (
	"errors"
)

type OrderAction interface {
	GetOrder(orderId int) (order Order)
	GetOrders() (orders []Order, err error)
	DeleteOrder(orderId int) (res bool, err error)
}


func (DBM *Manager) GetOrders() (orders []Order, err error) {
	orders = []Order{}
	err = DBM.db.Select(&orders, "SELECT * FROM orders")
	if len(orders) == 0 {
		err = errors.New("no orders")
	}
	return
}

func (DBM *Manager) DeleteOrder(orderId int) (res bool, err error) {
	_, err = DBM.db.Exec("DELETE FROM orders WHERE order_id = ?", orderId)
	if err != nil {
		panic(err)
	}
	return
}

func (DBM *Manager) GetOrder(orderId int) (order Order) {
	order = Order{}
	DBM.db.Get(&order, "SELECT * FROM users WHERE tg_id = ?", orderId)
	return
}