package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Product 定义结构体用来映射数据库表
type Product struct {
	gorm.Model
	Code  string
	Price uint
}

func main() {
	// 建立数据库连接
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// 迁移表结构
	db.AutoMigrate(&Product{})

	// 增加数据
	db.Create(&Product{Code: "D42", Price: 100})

	// 查找数据
	var product Product
	db.First(&product, 1)                 // find product with integer primary key
	db.First(&product, "code = ?", "D42") // find product with code D42

	// 更新数据 - update product's price to 200
	db.Model(&product).Update("Price", 200)
	// 更新数据 - update multiple fields
	db.Model(&product).Updates(Product{Price: 200, Code: "F42"}) // non-zero fields
	db.Model(&product).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})

	// 删除数据 - delete product
	db.Delete(&product, 1)
}
