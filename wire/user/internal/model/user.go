package model

import (
	"time"
)

// UserM user 表结构体映射
type UserM struct {
	ID        int64     `gorm:"column:id;primary_key"`
	Email     string    `gorm:"column:email"`
	Nickname  string    `gorm:"column:nickname"`
	Username  string    `gorm:"column:username;not null"`
	Password  string    `gorm:"column:password;not null"`
	CreatedAt time.Time `gorm:"column:createdAt"`
	UpdatedAt time.Time `gorm:"column:updatedAt"`
}

// TableName 映射 MySQL 表名
func (u *UserM) TableName() string {
	return "user"
}
