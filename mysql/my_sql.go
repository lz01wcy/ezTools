package mysql

import (
	"fmt"
	"github.com/Anveena/ezTools/password"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Info struct {
	Host              string
	Port              string
	Account           string
	PasswordBase64Str string
	DatabaseName      string
}

func NewDBEngine(dbInfo *Info, gormConfig *gorm.Config, dbModels ...any) *gorm.DB {
	if dbInfo.PasswordBase64Str == "" {
		panic("密码没配置")
	}
	var pwd string
	pwd, err := password.Decode(dbInfo.PasswordBase64Str)
	if err != nil {
		panic(fmt.Sprintf("密码配的不合适,需要一个神秘的字符串才能解析,错误:\n\t%s", err.Error()))
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		dbInfo.Account, pwd, dbInfo.Host, dbInfo.Port, dbInfo.DatabaseName)
	db, err := gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		panic(err)
	}
	for _, ptr := range dbModels {
		if err = db.Migrator().AutoMigrate(ptr); err != nil {
			panic(fmt.Sprintf("建表失败!错误:\n\t%s", err.Error()))
		}
	}
	return db
}
