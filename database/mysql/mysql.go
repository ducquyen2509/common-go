package storage

import (
	"fmt"
	cfg "github.com/ducquyen2509/common-go/config"
	"github.com/ducquyen2509/common-go/logger"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"time"
)

// 	To make this helper create tables, set environment variable "DB_MIGRATE" to true
func InitDatabase(config *cfg.Mysql) error {
	return initDatabase(config.Url, config.User, config.Pass, config.Name, config.Alias, config.Debug)
}

func initDatabase(addr, user, pass, dbName, alias string, debug bool) error {

	aliasDb := "default"
	if alias != "" {
		aliasDb = alias
	}
	var connectStr = fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?charset=utf8&collation=utf8mb4_general_ci&loc=%s",
		user, pass, addr, dbName, "Asia%2fBangkok")
	if err := orm.RegisterDriver("mysql", orm.DRMySQL); err != nil {
		return err
	}

	logger.Debugf("Connect database with connect string: %s", connectStr)
	orm.Debug, orm.DebugLog = debug, orm.NewLog(os.Stdout)
	if err := orm.RegisterDataBase(aliasDb, "mysql", connectStr, 10, 20); err != nil {
		return err
	}

	if db, err := orm.GetDB(aliasDb); err != nil {
		return err
	} else {
		db.SetConnMaxLifetime(time.Duration(300) * time.Second)
	}

	return nil
}
