package main

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

func CreateUser(db *gorm.DB) error {
	// 创建单条记录
	now := time.Now()
	email := "u1@jianghushinian.com"
	user := User{Name: "user1", Email: &email, Age: 18, Birthday: &now}
	// INSERT INTO `user` (`created_at`,`updated_at`,`deleted_at`,`name`,`email`,`age`,`birthday`,`member_number`,`activated_at`) VALUES ('2023-05-22 22:14:47.814','2023-05-22 22:14:47.814',NULL,'user1','u1@jianghushinian.com',18,'2023-05-22 22:14:47.812',NULL,NULL)
	result := db.Create(&user)      // 通过数据的指针来创建
	fmt.Printf("user: %+v\n", user) // user.ID 自动填充
	fmt.Printf("affected rows: %d\n", result.RowsAffected)
	fmt.Printf("error: %v\n", result.Error)
	if err := result.Error; err != nil {
		return err
	}

	// 批量插入
	now = time.Now()
	email2 := "u2@jianghushinian.com"
	email3 := "u3@jianghushinian.com"
	users := []User{
		{Name: "user2", Email: &email2, Age: 19, Birthday: &now},
		{Name: "user3", Email: &email3, Age: 20, Birthday: &now},
	}
	// INSERT INTO `user` (`created_at`,`updated_at`,`deleted_at`,`name`,`email`,`age`,`birthday`,`member_number`,`activated_at`) VALUES ('2023-05-22 22:14:47.834','2023-05-22 22:14:47.834',NULL,'user2','u2@jianghushinian.com',19,'2023-05-22 22:14:47.833',NULL,NULL),('2023-05-22 22:14:47.834','2023-05-22 22:14:47.834',NULL,'user3','u3@jianghushinian.com',20,'2023-05-22 22:14:47.833',NULL,NULL)
	result = db.Create(&users)
	fmt.Printf("affected rows: %d\n", result.RowsAffected)
	return result.Error
}

func CreatePost(db *gorm.DB) error {
	// 创建 post 时会自动创建 comment 和 tag，并开启事务
	var post Post
	post = Post{
		Title:   "post1",
		Content: "content1",
		Comments: []*Comment{
			{Content: "comment1", Post: &post},
			{Content: "comment2", Post: &post},
		},
		Tags: []*Tag{
			{Name: "tag1"},
			{Name: "tag2"},
		},
	}
	// BEGIN TRANSACTION;
	// INSERT INTO `tag` (`created_at`,`updated_at`,`deleted_at`,`name`) VALUES ('2023-05-22 22:56:52.923','2023-05-22 22:56:52.923',NULL,'tag1'),('2023-05-22 22:56:52.923','2023-05-22 22:56:52.923',NULL,'tag2') ON DUPLICATE KEY UPDATE `id`=`id`
	// INSERT INTO `post` (`created_at`,`updated_at`,`deleted_at`,`title`,`content`) VALUES ('2023-05-22 22:56:52.898','2023-05-22 22:56:52.898',NULL,'post1','content1') ON DUPLICATE KEY UPDATE `id`=`id`
	// INSERT INTO `comment` (`created_at`,`updated_at`,`deleted_at`,`content`,`post_id`) VALUES ('2023-05-22 22:56:52.942','2023-05-22 22:56:52.942',NULL,'comment1',1),('2023-05-22 22:56:52.942','2023-05-22 22:56:52.942',NULL,'comment2',1) ON DUPLICATE KEY UPDATE `post_id`=VALUES(`post_id`)
	// INSERT INTO `post_tags` (`post_id`,`tag_id`) VALUES (1,1),(1,2) ON DUPLICATE KEY UPDATE `post_id`=`post_id`
	// COMMIT;
	result := db.Create(&post)
	fmt.Printf("affected rows: %d\n", result.RowsAffected)
	if err := result.Error; err != nil {
		return err
	}
	fmt.Println(post)

	var post2 Post
	post2 = Post{
		Title:   "post2",
		Content: "content2",
		Comments: []*Comment{
			{Content: "comment22", Post: &post2},
		},
	}
	// BEGIN TRANSACTION;
	// INSERT INTO `post` (`created_at`,`updated_at`,`deleted_at`,`title`,`content`) VALUES ('2023-05-22 22:56:52.955','2023-05-22 22:56:52.955',NULL,'post2','content2')
	// INSERT INTO `comment` (`created_at`,`updated_at`,`deleted_at`,`content`,`post_id`) VALUES ('2023-05-22 22:56:52.958','2023-05-22 22:56:52.958',NULL,'comment22',2) ON DUPLICATE KEY UPDATE `post_id`=VALUES(`post_id`)
	// COMMIT;
	result = db.Create(&post2)
	fmt.Printf("affected rows: %d\n", result.RowsAffected)
	if err := result.Error; err != nil {
		return err
	}
	fmt.Println(post2)

	// 使用 db.Save() 创建记录
	var post3 Post
	post3 = Post{
		Title:   "post3",
		Content: "content3",
		Comments: []*Comment{
			{Content: "comment33", Post: &post3},
		},
		Tags: []*Tag{
			{Name: "tag3"},
		},
	}
	// BEGIN TRANSACTION;
	// INSERT INTO `tag` (`created_at`,`updated_at`,`deleted_at`,`name`) VALUES ('2023-05-22 23:17:53.189','2023-05-22 23:17:53.189',NULL,'tag3') ON DUPLICATE KEY UPDATE `id`=`id`
	// INSERT INTO `post` (`created_at`,`updated_at`,`deleted_at`,`title`,`content`) VALUES ('2023-05-22 23:17:53.189','2023-05-22 23:17:53.189',NULL,'post3','content3') ON DUPLICATE KEY UPDATE `id`=`id`
	// INSERT INTO `comment` (`created_at`,`updated_at`,`deleted_at`,`content`,`post_id`) VALUES ('2023-05-22 23:17:53.19','2023-05-22 23:17:53.19',NULL,'comment33',0) ON DUPLICATE KEY UPDATE `post_id`=VALUES(`post_id`)
	// INSERT INTO `post_tags` (`post_id`,`tag_id`) VALUES (0,0) ON DUPLICATE KEY UPDATE `post_id`=`post_id`
	// COMMIT;
	result = db.Save(&post3)
	fmt.Printf("affected rows: %d\n", result.RowsAffected)
	if err := result.Error; err != nil {
		return err
	}
	fmt.Println(post3)
	return nil
}
