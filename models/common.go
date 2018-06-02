package models

import (
	"strings"

	"github.com/jinzhu/gorm"
)

const (
	ErrNotFound               modelError = "models: resource not found"
	ErrInvalidID              modelError = "models: ID provided was invalid"
	ErrInvalidEmailOrPassword modelError = "models: Email or password is incorrect"
	ErrEmailRequired          modelError = "models: email is required"
	ErrEmailInvalid           modelError = "models: email is not valid"
	ErrEmailTaken             modelError = "models: email is already taken"
	ErrPasswordRequired       modelError = "models: password is required"
	ErrPasswordTooShort       modelError = "models: password must be at least 8 characters long"
	ErrTokenRequired          modelError = "models: Token is required"
	ErrTokenTooShort          modelError = "models: Token must be at least 32 bytes"
	ErrNoUserWithEmail        string     = "No user with the email is found"
)

func First(db *gorm.DB, dst interface{}) error {
	err := db.First(dst).Error
	if err == gorm.ErrRecordNotFound {
		return ErrNotFound
	}
	return err
}

type modelError string

func (e modelError) Error() string {
	return string(e)
}

// Creates public facing error message
func (e modelError) Public() string {
	cleanedStr := strings.Replace(string(e), "models: ", "", 1)
	split := strings.Split(cleanedStr, " ")
	split[0] = strings.Title(split[0]) // Capitalize first letter
	return strings.Join(split, " ")
}
