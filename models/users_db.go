package models

import (
	"jiji/utils"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

// userGorm represents our database interaction layer
// and implements the UserDB interface fully.
type userGorm struct {
	db   *gorm.DB
	hmac utils.HMAC
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
func (ug *userGorm) GetByToken(token string) (*User, error) {
	var user User
	tokenHash := ug.hmac.Hash(token)
	db := ug.db.Where("tokenHash = ?", tokenHash)
	err := First(db, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Create an user
func (ug *userGorm) Create(user *User) error {
	hasedBytes, err := bcrypt.GenerateFromPassword(
		[]byte(user.Password+userPwPepper),
		bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	user.PasswordHash = string(hasedBytes)
	user.Password = ""

	// Generate token, hash it, and save it
	if user.Token == "" {
		token, err := utils.GenerateToken()
		if err != nil {
			return err
		}
		user.Token = token
	}
	user.TokenHash = ug.hmac.Hash(user.Token)

	return ug.db.Create(user).Error
}

// Update an user
func (ug *userGorm) Update(user *User) error {
	if user.Token != "" {
		user.TokenHash = ug.hmac.Hash(user.Token)
	}
	return ug.db.Save(user).Error
}

// Delete an user
func (ug *userGorm) Delete(id uint) error {
	if id == 0 {
		return ErrInvalidID
	}
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

func NewUserGorm(connectionInfo string) (*userGorm, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	hmac := utils.NewHMAC(hmacSecretKey)
	return &userGorm{
		db:   db,
		hmac: hmac,
	}, nil
}
