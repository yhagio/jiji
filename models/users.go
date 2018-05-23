package models

import (
	"errors"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	ErrNotFound = errors.New("models: resource not found")
)

type User struct {
	gorm.Model
	Username string `gorm:"not null; unique_index"`
	Email    string `gorm:"not null; unique_index"`
}

type UserService struct {
	db *gorm.DB
}

func NewUserService(connectionInfo string) (*UserService, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	return &UserService{db: db}, nil
}

func (us *UserService) close() error {
	return us.db.Close()
}

// For development, testing only
func (us *UserService) DestructiveReset() {
	us.db.DropTableIfExists(&User{})
	us.db.AutoMigrate(&User{})
}

// Get an user by id
func (us *UserService) GetById(id uint) (*User, error) {
	var user User
	err := us.db.Where("id = ?", id).First(&user).Error

	switch err {
	case nil:
		return &user, nil
	case gorm.ErrRecordNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// Create an user
func (us *UserService) Create(user *User) error {
	return us.db.Create(user).Error
}
