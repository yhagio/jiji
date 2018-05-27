package models

import (
	"errors"
	"jiji/utils"

	"golang.org/x/crypto/bcrypt"
)

// userValidator is our validation layer that validates
// and normalizes data before passing it on to the next
// UserDB in our interface chain.
type userValidator struct {
	UserDB
	hmac utils.HMAC
}

func (uv *userValidator) GetById(id uint) (*User, error) {
	if id <= 0 {
		return nil, errors.New("Invalid ID")
	}
	return uv.UserDB.GetById(id)
}

func (uv *userValidator) GetByToken(token string) (*User, error) {
	tokenHash := uv.hmac.Hash(token)
	return uv.UserDB.GetByToken(tokenHash)
}

func (uv *userValidator) Create(user *User) error {
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
	user.TokenHash = uv.hmac.Hash(user.Token)
	return uv.UserDB.Create(user)
}

func (uv *userValidator) Update(user *User) error {
	if user.Token != "" {
		user.TokenHash = uv.hmac.Hash(user.Token)
	}
	return uv.UserDB.Update(user)
}

func (uv *userValidator) Delete(id uint) error {
	if id == 0 {
		return ErrInvalidID
	}
	return uv.UserDB.Delete(id)
}
