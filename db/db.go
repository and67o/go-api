package db

import (
	"fmt"
	"go-api/config"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type DBOperations interface {
	AddUser(User User) (id int64, err error)
	GetUser(tgId int) (res User)
}

type DBManager struct {
	db *sqlx.DB
}

var Db DBOperations

func init() {
	db, err := sqlx.Connect("mysql", auth())
	if err != nil {
		panic(err)
	} else {
		log.Print("Db connect")
	}
	Db = &DBManager{db: db}
}

func auth() string {
	conf := config.New()
	fmt.Println(11, conf)
	login := conf.Mysql.Login
	if login == "" {
		panic("no login")
	}
	pass := conf.Mysql.Pass
	if pass == "" {
		panic("no password")
	}
	dbName := conf.Mysql.DbName
	return login + ":" + pass + "@/" + dbName
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
	DBM.db.Get(&user, "select * from users Where tg_id = ?", tgId)
	return
}
