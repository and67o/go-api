package db

type BaseAction interface {
	Exec(sql string) (res interface{}, err error)
}

func(DBM *Manager) Exec(sql string) (interface{}, error)  {
	return DBM.db.Exec(sql)
}
