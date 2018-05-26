package models

import (
	"jiji/utils"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"golang.org/x/crypto/bcrypt"
)

var userPwPepper = "super-secret-pepper-for-password"

const hmacSecretKey = "secret-hmac-key"

type User struct {
	gorm.Model
	Username     string `gorm:"not null; unique_index"`
	Email        string `gorm:"not null; unique_index"`
	Password     string `gorm:"-"`
	PasswordHash string `gorm:"not null`
	Token        string `gorm:"-"`
	TokenHash    string `gorm:"not null; unique_index"`
}

type UserService struct {
	db   *gorm.DB
	hmac utils.HMAC
}

func NewUserService(connectionInfo string) (*UserService, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)

	hmac := utils.NewHMAC(hmacSecretKey)
	return &UserService{
		db:   db,
		hmac: hmac,
	}, nil
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

// Authentication
func (us *UserService) Authenticate(email, password string) (*User, error) {
	user, err := us.GetByEmail(email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(password+userPwPepper),
	)

	if err != nil {
		return nil, ErrInvalidEmailOrPassword
	}

	return user, nil
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
	user.TokenHash = us.hmac.Hash(user.Token)

	return us.db.Create(user).Error
}

// Update an user
func (us *UserService) Update(user *User) error {
	if user.Token != "" {
		user.TokenHash = us.hmac.Hash(user.Token)
	}
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
