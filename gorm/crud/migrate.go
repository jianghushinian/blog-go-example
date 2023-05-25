package main

import "gorm.io/gorm"

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&User{}, &Post{}, &Comment{}, &Tag{})
}
