package main

import (
	"fmt"

	"gorm.io/gorm"
)

func DeleteUser(db *gorm.DB) error {
	// 软删除（逻辑删除）
	var user User
	// UPDATE `user` SET `deleted_at`='2023-05-22 22:46:45.086' WHERE name = 'JiangHu' AND `user`.`deleted_at` IS NULL
	result := db.Where("name = ?", "JiangHu").Delete(&user)
	fmt.Printf("affected rows: %d\n", result.RowsAffected)
	if err := result.Error; err != nil {
		return err
	}

	// 无法直接查询到被软删除的记录
	// SELECT * FROM `user` WHERE name = 'JiangHu' AND `user`.`deleted_at` IS NULL ORDER BY `user`.`id` LIMIT 1
	result = db.Where("name = ?", "JiangHu").First(&user)
	if err := result.Error; err != nil {
		fmt.Println(err) // record not found
		// return err
	}
	fmt.Println(user)

	// 使用 Unscoped 能够查询被软删除的记录
	// SELECT * FROM `user` WHERE name = 'JiangHu' ORDER BY `user`.`id` LIMIT 1
	result = db.Unscoped().Where("name = ?", "JiangHu").First(&user)
	fmt.Printf("affected rows: %d\n", result.RowsAffected)
	if err := result.Error; err != nil {
		return err
	}
	fmt.Println(user)

	// 永久删除（物理删除）
	// DELETE FROM `user` WHERE name = 'JiangHu' AND `user`.`id` = 1
	result = db.Unscoped().Where("name = ?", "JiangHu").Delete(&user)
	fmt.Printf("affected rows: %d\n", result.RowsAffected)
	return result.Error
}

func DeletePost(db *gorm.DB) error {
	// 删除 post，不会影响 comment、tag
	var post Post
	// UPDATE `post` SET `deleted_at`='2023-05-23 09:09:58.534' WHERE id = 1 AND `post`.`deleted_at` IS NULL
	result := db.Where("id = ?", 1).Delete(&post)
	fmt.Printf("affected rows: %d\n", result.RowsAffected)
	if err := result.Error; err != nil {
		return err
	}

	var post2 Post
	// SELECT * FROM `post` WHERE `post`.`deleted_at` IS NULL ORDER BY `post`.`id` DESC LIMIT 1
	// SELECT * FROM `comment` WHERE `comment`.`post_id` = 6 AND `comment`.`deleted_at` IS NULL
	result = db.Preload("Comments").Last(&post2)
	if err := result.Error; err != nil {
		return err
	}
	fmt.Println(post2)

	// 删除关联
	// UPDATE `comment` SET `post_id`=NULL WHERE `comment`.`post_id` = 6 AND `comment`.`id` IN (NULL) AND `comment`.`deleted_at` IS NULL
	err := db.Model(&post2).Association("Comments").Delete(post2.Comments)
	if err != nil {
		return err
	}
	fmt.Println(post2)
	for _, comment := range post2.Comments {
		fmt.Println(comment)
	}
	return nil
}
