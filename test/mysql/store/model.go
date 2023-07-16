package store

type User struct {
	ID   int    `gorm:"id"`
	Name string `gorm:"name"`
}

func (u *User) TableName() string {
	return "user"
}
