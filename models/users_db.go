package models

import (
	"github.com/jinzhu/gorm"
)

// userGorm represents our database interaction layer
// and implements the UserDB interface fully.
type userGorm struct {
	db *gorm.DB
}

var _ UserDB = &userGorm{}

type UserDB interface {
	// Reader
	GetById(id uint) (*User, error)
	GetByEmail(email string) (*User, error)
	GetByToken(token string) (*User, error)
	// Writer
	Create(user *User) error
	Update(user *User) error
	Delete(id uint) error

	// Close DB connection
	Close() error
	// Migration tasks
	AutoMigrate() error
	DestructiveReset() error
}

// Get an user by id
func (ug *userGorm) GetById(id uint) (*User, error) {
	var user User
	db := ug.db.Where("id = ?", id)
	err := First(db, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Get an user by email
func (ug *userGorm) GetByEmail(email string) (*User, error) {
	var user User
	db := ug.db.Where("email = ?", email)
	err := First(db, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Get an user by token
func (ug *userGorm) GetByToken(tokenHash string) (*User, error) {
	var user User
	db := ug.db.Where("tokenHash = ?", tokenHash)
	err := First(db, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Create an user
func (ug *userGorm) Create(user *User) error {
	return ug.db.Create(user).Error
}

// Update an user
func (ug *userGorm) Update(user *User) error {
	return ug.db.Save(user).Error
}

// Delete an user
func (ug *userGorm) Delete(id uint) error {
	user := &User{
		Model: gorm.Model{ID: id},
	}
	return ug.db.Delete(user).Error
}

func (ug *userGorm) Close() error {
	return ug.db.Close()
}

// For development, testing only
// Recreate user table
func (ug *userGorm) DestructiveReset() error {
	err := ug.db.DropTableIfExists(&User{}).Error
	if err != nil {
		return err
	}
	return ug.AutoMigrate()
}

// Auto-migrate user table
func (ug *userGorm) AutoMigrate() error {
	err := ug.db.AutoMigrate(&User{}).Error
	if err != nil {
		return err
	}
	return nil
}
