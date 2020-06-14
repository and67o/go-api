package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load("go-api/.env"); err != nil {
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
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

func getEnvAsBool(name string, defaultVal bool) bool {
	valStr := getEnv(name, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}

	return defaultVal
}

func getEnvAsInt(name string, defaultVal int) int {
	valStr := getEnv(name, "")
	if val, err := strconv.Atoi(valStr); err == nil {
		return val
	}

	return defaultVal
}
