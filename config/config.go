package config

import (
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load("/var/www/golang/src/go-api/.env"); err != nil {
		panic("No .env file found")
	}
}

func New() *Config {
	return &Config{
		Mysql: MysqlConfig{
			Login:  getEnv("MYSQL_LOGIN", ""),
			Pass:   getEnv("MYSQL_PASS", ""),
			DbName: getEnv("MYSQL_DB_NAME", ""),
		},
		Redis: RedisConfig{
			Addr: "localhost:6379",
			Pass: "",
			Db: 0,
		},
	}
}