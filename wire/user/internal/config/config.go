package config

import "github.com/jianghushinian/blog-go-example/wire/user/pkg/db"

type Config struct {
	MySQL db.MySQLOptions `json:"mysql" yaml:"mysql"`
}
