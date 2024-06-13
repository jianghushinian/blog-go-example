package main

import (
	user "github.com/jianghushinian/blog-go-example/wire/user/internal"
	"github.com/jianghushinian/blog-go-example/wire/user/internal/config"
	"github.com/jianghushinian/blog-go-example/wire/user/pkg/db"
)

func main() {
	cfg := &config.Config{
		MySQL: db.MySQLOptions{
			Address:  "127.0.0.1:3306",
			Database: "user",
			Username: "root",
			Password: "123456",
		},
	}

	app, cleanup, err := user.NewApp(cfg)
	if err != nil {
		panic(err)
	}

	defer cleanup()
	app.Run()
}
