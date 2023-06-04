package main

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type User struct {
	ID int
	// 如果 Name 为 string 类型，当数据库中为 NULL 值时，
	// Scan 方法会报错：sql: Scan error on column index 1, name "name": converting NULL to string is unsupported
	// |         value          | value for MySQL |
	// |------------------------|-----------------|
	// | {String:n1 Valid:true} |      'n1'       |
	// | {String: Valid:true}   |       ''        |
	// | {String: Valid:false}  |      NULL       |
	Name  sql.NullString
	Email string
	Age   int
	// 使用指针类型可以 Scan 数据库中的 NULL 类型为 nil
	// 业务端也可以用来检查字段是否有值
	Birthday *time.Time
	// 实现了 sql.Scanner、driver.Valuer 类型
	// |            value             |          value for MySQL         |
	// |------------------------------|----------------------------------|
	// | {Month:0 Year:0}             | {"month":0 "year":0}/NULL        |
	// | {Month:100000 Year:10000000} | {"month":100000,"year":10000000} |
	Salary Salary
	// 使用 time.Time 映射 MySQL 中的 datetime 类型字段，DSN 需要指定 parseTime=true 参数
	// 2023-06-03 08:48:35 +0800 CST
	CreatedAt time.Time
	// 也可以使用 string 映射 MySQL 中的 datetime 类型字段
	// 2023-06-03T08:48:35+08:00
	UpdatedAt string
}

type Salary struct {
	Month int `json:"month"`
	Year  int `json:"year"`
}

// Scan implements sql.Scanner
func (s *Salary) Scan(src any) error {
	if src == nil {
		return nil
	}

	var buf []byte
	switch v := src.(type) {
	case []byte:
		buf = v
	case string:
		buf = []byte(v)
	default:
		return fmt.Errorf("invalid type: %T", src)
	}

	err := json.Unmarshal(buf, s)
	return err
}

// Value implements driver.Valuer
func (s Salary) Value() (driver.Value, error) {
	v, err := json.Marshal(s)
	return string(v), err
}
