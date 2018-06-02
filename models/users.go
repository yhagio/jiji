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

// UserService is a set of methods used to manipulate and
// work with the user model
type UserService interface {
	Authenticate(email, password string) (*User, error)
	UserDB
}

func NewUserService(db *gorm.DB) UserService {
	ug := &userGorm{db}
	hmac := utils.NewHMAC(hmacSecretKey)
	uv := newUserValidator(ug, hmac)

	return &userService{
		UserDB: uv,
	}
}

var _ UserService = &userService{}

type userService struct {
	UserDB
}

// Authenticate user. Checks email and password.
func (us *userService) Authenticate(email, password string) (*User, error) {
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
