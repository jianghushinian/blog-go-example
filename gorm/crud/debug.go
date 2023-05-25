package main

import (
	"fmt"

	"gorm.io/gorm"
)

func DebugLogger(db *gorm.DB) error {
	// 1. 在打开连接时设置日志级别为 Info，可以打印所有的 SQL 记录
	// db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
	// 	Logger:logger.Default.LogMode(logger.Info),
	// })

	// 2. 打印单条 SQL 记录
	db.Debug().First(&User{})

	// 3. 打印慢查询 SQL 记录
	// slowLogger := logger.New(
	// 	log.New(os.Stdout, "\r\n", log.LstdFlags),
	// 	logger.Config{
	// 		// 设定慢查询时间阈值为 1ms（默认值：200 * time.Millisecond）
	// 		SlowThreshold: 1 * time.Microsecond,
	// 		// 设置日志级别
	// 		LogLevel: logger.Warn,
	// 	},
	// )
	//
	// db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
	// 	Logger: slowLogger,
	// })
	return nil
}

func DebugDryRun(db *gorm.DB) error {
	// 全局开启「空跑」模式
	// db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
	// 	DryRun: true,
	// })

	// 在不执行的情况下生成 SQL 及其参数，可以用于准备或测试生成的 SQL，详情请参考 Session
	var user User
	stmt := db.Session(&gorm.Session{DryRun: true}).First(&user, 1).Statement
	fmt.Println(stmt.SQL.String()) // => SELECT * FROM `users` WHERE `id` = $1 ORDER BY `id`
	fmt.Println(stmt.Vars)         // => []interface{}{1}
	return nil
}
