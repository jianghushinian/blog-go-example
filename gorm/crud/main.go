package main

func main() {
	// 创建连接
	db, err := ConnectMySQL("localhost", "3306", "root", "password", "gorm")
	if err != nil {
		panic(err)
	}
	// 连接池设置
	if err := SetConnect(db); err != nil {
		panic(err)
	}

	// 迁移数据
	if err := Migrate(db); err != nil {
		panic(err)
	}

	// 单表操作
	{
		// 创建记录
		if err := CreateUser(db); err != nil {
			panic(err)
		}

		// 更新记录
		if err := UpdateUser(db); err != nil {
			panic(err)
		}

		// 查询记录
		if err := ReadUser(db); err != nil {
			panic(err)
		}

		// 删除记录
		if err := DeleteUser(db); err != nil {
			panic(err)
		}
	}

	// 关联表操作
	{
		if err := CreatePost(db); err != nil {
			panic(err)
		}

		if err := ReadPost(db); err != nil {
			panic(err)
		}

		if err := UpdatePost(db); err != nil {
			panic(err)
		}

		if err := DeletePost(db); err != nil {
			panic(err)
		}
	}

	// 事务
	{
		if err := TransactionPost(db); err != nil {
			panic(err)
		}

		if err := TransactionPostWithManually(db); err != nil {
			panic(err)
		}
	}

	// 钩子

	// 原生 SQL
	{
		if err := Raw(db); err != nil {
			panic(err)
		}
	}

	// 调试
	{
		if err := DebugLogger(db); err != nil {
			panic(err)
		}

		if err := DebugDryRun(db); err != nil {
			panic(err)
		}
	}
}
