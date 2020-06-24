package db

import "fmt"

func (DBM *DBManager) GetOrders() (orders []Order) {
	orders = []Order{}
	err := DBM.db.Select(&orders, "SELECT * FROM orders")
	for i, v := range orders {
		fmt.Println(i, v)
	}
	fmt.Println(err)
	return
}

func (DBM *DBManager) DeleteOrder(orderId int) (res bool, err error) {
	_, err = DBM.db.Exec("DELETE FROM orders WHERE order_id = ?", orderId)
	if err != nil {
		panic(err)
	}
	return
}

func (DBM *DBManager) GetOrder(orderId int) (order Order) {
	order = Order{}
	DBM.db.Get(&order, "SELECT * FROM users WHERE tg_id = ?", orderId)
	return
}

func (DBM *DBManager) AddUser(user User) (id int64, err error) {
	res, err := DBM.db.Exec("INSERT INTO users (name,tg_id) VALUES(?, ?)", user.Name, user.TgId)
	if err != nil {
		panic(err)
	}
	id, err = res.LastInsertId()
	if err != nil {
		panic(err)
	}
	return
}

func (DBM *DBManager) GetUser(tgId int) (user User) {
	user = User{}
	DBM.db.Get(&user, "SELECT * FROM users WHERE tg_id = ?", tgId)
	return
}

func (DBM *DBManager) GetUsers() (users []User) {
	users = []User{}
	err := DBM.db.Select(&users, "SELECT * FROM users")
	for i, v := range users {
		fmt.Println(i, v)
	}
	fmt.Println(err)
	return
}

func (DBM *DBManager) DeleteUser(userId int) (res bool, err error) {
	_, err = DBM.db.Exec("DELETE FROM users WHERE id = ?", userId)
	if err != nil {
		panic(err)
	}
	return
}

// func (DBM *DBManager) UpdateUser(userId int) (res bool, err error) {
// 	_, err = DBM.db.Exec("UPDAE users set WHERE id = ?", userId)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return
// }
