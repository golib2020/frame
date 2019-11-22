package f

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

const (
	mysqlInstancePrefix = `mysql`
	mysqlDefaultGroup   = `default`
)

//Mysql mysql实例
func Mysql(name ...string) *sql.DB {
	group := mysqlDefaultGroup
	if len(name) > 0 && name[0] != "" {
		group = name[0]
	}
	key := fmt.Sprintf("%s.%s", mysqlInstancePrefix, group)
	return Instance().GetOrSetFunc(key, func() interface{} {
		return mysqlInit(group)
	}).(*sql.DB)
}

func mysqlInit(group string) *sql.DB {
	conf := Config().Sub(fmt.Sprintf("db.%s", group))
	db, err := sql.Open(conf.GetString("driver"), conf.GetString("dsn"))
	if err != nil {
		log.Panic(err)
	}
	if err := db.Ping(); err != nil {
		log.Panic(err)
	}
	db.SetMaxOpenConns(conf.GetInt("max.open"))
	db.SetMaxIdleConns(conf.GetInt("max.idle"))
	db.SetConnMaxLifetime(conf.GetDuration("max.life"))
	return db
}
