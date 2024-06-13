package db

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// MySQLOptions MySQL 数据库的选项
type MySQLOptions struct {
	Address  string
	Database string
	Username string
	Password string
}

// DSN 从 MySQLOptions 返回 DSN
func (o *MySQLOptions) DSN() string {
	return fmt.Sprintf(`%s:%s@tcp(%s)/%s?charset=utf8&parseTime=%t&loc=%s`,
		o.Username,
		o.Password,
		o.Address,
		o.Database,
		true,
		"Local")
}

// NewMySQL 根据选项构造 *gorm.DB
func NewMySQL(opts *MySQLOptions) (*gorm.DB, func(), error) {
	// 可以用来释放资源，这里仅作为示例使用，没有释放任何资源，因为 gorm 内部已经帮我们做了
	cleanFunc := func() {}

	db, err := gorm.Open(mysql.Open(opts.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	return db, cleanFunc, err
}
