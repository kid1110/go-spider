package service

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var (
	MyDB       *sql.DB
	initDBOnce sync.Once
)

func DBInit() {
	fmt.Println(initDBOnce)
	initDBOnce.Do(func() {
		ConnectDB()
	})
}

func ConnectDB() {
	fmt.Println(GetConfiguration().Mysql.Dsn)
	db, err := sql.Open("mysql", GetConfiguration().Mysql.Dsn)
	if err != nil {
		Mylog.Panicln("数据库连接错误:" + err.Error())
	}
	//最大空闲数
	db.SetMaxIdleConns(15)
	db.SetMaxOpenConns(50)
	db.SetConnMaxLifetime(5 * time.Minute)
	MyDB = db
}

func DisConnectDB() {
	if err := MyDB.Close(); err != nil {
		Mylog.Panicln("关闭数据库出错:" + err.Error())
	}

}
