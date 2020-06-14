package config

type MysqlConfig struct {
	Login  string
	Pass   string
	DbName string
}

type Config struct {
	Mysql MysqlConfig
}
