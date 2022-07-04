package db

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
)

const MaxMysqlConnect = 8

var (
	OrmEngine *xorm.Engine
)

// InitOrmMysql 初始化orm
func InitOrmMysql(user string, pwd string, host string, port int, dbs string) error {
	connInfo := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", user, pwd, host, port, dbs)
	temp, err := xorm.NewEngine("mysql", connInfo)
	if err != nil {
		return err
	}
	OrmEngine = temp
	OrmEngine.Logger().SetLevel(core.LOG_WARNING)
	if err := OrmEngine.Ping(); err != nil {
		return err
	}
	OrmEngine.SetMaxOpenConns(MaxMysqlConnect)
	OrmEngine.SetMaxIdleConns(MaxMysqlConnect)

	return nil
}
