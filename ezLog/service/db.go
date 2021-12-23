package main

import (
	"database/sql"
	"fmt"
	"github.com/Anveena/ezTools/ezLog/ezLogPB"
	"github.com/Anveena/ezTools/ezPasswordEncoder"
	_ "github.com/go-sql-driver/mysql"
	"runtime"
	"time"
)

func startDBWritingThread() {
	runtime.LockOSThread()
	if ezLSConfig.MySQLConf.PasswordBase64Str == "" {
		panic(fmt.Errorf("MySQL未配置密码"))
	}
	password, err := ezPasswordEncoder.GetPasswordFromEncodedStr(ezLSConfig.MySQLConf.PasswordBase64Str)
	if err != nil {
		panic(fmt.Errorf("MySQL密码解析不出来,错误信息:%s", err.Error()))
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		ezLSConfig.MySQLConf.Account, password, ezLSConfig.MySQLConf.Host, ezLSConfig.MySQLConf.Port, ezLSConfig.MySQLConf.DatabaseName)
	db, err := sql.Open("mysql", dsn)
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	insertSQL := ""
	if err = _creatTable(db, &insertSQL); err != nil {
		panic(fmt.Errorf("建表失败,错误信息:%s", err.Error()))
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
		if err = _creatTable(db, &insertSQL); err != nil {
			fmt.Printf("建表失败,错误信息:%s", err.Error())
		}
	}()
	i := 0
	msgArr := make([]*ezLogPB.EZLogReq, howManyLogsToInsertDBOnce)
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
			stmtIns, err := db.Prepare(insertSQL) // ? = placeholder
			if err != nil {
				fmt.Printf("插入数据失败了!错误:\n\t%s\n", err.Error())
			}
			for j := 0; j < i; j++ {
				if _, err = stmtIns.Exec(msgArr[j].Level, msgArr[j].AppName, msgArr[j].FileName, msgArr[j].FileLine, msgArr[j].Tag, msgArr[j].Time.AsTime(), msgArr[j].Content); err != nil {
					fmt.Printf("插入数据失败了!错误:\n\t%s\n", err.Error())
					break
				}
			}
			if err = stmtIns.Close(); err != nil {
				fmt.Printf("插入数据失败了!错误:\n\t%s\n", err.Error())
			}
		}
		select {
		case <-tickerToNewTable.C:
			toDelDate := time.Now().Add(-time.Hour * 23 * time.Duration(ezLSConfig.HowManyDaysThatLogsShouldSave))
			toDelTableName := _getTableName(toDelDate)
			if _, e := db.Exec(fmt.Sprintf("drop table if exists %s", toDelTableName)); e != nil {
				fmt.Printf("删除表失败了!错误:\n\t%s\n", e.Error())
			}
			if err = _creatTable(db, &insertSQL); err != nil {
				fmt.Printf("建表失败,错误信息:%s", err.Error())
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
func _creatTable(db *sql.DB, insertSQL *string) error {
	tableName := _getTableName(time.Now())
	if _, err := db.Exec(fmt.Sprintf(
		"create table if not exists `%s`"+
			"("+
			"    `id`      bigint auto_increment comment '无意义',"+
			"    `level`   int comment '日志等级:1-debug,2-info,3-err,4-ding,5-ding_list,6-ding_all',"+
			"    `name`    char(32) comment 'app名字',"+
			"    `file`    char(255) comment '代码文件',"+
			"    `line`    int comment '代码行',"+
			"    `tag`     char(127) comment '日志标签',"+
			"    `time`    datetime(3) null comment '日志时间',"+
			"    `content` text(2048) comment '具体日志',"+
			"    primary key (`id`),"+
			"    index tag_with_level_of_app (`tag`, `level`, `name`),"+
			"    index file_name_of_app_query (`file`, `name`),"+
			"    index file_line_of_app_query (`line`, `name`)"+
			")", tableName)); err != nil {
		return err
	}
	*insertSQL = fmt.Sprintf("insert into `%s` (level, name, file, line, tag, time, content) values (?,?,?,?,?,?,?)", tableName)
	return nil
}
