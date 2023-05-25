package main

import (
	"database/sql"
	"fmt"

	"gorm.io/gorm"
)

type UserResult struct {
	ID   int
	Name string
	Age  int
}

func Raw(db *gorm.DB) error {
	// 原生查询 SQL 和 Scan
	var userRes UserResult
	db.Raw(`SELECT id, name, age FROM user WHERE id = ?`, 3).Scan(&userRes)
	fmt.Printf("affected rows: %d\n", db.RowsAffected)
	if err := db.Error; err != nil {
		return err
	}
	fmt.Println(userRes)

	var sumage int
	db.Raw(`SELECT SUM(age) as sumage FROM user WHERE member_number ?`, gorm.Expr("IS NULL")).Scan(&sumage)
	fmt.Printf("affected rows: %d\n", db.RowsAffected)
	if err := db.Error; err != nil {
		return err
	}
	fmt.Println(sumage)

	// Exec 原生 SQL
	// db.Exec("DROP TABLE user")
	db.Exec("UPDATE user SET age = ? WHERE id IN ?", 18, []int64{1, 2})
	// 使用表达式
	db.Exec(`UPDATE user SET age = ? WHERE name = ?`, gorm.Expr("age * ? + ?", 1, 2), "Jianghu")

	// 命名参数
	var post Post
	db.Where("title LIKE @name OR content LiKE @name", sql.Named("name", "%Hello%")).Find(&post)

	// DryRun 模式
	var user User
	stmt := db.Session(&gorm.Session{DryRun: true}).First(&user, 1).Statement
	fmt.Println(stmt.SQL.String()) // SQL: SELECT * FROM `user` WHERE `user`.`id` = ? AND `user`.`deleted_at` IS NULL ORDER BY `user`.`id` LIMIT 1
	fmt.Println(stmt.Vars)         // 参数: [1]
	return nil
}
