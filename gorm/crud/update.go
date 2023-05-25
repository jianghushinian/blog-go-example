package main

import (
	"fmt"

	"gorm.io/gorm"
)

func UpdateUser(db *gorm.DB) error {
	var user User
	// 先查询用户
	// SELECT * FROM `user` WHERE `user`.`deleted_at` IS NULL ORDER BY `user`.`id` LIMIT 1
	result := db.First(&user)
	fmt.Printf("affected rows: %d\n", result.RowsAffected)
	if err := result.Error; err != nil {
		return err
	}

	// 更新除主键的全部列（即使字段是零值）
	user.Name = "John"
	user.Age = 20
	// UPDATE `user` SET `created_at`='2023-05-22 22:14:47.814',`updated_at`='2023-05-22 22:24:34.201',`deleted_at`=NULL,`name`='John',`email`='u1@jianghushinian.com',`age`=20,`birthday`='2023-05-22 22:14:47.813',`member_number`=NULL,`activated_at`=NULL WHERE `user`.`deleted_at` IS NULL AND `id` = 1
	result = db.Save(&user)
	fmt.Printf("affected rows: %d\n", result.RowsAffected)
	if err := result.Error; err != nil {
		return err
	}

	// 更新指定列
	// UPDATE `user` SET `name`='Jianghushinian',`updated_at`='2023-05-22 22:24:34.215' WHERE `user`.`deleted_at` IS NULL AND `id` = 1
	result = db.Model(&user).Update("name", "Jianghushinian")
	fmt.Printf("affected rows: %d\n", result.RowsAffected)
	if err := result.Error; err != nil {
		return err
	}

	// 使用表达式
	// UPDATE `user` SET `age`=age + 1,`updated_at`='2023-05-22 22:24:34.219' WHERE `user`.`deleted_at` IS NULL AND `id` = 1
	result = db.Model(&user).Update("age", gorm.Expr("age + ?", 1))
	fmt.Printf("affected rows: %d\n", result.RowsAffected)
	if err := result.Error; err != nil {
		return err
	}

	// 更新多个指定列（只会更新非零值字段）
	// UPDATE `user` SET `updated_at`='2023-05-22 22:29:35.19',`name`='JiangHu' WHERE `user`.`deleted_at` IS NULL AND `id` = 1
	result = db.Model(&user).Updates(User{Name: "JiangHu", Age: 0})
	fmt.Printf("affected rows: %d\n", result.RowsAffected)
	if err := result.Error; err != nil {
		return err
	}

	//  更新多个指定列（更新零值字段）
	// UPDATE `user` SET `age`=0,`name`='JiangHu',`updated_at`='2023-05-22 22:29:35.623' WHERE `user`.`deleted_at` IS NULL AND `id` = 1
	result = db.Model(&user).Updates(map[string]interface{}{"name": "JiangHu", "age": 0})
	fmt.Printf("affected rows: %d\n", result.RowsAffected)
	return result.Error
}

func UpdatePost(db *gorm.DB) error {
	var post Post
	// SELECT * FROM `post` WHERE `post`.`deleted_at` IS NULL ORDER BY `post`.`id` LIMIT 1
	result := db.First(&post)
	if err := result.Error; err != nil {
		return err
	}
	fmt.Println(post)
	for i, comment := range post.Comments {
		fmt.Println(i, comment)
	}

	// 替换关联
	comment := Comment{
		Content: "comment3",
	}
	// INSERT INTO `comment` (`created_at`,`updated_at`,`deleted_at`,`content`,`post_id`) VALUES ('2023-05-23 09:07:42.852','2023-05-23 09:07:42.852',NULL,'comment3',1) ON DUPLICATE KEY UPDATE `post_id`=VALUES(`post_id`)
	// UPDATE `post` SET `updated_at`='2023-05-23 09:07:42.846' WHERE `post`.`deleted_at` IS NULL AND `id` = 1
	// UPDATE `comment` SET `post_id`=NULL WHERE `comment`.`id` <> 8 AND `comment`.`post_id` = 1 AND `comment`.`deleted_at` IS NULL
	err := db.Model(&post).Association("Comments").Replace([]*Comment{&comment})
	if err != nil {
		return err
	}
	fmt.Println(post)
	for i, comment := range post.Comments {
		fmt.Println(i, comment)
	}

	// var post1 Post
	// result = db.First(&post1)
	// if err := result.Error; err != nil {
	// 	return err
	// }
	// fmt.Println(post1)
	// for _, comment := range post1.Comments {
	// 	fmt.Println(comment)
	// }
	return nil
}
