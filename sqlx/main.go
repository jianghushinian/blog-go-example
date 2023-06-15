package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db := Connect()
	// defer db.Close()

	// 插入记录
	id, err := MustCreateUser(db)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("id:", id)

	// 查询记录
	users, err := QueryxUsers(db)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("users:", users)

	user, err := QueryRowxUser(db, 1)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("user:", user)

	user, err = GetUser(db, 1)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("user:", user)

	users, err = SelectUsers(db)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("users:", users)

	// 事务
	if err := MustTransaction(db); err != nil {
		log.Fatal(err)
	}

	// 使用预处理语句查询记录
	user, err = PreparexGetUser(db)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("user:", user)

	// 具名参数
	if err := NamedExec(db); err != nil {
		log.Fatal(err)
	}

	users, err = NamedQuery(db)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("users:", users)

	// In 查询
	users, err = SqlxIn(db, []int64{1, 2, 3})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("users:", users)

	// Unsafe
	user, err = Unsafe(db)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("user:", user)

	// MapScan
	mapRes, err := MapScan(db)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("mapRes:", mapRes)
	for k, v := range mapRes[0] {
		log.Printf("%s:\t%s\n", k, v)
	}

	// SliceScan
	sliceRes, err := SliceScan(db)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("sliceRes:", sliceRes)
	for _, v := range sliceRes[0] {
		log.Printf("%s\n", v)
	}

	// 控制字段名称映射
	user, err = MapperFuncUseToUpper(db)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("user:", user)

	user, err = MapperFuncUseJsonTag(db)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("user:", user)
}

func Connect() *sqlx.DB {
	var (
		db  *sqlx.DB
		err error
		dsn = "user:password@tcp(127.0.0.1:3306)/demo?charset=utf8mb4&parseTime=true&loc=Local"
	)

	// 1. 使用 sqlx.Open 连接数据库
	db, err = sqlx.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	// 2. 使用 sqlx.Open 变体方法 sqlx.MustOpen 连接数据库，如果出现错误直接 panic
	db = sqlx.MustOpen("mysql", dsn)

	// 3. 如果已经有了 *sql.DB 对象，则可以使用 sqlx.NewDb 连接数据库，得到 *sqlx.DB 对象
	sqlDB, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	db = sqlx.NewDb(sqlDB, "mysql")

	// 以上 3 中方式，需要调用 db.Ping() 方法，确保能够正常连接数据库
	ctx := context.Background()
	if err := db.PingContext(ctx); err != nil {
		log.Fatal(err)
	}

	// 4. 使用 sqlx.Connect 连接数据库，等价于 sqlx.Open + db.Ping
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	// 5. 使用 sqlx.Connect 变体方法 sqlx.MustConnect 连接数据库，如果出现错误直接 panic
	db = sqlx.MustConnect("mysql", dsn)

	// 连接 SQLite/PostgreSQL
	// db = sqlx.MustConnect("sqlite3", ":memory:")
	// db = sqlx.MustConnect("postgres", "user=postgres password=password dbname=postgres host=localhost port=5432 sslmode=disable")
	return db
}

func MustCreateUser(db *sqlx.DB) (int64, error) {
	birthday := time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local)
	user := User{
		Name:     sql.NullString{String: "jianghushinian", Valid: true},
		Email:    "jianghushinian007@outlook.com",
		Age:      10,
		Birthday: birthday,
		Salary: Salary{
			Month: 100000,
			Year:  10000000,
		},
	}

	res := db.MustExec(
		`INSERT INTO user(name, email, age, birthday, salary) VALUES(?, ?, ?, ?, ?)`,
		user.Name, user.Email, user.Age, user.Birthday, user.Salary,
	)
	return res.LastInsertId()
}

func QueryxUsers(db *sqlx.DB) ([]User, error) {
	var us []User
	rows, err := db.Queryx("SELECT * FROM user")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var u User
		// sqlx 提供了便捷方法可以将查询结果直接扫描到结构体
		err = rows.StructScan(&u)
		if err != nil {
			return nil, err
		}
		us = append(us, u)
	}
	return us, nil
}

func QueryRowxUser(db *sqlx.DB, id int) (User, error) {
	var u User
	err := db.QueryRowx("SELECT * FROM user WHERE id = ?", id).StructScan(&u)
	return u, err
}

func GetUser(db *sqlx.DB, id int) (User, error) {
	var u User
	// 查询记录扫描数据到 struct
	err := db.Get(&u, "SELECT * FROM user WHERE id = ?", id)
	return u, err
}

func SelectUsers(db *sqlx.DB) ([]User, error) {
	var us []User
	// 查询记录扫描数据到 slice
	err := db.Select(&us, "SELECT * FROM user")
	return us, err
}

func MustTransaction(db *sqlx.DB) error {
	tx := db.MustBegin()
	tx.MustExec("UPDATE user SET age = 25 WHERE id = ?", 1)
	return tx.Commit()
}

func Transaction(db *sqlx.DB, id int64, name string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	res, err := tx.Exec("UPDATE user SET name = ? WHERE id = ?", name, id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	fmt.Printf("rowsAffected: %d\n", rowsAffected)

	return tx.Commit()
}

func PreparexGetUser(db *sqlx.DB) (User, error) {
	stmt, err := db.Preparex(`SELECT * FROM user WHERE id = ?`)
	if err != nil {
		return User{}, err
	}

	var u User
	err = stmt.Get(&u, 1)
	return u, err
}

func SqlxIn(db *sqlx.DB, ids []int64) ([]User, error) {
	query, args, err := sqlx.In("SELECT * FROM user WHERE id IN (?)", ids)
	if err != nil {
		return nil, err
	}

	// 将 SQL 查询中的占位符`?` 转换为适合特定数据库的占位符，如 SQLite 中的 `?`，MySQL 中的 `?` 或 PostgreSQL 中的 `$1`
	// e.g., SELECT * FROM user WHERE id IN (?, ?, ?) => SELECT * FROM user WHERE id IN ($1, $2, $3)
	query = db.Rebind(query)
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var us []User
	for rows.Next() {
		var user User
		err = rows.Scan(&user.ID, &user.Name, &user.Email, &user.Age,
			&user.Birthday, &user.Salary, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		us = append(us, user)
	}
	return us, nil
}

func NamedExec(db *sqlx.DB) error {
	m := map[string]interface{}{
		"email": "jianghushinian007@outlook.com",
		"age":   18,
	}
	result, err := db.NamedExec(`UPDATE user SET age = :age WHERE email = :email`, m)
	if err != nil {
		return err
	}
	fmt.Println(result.RowsAffected())
	return nil
}

func NamedQuery(db *sqlx.DB) ([]User, error) {
	u := User{
		Email: "jianghushinian007@outlook.com",
		Age:   18,
	}
	rows, err := db.NamedQuery("SELECT * FROM user WHERE email = :email OR age = :age", u)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.StructScan(&user)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func Unsafe(db *sqlx.DB) (User, error) {
	var user struct {
		ID    int
		Name  string
		Email string
		// 没有 Age 属性
	}
	// 不使用 Unsafe 则报错：missing destination name age in *struct { ID int; Name string; Email string }
	udb := db.Unsafe()
	err := udb.Get(&user, "SELECT id, name, email, age FROM user WHERE id = ?", 1)
	if err != nil {
		return User{}, err
	}
	return User{
		ID:    user.ID,
		Name:  sql.NullString{String: user.Name},
		Email: user.Email,
	}, nil
}

func MapScan(db *sqlx.DB) ([]map[string]interface{}, error) {
	rows, err := db.Queryx("SELECT * FROM user")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []map[string]interface{}
	for rows.Next() {
		r := make(map[string]interface{})
		err := rows.MapScan(r)
		if err != nil {
			return nil, err
		}
		res = append(res, r)
	}
	return res, err
}

func SliceScan(db *sqlx.DB) ([][]interface{}, error) {
	rows, err := db.Queryx("SELECT * FROM user")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res [][]interface{}
	for rows.Next() {
		// cols is an []interface{} of all the column results
		cols, err := rows.SliceScan()
		if err != nil {
			return nil, err
		}
		res = append(res, cols)
	}
	return res, err
}

func MapperFuncUseToUpper(db *sqlx.DB) (User, error) {
	copyDB := sqlx.NewDb(db.DB, db.DriverName())
	copyDB.MapperFunc(strings.ToUpper)

	var user User
	err := copyDB.Get(&user, "SELECT id as ID, name as NAME, email as EMAIL FROM user WHERE id = ?", 1)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func MapperFuncUseJsonTag(db *sqlx.DB) (User, error) {
	copyDB := sqlx.NewDb(db.DB, db.DriverName())
	// Create a new mapper which will use the struct field tag "json" instead of "db"
	copyDB.Mapper = reflectx.NewMapperFunc("json", strings.ToLower)

	var user User
	// json tag
	err := copyDB.Get(&user, "SELECT id, name as username, email FROM user WHERE id = ?", 1)
	if err != nil {
		return User{}, err
	}
	return user, nil
}
