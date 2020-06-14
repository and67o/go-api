package db

import (
	"go-api/config"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type DBOperations interface {
	AddUser(User User) (id int64, err error)
	GetUser(tgId int) (res User)
	GetOrders() (orders Order)
	GetOrder(orderId int) (order Order)
	DeleteOrder(orderId int) (res bool, err error)
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
