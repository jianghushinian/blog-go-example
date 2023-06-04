package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// 建立连接
	db, err := sql.Open("mysql",
		"user:password@tcp(127.0.0.1:3306)/demo?charset=utf8mb4&parseTime=true&loc=Local")
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = db.Close() }()

	// 配置 sql.DB 连接池参数
	// ref: https://www.alexedwards.net/blog/configuring-sqldb
	db.SetMaxOpenConns(25)                 // 设置最大的并发连接数（in-use + idle）
	db.SetMaxIdleConns(25)                 // 设置最大的空闲连接数（idle）
	db.SetConnMaxLifetime(5 * time.Minute) // 设置连接的最大生命周期

	// 检查数据库连接
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	// db.PingContext(context.Background())

	// 插入记录
	if id, err := CreateUser(db); err != nil {
		log.Fatal(err)
	} else {
		log.Println("id:", id)
	}
	if ids, err := CreateUsers(db); err != nil {
		log.Fatal(err)
	} else {
		log.Println("ids:", ids)
	}

	// 查询多条记录
	if users, err := GetUsers(db); err != nil {
		log.Fatal(err)
	} else {
		for _, user := range users {
			log.Printf("user: %+v\n", user)
		}
	}

	// 查询单条记录
	id := int64(1)
	if user, err := GetUser(db, id); err != nil {
		log.Fatal(err)
	} else {
		log.Printf("user: %+v\n", user)
	}

	// 预处理查询
	user, err := GetUserByPreparedStatement(db, id)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("user: %+v\n", user)
	}

	// 更新记录
	id = int64(3)
	name := "jianghushinian007"
	if err := UpdateUserName(db, id, name); err != nil {
		log.Fatal(err)
	}

	// 删除记录
	id = int64(3)
	if err := DeleteUser(db, id); err != nil {
		log.Fatal(err)
	}

	// 事务
	id = int64(2)
	name = "jianghushinian"
	if err := Transaction(db, id, name); err != nil {
		log.Fatal(err)
	}

	// 处理 NULL
	id = int64(2)
	if user, err := HandleNull(db, id); err != nil {
		log.Fatal(err)
	} else {
		log.Printf("user: %+v\n", user)
	}

	// 处理未知列
	id = int64(2)
	if res, err := HandleUnknownColumns(db, id); err != nil {
		log.Fatal(err)
	} else {
		for k, v := range res {
			rv := reflect.ValueOf(v)
			switch t := reflect.TypeOf(v); t.Kind() {
			case reflect.Ptr:
				log.Printf("key: %d, value: %+v\n", k, rv.Elem())
			default:
				log.Printf("key: %d, value: %+v\n", k, v)
			}
		}
	}
}

func CreateUser(db *sql.DB) (int64, error) {
	birthday := time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local)
	user := User{
		Name:     sql.NullString{String: "jianghushinian007", Valid: true},
		Email:    "jianghushinian007@outlook.com",
		Age:      10,
		Birthday: &birthday,
		Salary: Salary{
			Month: 100000,
			Year:  10000000,
		},
	}
	res, err := db.Exec(`INSERT INTO user(name, email, age, birthday, salary) VALUES(?, ?, ?, ?, ?)`,
		user.Name, user.Email, user.Age, user.Birthday, user.Salary)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func CreateUsers(db *sql.DB) ([]int64, error) {
	// 预处理
	stmt, err := db.Prepare("INSERT INTO user(name, email, age, birthday, salary) VALUES(?, ?, ?, ?, ?)")
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	birthday := time.Date(2000, 2, 2, 0, 0, 0, 0, time.Local)
	users := []User{
		{
			Name:     sql.NullString{String: "", Valid: true},
			Email:    "jianghushinian007@gmail.com",
			Age:      20,
			Birthday: &birthday,
			Salary: Salary{
				Month: 200000,
				Year:  20000000,
			},
		},
		{
			Name:  sql.NullString{String: "", Valid: false},
			Email: "jianghushinian007@163.com",
			Age:   30,
		},
	}

	var ids []int64
	for _, user := range users {
		res, err := stmt.Exec(user.Name, user.Email, user.Age, user.Birthday, user.Salary)
		if err != nil {
			return nil, err
		}
		id, err := res.LastInsertId()
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func GetUser(db *sql.DB, id int64) (User, error) {
	var user User
	row := db.QueryRow("SELECT * FROM user WHERE id = ?", id)
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Age,
		&user.Birthday, &user.Salary, &user.CreatedAt, &user.UpdatedAt)
	switch {
	case err == sql.ErrNoRows:
		return user, fmt.Errorf("no user with id %d", id)
	case err != nil:
		return user, err
	}
	// 处理错误
	if err := row.Err(); err != nil {
		return user, err
	}
	return user, nil
}

func GetUsers(db *sql.DB) ([]User, error) {
	rows, err := db.Query("SELECT * FROM user;")
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Age,
			&user.Birthday, &user.Salary, &user.CreatedAt, &user.UpdatedAt); err != nil {
			log.Println(err.Error())
			continue
		}
		users = append(users, user)
	}
	// 处理错误
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func GetUserByPreparedStatement(db *sql.DB, id int64) (User, error) {
	var user User
	stmt, err := db.Prepare("SELECT * FROM user WHERE id = ?")
	if err != nil {
		return user, err
	}
	defer func() { _ = stmt.Close() }()

	err = stmt.QueryRow(id).Scan(&user.ID, &user.Name, &user.Email, &user.Age,
		&user.Birthday, &user.Salary, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func UpdateUserName(db *sql.DB, id int64, name string) error {
	ctx := context.Background()
	res, err := db.ExecContext(ctx, "UPDATE user SET name = ? WHERE id = ?", name, id)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		// 如果新的 name 等于原 name，也会执行到这里
		return fmt.Errorf("no user with id %d", id)
	}
	return nil
}

func DeleteUser(db *sql.DB, id int64) error {
	ctx := context.Background()
	res, err := db.ExecContext(ctx, "DELETE FROM user WHERE id = ?", id)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return fmt.Errorf("no user with id %d", id)
	}
	return nil
}

func Transaction(db *sql.DB, id int64, name string) error {
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}

	_, execErr := tx.ExecContext(ctx, "UPDATE user SET name = ? WHERE id = ?", name, id)
	if execErr != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Fatalf("update failed: %v, unable to rollback: %v\n", execErr, rollbackErr)
		}
		log.Fatalf("update failed: %v", execErr)
	}
	return tx.Commit()
}

func HandleNull(db *sql.DB, id int64) (User, error) {
	var user User
	// row := db.QueryRow("SELECT id, name, updated_at FROM user WHERE id = ?", id)
	// 以上代码，如果 updated_at 为 NULL 则得到报错：sql: Scan error on column index 2, name "updated_at": converting NULL to string is unsupported
	// 使用 COALESCE(updated_at, '') 将 NULL 转成 ''
	row := db.QueryRow("SELECT id, name, COALESCE(updated_at, '') FROM user WHERE id = ?", id)
	if err := row.Scan(&user.ID, &user.Name, &user.UpdatedAt); err != nil {
		return user, err
	}
	if err := row.Err(); err != nil {
		return user, err
	}
	return user, nil
}

func HandleUnknownColumns(db *sql.DB, id int64) ([]interface{}, error) {
	var res []interface{}
	rows, err := db.Query("SELECT * FROM user WHERE id = ?", id)
	if err != nil {
		return res, err
	}
	defer func() { _ = rows.Close() }()

	// 如果不知道列名称，可以使用 rows.Columns() 查找列名称列表
	cols, err := rows.Columns()
	if err != nil {
		return res, err
	}

	fmt.Printf("columns: %v\n", cols) // [id name email age birthday salary created_at updated_at]
	fmt.Printf("columns length: %d\n", len(cols))

	// 获取列类型信息
	types, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}
	for _, typ := range types {
		// id: &{name:id hasNullable:true hasLength:false hasPrecisionScale:false nullable:false length:0 databaseType:INT precision:0 scale:0 scanType:0x1045d68a0}
		fmt.Printf("%s: %+v\n", typ.Name(), typ)
	}

	res = []interface{}{
		new(int),            // id
		new(sql.NullString), // name
		new(string),         // email
		new(int),            // age
		new(time.Time),      // birthday
		new(Salary),         // salary
		new(time.Time),      // created_at
		// 如果不知道列类型，可以使用 sql.RawBytes，它实际上是 []byte 的别名
		new(sql.RawBytes), // updated_at
	}

	for rows.Next() {
		if err := rows.Scan(res...); err != nil {
			return res, err
		}
	}
	return res, rows.Err()
}
