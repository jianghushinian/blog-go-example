package store

import "gorm.io/gorm"

type UserStore interface {
	Create(user *User) error
	Get(id int) (*User, error)
}

func NewUserStore(db *gorm.DB) UserStore {
	return &userStore{db}
}

type userStore struct {
	db *gorm.DB
}

func (s *userStore) Create(user *User) error {
	return s.db.Create(user).Error
}

func (s *userStore) Get(id int) (*User, error) {
	var user User
	err := s.db.First(&user, id).Error
	return &user, err
}
