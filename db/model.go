package db

type User struct {
	Id   int    `db:"id"`
	TgId int    `db:"tg_id"`
	Name string `db:"name"`
}

type Order struct {
	Id     int `db:"id"`
	UserId int `db:"user_id"`
	Status int `db:"status"`
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
	UpdatedAt string `db:"updated_at"`
	CreatedAt string `db:"created_at"`
}

type Drinks struct {
	Id   int    `db:"id"`
	Name string `db:"name"`
	TimeStamps
}
