package db

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"go-api/config"
	"log"
)

type Operations interface {
	UserAction
	OrderAction
	BaseAction
}

type Manager struct {
	db *sqlx.DB
}

var Db Operations

func init() {
	db, err := sqlx.Connect("mysql", auth())
	if err != nil {
		panic(err)
	} else {
		log.Print("Db connect")
	}
	Db = &Manager{db: db}
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
