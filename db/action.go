package db

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

func (DBM *DBManager) GetOrders() (orders Order) {
	orders = Order{}
	DBM.db.Get(&orders, "SELECT * FROM users")
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

func (DBM *DBManager) GetUser(tgId int) (user User) {
	user = User{}
	DBM.db.Get(&user, "SELECT * FROM users WHERE tg_id = ?", tgId)
	return
}
