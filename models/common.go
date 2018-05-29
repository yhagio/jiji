package models

import (
	"errors"

	"github.com/jinzhu/gorm"
)

var (
	ErrNotFound               = errors.New("models: resource not found")
	ErrInvalidID              = errors.New("models: ID provided was invalid")
	ErrInvalidEmailOrPassword = errors.New("model: Email or Password is incorrect")
	ErrEmailRequired          = errors.New("models: email is required")
	ErrEmailInvalid           = errors.New("models: email is not valid")
	ErrEmailTaken             = errors.New("models: email is already taken")
	ErrPasswordRequired       = errors.New("models: password is required")
	ErrPasswordTooShort       = errors.New("models: password must be at least 8 characters long")
	ErrTokenRequired          = errors.New("models: Token is required")
	ErrTokenTooShort          = errors.New("models: Token must be at least 32 bytes")
)

func First(db *gorm.DB, dst interface{}) error {
	err := db.First(dst).Error
	if err == gorm.ErrRecordNotFound {
		return ErrNotFound
	}
	return err
}
