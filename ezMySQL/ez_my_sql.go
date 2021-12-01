package ezMySQL

import (
	"fmt"
	"github.com/Anveena/ezTools/ezPasswordEncoder"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Info struct {
	Host              string
	Port              string
	Account           string
	PasswordBase64Str string
	DatabaseName      string
}

func NewDBEngine(dbInfo *Info, dbModels ...interface{}) (*gorm.DB, error) {
	if dbInfo.PasswordBase64Str == "" {
		return nil, fmt.Errorf("密码没配置")
	}
	var password string
	password, err := ezPasswordEncoder.GetPasswordFromEncodedStr(dbInfo.PasswordBase64Str)
	if err != nil {
		return nil, fmt.Errorf("密码配的不合适,需要一个神秘的字符串才能解析,错误:\n\t%s", err.Error())
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		dbInfo.Account, password, dbInfo.Host, dbInfo.Port, dbInfo.DatabaseName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		return nil, err
	}
	for _, ptr := range dbModels {
		if !db.Migrator().HasTable(ptr) {
			if err := db.Migrator().CreateTable(ptr); err != nil {
				return nil, fmt.Errorf("建表失败!错误:\n\t%s", err.Error())
			}
		}
	}
	return db, nil
}
