package main

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

/*
|钩子函数|执行时机|
|:------:|:------:|
| BeforeSave | 调用	Save 前 |
| AfterSave | 调用 Save 后 |
| BeforeCreate | 插入记录前 |
| AfterCreate | 插入记录后 |
| BeforeUpdate | 更新记录前 |
| AfterUpdate | 更新记录后 |
| BeforeDelete | 删除记录前 |
| AfterDelete | 删除记录后 |
| AfterFind | 查询记录后 |
*/

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	// u.UUID = uuid.New()
	fmt.Println("BeforeCreate")
	if u.Name == "admin" {
		return errors.New("invalid name")
	}
	return nil
}

func (u *User) AfterCreate(tx *gorm.DB) (err error) {
	fmt.Println("AfterCreate")
	return nil
}

func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
	fmt.Println("BeforeUpdate")
	return nil
}

func (u *User) AfterUpdate(tx *gorm.DB) (err error) {
	fmt.Println("AfterUpdate")
	return nil
}

func (u *User) BeforeDelete(tx *gorm.DB) (err error) {
	fmt.Println("BeforeDelete")
	return nil
}

func (u *User) AfterDelete(tx *gorm.DB) (err error) {
	fmt.Println("AfterDelete")
	return nil
}

func (u *User) AfterFind(tx *gorm.DB) (err error) {
	fmt.Println("AfterFind")
	return nil
}

func (u *User) BeforeSave(tx *gorm.DB) (err error) {
	fmt.Println("BeforeSave")
	return nil
}

func (u *User) AfterSave(tx *gorm.DB) (err error) {
	fmt.Println("AfterSave")
	return nil
}
