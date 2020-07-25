package config

type MysqlConfig struct {
	Login  string
	Pass   string
	DbName string
}

type RedisConfig struct {
	Addr string
	Pass string
	Db int
}

type Config struct {
	Mysql MysqlConfig
	Redis RedisConfig
}
