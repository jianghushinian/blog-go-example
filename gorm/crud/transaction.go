package main

import "gorm.io/gorm"

func TransactionPost(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		post := Post{
			Title: "Hello World",
		}
		if err := tx.Create(&post).Error; err != nil {
			return err
		}
		comment := Comment{
			Content: "Hello World",
			PostID:  post.ID,
		}
		if err := tx.Create(&comment).Error; err != nil {
			return err
		}
		return nil
	})
}

func TransactionPostWithManually(db *gorm.DB) error {
	// 手动事务，事务一旦开始，你就应该使用 tx 处理数据
	tx := db.Begin()

	post := Post{
		Title: "Hello World Manually",
	}
	if err := tx.Create(&post).Error; err != nil {
		tx.Rollback()
		return err
	}
	comment := Comment{
		Content: "Hello World Manually",
		PostID:  post.ID,
	}
	if err := tx.Create(&comment).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
