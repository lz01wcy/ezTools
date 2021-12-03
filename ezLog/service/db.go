package main

import (
	"fmt"
	"github.com/Anveena/ezTools/ezPasswordEncoder"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"runtime"
	"time"
)

type logModel struct {
	ID       int64     `gorm:"column:id;comment:无意义;autoIncrement;primaryKey"`
	Level    int32     `gorm:"colum:lv;comment:日志等级:1-debug,2-info,3-err,4-ding,5-ding_list,6-ding_all;index:tag_with_level_of_app,priority:2"`
	AppName  string    `gorm:"type:char(32);column:from;comment:App名字;index:file_name_of_app_query,priority:2;index:file_line_of_app_query,priority:2;index:tag_with_level_of_app,priority:3"`
	FileName string    `gorm:"type:char(255);column:file;comment:代码文件;index:file_name_of_app_query,priority:1"`
	FileLine int32     `gorm:"column:line;comment:代码行;index:file_line_of_app_query,priority:1"`
	Tag      string    `gorm:"type:char(127);column:tag;comment:日志标签;index:tag_with_level_of_app,priority:1"`
	Time     time.Time `gorm:"column:time;comment:日志时间"`
	Content  string    `gorm:"type:text(2048);column:content;comment:具体日志"`
}

func (lm *logModel) UpdateToDB() {
	logModelChan <- lm
}
func (lm *logModel) TableName() string {
	return _getTableName(time.Now())
}
func startDBWritingThread() (err error) {
	runtime.LockOSThread()
	if ezLSConfig.MySQLConf.PasswordBase64Str == "" {
		return fmt.Errorf("MySQL未配置密码")
	}
	var password string
	password, err = ezPasswordEncoder.GetPasswordFromEncodedStr(ezLSConfig.MySQLConf.PasswordBase64Str)
	if err != nil {
		return fmt.Errorf("MySQL密码解析不出来,错误信息:%s", err.Error())
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		ezLSConfig.MySQLConf.Account, password, ezLSConfig.MySQLConf.Host, ezLSConfig.MySQLConf.Port, ezLSConfig.MySQLConf.DatabaseName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		return err
	}
	if err = _creatTable(&db); err != nil {
		return err
	}
	howManyLogsToInsertDBOnce := ezLSConfig.HowManyLogsToInsertDBOnce
	tickerToWrite := time.NewTicker(time.Second * time.Duration(ezLSConfig.HowOftenToInsertDBInSeconds))
	now := time.Now()
	next := now.Add(time.Hour * 24)
	next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())
	tickerToNewTable := time.NewTicker(time.Hour * 24 * 365)
	timer := time.NewTimer(next.Sub(time.Now()))
	go func() {
		<-timer.C
		tickerToNewTable.Reset(time.Hour * 24)
		if err = _creatTable(&db); err != nil {
			fmt.Printf("第一次建表失败,错误信息:%s", err.Error())
		}
	}()
	i := 0
	msgArr := make([]*logModel, howManyLogsToInsertDBOnce)
	for {
		i = 0
	outer:
		for ; i < howManyLogsToInsertDBOnce; i++ {
			select {
			case <-tickerToWrite.C:
				break outer
			case msgArr[i] = <-logModelChan:
				break
			}
		}
		if i > 0 {
			if e := db.Create(msgArr[:i]).Error; e != nil {
				fmt.Printf("插入数据失败了!错误:\n\t%s\n", e.Error())
			}
		}
		select {
		case <-tickerToNewTable.C:
			toDelDate := time.Now().Add(-time.Hour * 23 * time.Duration(ezLSConfig.HowManyDaysThatLogsShouldSave))
			toDelTableName := _getTableName(toDelDate)
			if e := db.Exec(fmt.Sprintf("drop table if exists %s", toDelTableName)).Error; e != nil {
				fmt.Printf("删除表失败了!错误:\n\t%s\n", e.Error())
			}
			if err = _creatTable(&db); err != nil {
				return err
			}
			break
		default:
			break
		}
	}
}
func _getTableName(t time.Time) string {
	return fmt.Sprintf("logs_of_%d_%02d_%02d", t.Year(), t.Month(), t.Day())
}
func _creatTable(db **gorm.DB) error {
	sc := (*db).Scopes(func(db *gorm.DB) *gorm.DB {
		return db.Table(_getTableName(time.Now()))
	})
	*db = sc
	if !sc.Migrator().HasTable(&logModel{}) {
		if err := sc.Migrator().CreateTable(&logModel{}); err != nil {
			return fmt.Errorf("db table create failed with err:%s", err.Error())
		}
	}
	return nil
}
