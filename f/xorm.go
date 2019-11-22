package f

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

const (
	XormInstancePrefix = `xorm`
	XormDefaultGroup   = `default`
)

func Xorm(name ...string) *xorm.Engine {
	group := XormDefaultGroup
	if len(name) > 0 && name[0] != "" {
		group = name[0]
	}
	key := fmt.Sprintf("%s.%s", XormInstancePrefix, group)
	return Instance().GetOrSetFunc(key, func() interface{} {
		return xormInit(group)
	}).(*xorm.Engine)
}

func xormInit(group string) *xorm.Engine {
	conf := Config().Sub(fmt.Sprintf("db.%s", group))
	db, err := xorm.NewEngine(conf.GetString("driver"), conf.GetString("dsn"))
	if err != nil {
		log.Panic(err)
	}
	db.Logger().SetLevel(2)
	if err = db.Ping(); err != nil {
		log.Panic(err)
	}
	db.SetMaxOpenConns(conf.GetInt("max.open"))
	db.SetMaxIdleConns(conf.GetInt("max.idle"))
	db.SetConnMaxLifetime(conf.GetDuration("max.life"))
	return db
}
