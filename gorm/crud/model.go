package main

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

/*
type User struct {
	ID           uint
	Name         string
	Email        *string
	Age          uint8
	Birthday     *time.Time
	MemberNumber sql.NullString
	ActivatedAt  sql.NullTime
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
*/

type User struct {
	gorm.Model
	Name         string         `gorm:"column:name"`
	Email        *string        `gorm:"column:email"`
	Age          uint8          `gorm:"column:age"`
	Birthday     *time.Time     `gorm:"column:birthday"`
	MemberNumber sql.NullString `gorm:"column:member_number"`
	ActivatedAt  sql.NullTime   `gorm:"column:activated_at"`
}

func (u *User) TableName() string {
	return "user"
}

type Post struct {
	gorm.Model
	Title    string     `gorm:"column:title"`
	Content  string     `gorm:"column:content"`
	Comments []*Comment `gorm:"foreignKey:PostID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;references:ID"`
	Tags     []*Tag     `gorm:"many2many:post_tags"`
}

func (p *Post) TableName() string {
	return "post"
}

type Comment struct {
	gorm.Model
	Content string `gorm:"column:content"`
	PostID  uint   `gorm:"column:post_id"`
	Post    *Post
}

func (c *Comment) TableName() string {
	return "comment"
}

type Tag struct {
	gorm.Model
	Name string  `gorm:"column:name"`
	Post []*Post `gorm:"many2many:post_tags"`
}

func (t *Tag) TableName() string {
	return "tag"
}
