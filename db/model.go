package db

import "database/sql"

type User struct {
	Id   int    `db:"id" json:"id"`
	TgId int    `db:"tg_id" json:"tg_id"`
	Name string `db:"name" json:"name"`
	TimeStamps
}

type Order struct {
	Id     int `db:"id" json:"id"`
	UserId int `db:"user_id" json:"user_id"`
	Status int `db:"status" json:"status"`
	TimeStamps
}

// type OrdersData struct {
// 	Id      int `db:"id"`
// 	OrderId int `db:"order_id"`
// 	Factory int `db:"factory"`
// 	DrinkId int `db:"drink_id"`
// 	count   int `db:"count"`
// }

type Factory struct {
	Id   int    `db:"id"`
	Name string `db:"name"`
	TimeStamps
}

type TimeStamps struct {
	UpdatedAt sql.NullString `db:"updated_at" json:"updated_at"`
	CreatedAt string `db:"created_at" json:"created_at"`
}

type Drinks struct {
	Id   int    `db:"id"`
	Name string `db:"name"`
	TimeStamps
}
