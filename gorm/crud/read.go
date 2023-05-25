package main

import (
	"fmt"

	"gorm.io/gorm"
)

func ReadUser(db *gorm.DB) error {
	// 查询第一条记录（主键升序）
	var user User
	// SELECT * FROM `user` WHERE `user`.`deleted_at` IS NULL ORDER BY `user`.`id` LIMIT 1
	result := db.First(&user)
	fmt.Printf("affected rows: %d\n", result.RowsAffected)
	if err := result.Error; err != nil {
		return err
	}
	fmt.Println(user)

	// 查询最后一条记录
	var lastUser User
	// SELECT * FROM `user` WHERE `user`.`deleted_at` IS NULL ORDER BY `user`.`id` DESC LIMIT 1
	result = db.Last(&lastUser)
	fmt.Printf("affected rows: %d\n", result.RowsAffected)
	if err := result.Error; err != nil {
		return err
	}
	fmt.Println(lastUser)

	// 条件查询
	var users []User
	// SELECT * FROM `user` WHERE name != 'unknown' AND `user`.`deleted_at` IS NULL
	result = db.Where("name != ?", "unknown").Find(&users)
	fmt.Printf("affected rows: %d\n", result.RowsAffected)
	if err := result.Error; err != nil {
		return err
	}
	fmt.Println(users)

	// 查询指定字段
	var user2 User
	// SELECT `name`,`age` FROM `user` WHERE `user`.`deleted_at` IS NULL ORDER BY `user`.`id` LIMIT 1
	result = db.Select("name", "age").First(&user2)
	fmt.Printf("affected rows: %d\n", result.RowsAffected)
	if err := result.Error; err != nil {
		return err
	}
	fmt.Println(user2)

	// 排序
	var users2 []User
	// SELECT * FROM `user` WHERE `user`.`deleted_at` IS NULL ORDER BY id desc
	result = db.Order("id desc").Find(&users2)
	fmt.Printf("affected rows: %d\n", result.RowsAffected)
	if err := result.Error; err != nil {
		return err
	}
	fmt.Println(users2)

	// Limit & Offset
	var users3 []User
	// SELECT * FROM `user` WHERE `user`.`deleted_at` IS NULL LIMIT 2 OFFSET 1
	result = db.Limit(2).Offset(1).Find(&users3)
	fmt.Printf("affected rows: %d\n", result.RowsAffected)
	if err := result.Error; err != nil {
		return err
	}
	fmt.Println(users3)
	// 使用 -1 可以取消限制条件
	var users4 []User
	var users5 []User
	// SELECT * FROM `user` WHERE `user`.`deleted_at` IS NULL LIMIT 2 OFFSET 1; (users4)
	// SELECT * FROM `user` WHERE `user`.`deleted_at` IS NULL; (users5)
	result = db.Limit(2).Offset(1).Find(&users4).Limit(-1).Offset(-1).Find(&users5)
	fmt.Printf("affected rows: %d\n", result.RowsAffected)
	if err := result.Error; err != nil {
		return err
	}
	fmt.Println(users4)
	fmt.Println(users5)

	// Count
	var count int64
	// SELECT count(*) FROM `user` WHERE `user`.`deleted_at` IS NULL
	result = db.Model(&User{}).Count(&count)
	fmt.Printf("affected rows: %d\n", result.RowsAffected)
	if err := result.Error; err != nil {
		return err
	}
	fmt.Println(count)

	// 子查询
	var avgages []float64
	// SELECT AVG(age) as avgage FROM `user` WHERE `user`.`deleted_at` IS NULL GROUP BY `name` HAVING AVG(age) > (SELECT AVG(age) FROM `user` WHERE name LIKE 'user%')
	subQuery := db.Select("AVG(age)").Where("name LIKE ?", "user%").Table("user")
	result = db.Model(&User{}).Select("AVG(age) as avgage").Group("name").Having("AVG(age) > (?)", subQuery).Find(&avgages)
	if err := result.Error; err != nil {
		return err
	}
	fmt.Println(avgages)
	return nil
}

func ReadPost(db *gorm.DB) error {
	var (
		post     Post
		comments []*Comment
	)

	// 根据 post 的 id 查询 comments
	post.ID = 1
	// SELECT * FROM `comment` WHERE `comment`.`post_id` = 1 AND `comment`.`deleted_at` IS NULL
	err := db.Model(&post).Association("Comments").Find(&comments)
	for i, comment := range comments {
		fmt.Println(i, comment)
	}
	if err != nil {
		return err
	}

	// 预加载
	post2 := Post{}
	// SELECT * FROM `post` WHERE `post`.`deleted_at` IS NULL ORDER BY `post`.`id` LIMIT 1
	// SELECT * FROM `comment` WHERE `comment`.`post_id` = 1 AND `comment`.`deleted_at` IS NULL
	// SELECT * FROM `post_tags` WHERE `post_tags`.`post_id` = 1
	// SELECT * FROM `tag` WHERE `tag`.`id` IN (1,2) AND `tag`.`deleted_at` IS NULL
	err = db.Preload("Comments").Preload("Tags").First(&post2).Error
	if err != nil {
		return err
	}
	fmt.Println(post2)
	for i, comment := range post2.Comments {
		fmt.Println(i, comment)
	}
	for i, tag := range post2.Tags {
		fmt.Println(i, tag)
	}

	// JOIN
	type PostComment struct {
		Title   string
		Comment string
	}
	postComment := PostComment{}
	post3 := Post{}
	post3.ID = 3
	// SELECT post.title, comment.Content AS comment FROM `post` LEFT JOIN comment ON comment.post_id = post.id WHERE `post`.`deleted_at` IS NULL AND `post`.`id` = 3
	result := db.Model(&post3).Select("post.title, comment.Content AS comment").Joins("LEFT JOIN comment ON comment.post_id = post.id").Scan(&postComment)
	if err := result.Error; err != nil {
		return err
	}
	fmt.Println(postComment)

	// JOIN 预加载
	var comments2 []*Comment
	// SELECT `comment`.`id`,`comment`.`created_at`,`comment`.`updated_at`,`comment`.`deleted_at`,`comment`.`content`,`comment`.`post_id`,`Post`.`id` AS `Post__id`,`Post`.`created_at` AS `Post__created_at`,`Post`.`updated_at` AS `Post__updated_at`,`Post`.`deleted_at` AS `Post__deleted_at`,`Post`.`title` AS `Post__title`,`Post`.`content` AS `Post__content` FROM `comment` LEFT JOIN `post` `Post` ON `comment`.`post_id` = `Post`.`id` AND `Post`.`deleted_at` IS NULL WHERE `comment`.`deleted_at` IS NULL
	result = db.Joins("Post").Find(&comments2)
	if err := result.Error; err != nil {
		return err
	}
	for i, comment := range comments2 {
		fmt.Println(i, comment)
		fmt.Println(i, comment.Post)
	}
	return nil
}
