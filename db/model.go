package db

import "database/sql"

type User struct {
	Id   int64    `db:"id" json:"id"`
	TgId int    `db:"tg_id" json:"tg_id"`
	Name string `db:"name" json:"name"`
	CreatedAt string `db:"created_at" json:"created_at"`
	UpdatedAt sql.NullString `db:"updated_at" json:"updated_at"`
}

type Order struct {
	Id     int `db:"id" json:"id"`
	UserId int `db:"user_id" json:"user_id"`
	Status int `db:"status" json:"status"`
	CreatedAt string `db:"created_at" json:"created_at"`
	UpdatedAt sql.NullString `db:"updated_at" json:"updated_at"`
}

type Factory struct {
	Id   int    `db:"id"`
	Name string `db:"name"`
	CreatedAt string `db:"created_at" json:"created_at"`
	UpdatedAt sql.NullString `db:"updated_at" json:"updated_at"`
}

type Drinks struct {
	Id   int    `db:"id"`
	Name string `db:"name"`
	CreatedAt string `db:"created_at" json:"created_at"`
	UpdatedAt sql.NullString `db:"updated_at" json:"updated_at"`
}
