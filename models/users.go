package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
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

func (us *UserService) Close() error {
	return us.db.Close()
}

// For development, testing only
// Recreate user table
func (us *UserService) DestructiveReset() error {
	err := us.db.DropTableIfExists(&User{}).Error
	if err != nil {
		return err
	}
	return us.AutoMigrate()
}

// Auto-migrate user table
func (us *UserService) AutoMigrate() error {
	err := us.db.AutoMigrate(&User{}).Error
	if err != nil {
		return err
	}
	return nil
}

// Get an user by id
func (us *UserService) GetById(id uint) (*User, error) {
	var user User
	db := us.db.Where("id = ?", id)
	err := First(db, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Get an user by email
func (us *UserService) GetByEmail(email string) (*User, error) {
	var user User
	db := us.db.Where("email = ?", email)
	err := First(db, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Create an user
func (us *UserService) Create(user *User) error {
	return us.db.Create(user).Error
}

// Update an user
func (us *UserService) Update(user *User) error {
	return us.db.Save(user).Error
}

// Delete an user
func (us *UserService) Delete(id uint) error {
	if id == 0 {
		return ErrInvalidID
	}
	user := &User{
		Model: gorm.Model{ID: id},
	}
	return us.db.Delete(user).Error
}
