package multi

import (
	"database/sql"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

// NOTE: 多个依赖都有清理函数情况

// NewDatabaseConnection 创建数据库连接
func NewDatabaseConnection() (*sql.DB, func(), error) {
	db, err := sql.Open("mysql", "user:password@/dbname")
	if err != nil {
		return nil, nil, err
	}
	cleanup := func() {
		db.Close()
	}
	return db, cleanup, nil
}

// NewLogFile 创建一个用于日志记录的文件
func NewLogFile() (*os.File, func(), error) {
	file, err := os.Create("app.log")
	if err != nil {
		return nil, nil, err
	}
	cleanup := func() {
		file.Close()
	}
	return file, cleanup, nil
}

// App 应用程序结构
type App struct {
	DB  *sql.DB
	Log *os.File
}
