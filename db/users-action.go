package db

import (
	"errors"
	"time"
)

type UserAction interface {
	AddUser(name string, tgId int) (user User, err error)
	GetUser(tgId int) (res User, err error)
	GetUsers() (users []User, err error)
	DeleteUsers() (res bool, err error)
}

func (DBM *Manager) AddUser(name string, tgId int) (user User, err error) {
	createdAt := time.Now().Format("2006-01-02 15:04:05")

	res, err := DBM.db.Exec("INSERT INTO users (name, tg_id, created_at) VALUES(?, ?, ?)", name, tgId, createdAt)
	if err != nil {
		panic(err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		panic(err)
	}
	return User{Id: id, TgId: tgId, Name: name, CreatedAt: createdAt}, nil
}

func (DBM *Manager) GetUser(userId int) (user User, err error) {
	user = User{}
	DBM.db.Get(&user, "SELECT * FROM users WHERE id = ?", userId)
	if user.TgId == 0 {
		err = errors.New("no user")
	}
	return
}

func (DBM *Manager) GetUsers() (users []User, err error) {
	err = DBM.db.Select(&users, "SELECT * FROM users")
	if len(users) == 0 {
		err = errors.New("no users")
	}
	return
}

func (DBM *Manager) DeleteUser(userId int) (res bool, err error) {
	_, err = DBM.db.Exec("DELETE FROM users WHERE id = ?", userId)
	if err != nil {
		panic(err)
	}
	return
}

func (DBM *Manager) DeleteUsers() (res bool, err error) {
	_, err = DBM.db.Exec("DELETE FROM users where 1")
	if err != nil {
		panic(err)
	}
	return
}
