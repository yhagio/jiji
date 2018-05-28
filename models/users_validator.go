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
	user := User{Token: token}
	err := userValidationFuncs(&user, uv.hmacHashToken)
	if err != nil {
		return nil, err
	}
	return uv.UserDB.GetByToken(user.TokenHash)
}

func (uv *userValidator) Create(user *User) error {
	err := userValidationFuncs(user,
		uv.generatePasswordHash,
		uv.setTokenIfNotSet,
		uv.hmacHashToken)
	if err != nil {
		return err
	}

	return uv.UserDB.Create(user)
}

func (uv *userValidator) Update(user *User) error {
	err := userValidationFuncs(user,
		uv.generatePasswordHash,
		uv.hmacHashToken)
	if err != nil {
		return err
	}

	return uv.UserDB.Update(user)
}

func (uv *userValidator) Delete(id uint) error {
	var user User
	user.ID = id
	err := userValidationFuncs(&user, uv.idGreaterThan(0))
	if err != nil {
		return err
	}

	return uv.UserDB.Delete(id)
}

func (uv *userValidator) generatePasswordHash(user *User) error {
	// If password is not changed, do nothing
	if user.Password == "" {
		return nil
	}

	hasedBytes, err := bcrypt.GenerateFromPassword(
		[]byte(user.Password+userPwPepper),
		bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	user.PasswordHash = string(hasedBytes)
	user.Password = ""
	return nil
}

func (uv *userValidator) setTokenIfNotSet(user *User) error {
	if user.Token != "" {
		return nil
	}

	token, err := utils.GenerateToken()
	if err != nil {
		return err
	}
	user.Token = token

	return nil
}

func (uv *userValidator) hmacHashToken(user *User) error {
	if user.Token == "" {
		return nil
	}
	user.TokenHash = uv.hmac.Hash(user.Token)
	return nil
}

// Closure way, dynamically validates with argument
func (uv *userValidator) idGreaterThan(num uint) userValidationFunc {
	return userValidationFunc(func(user *User) error {
		if user.ID <= num {
			return ErrInvalidID
		}
		return nil
	})
}

// Reusable validation functions runner / helper

type userValidationFunc func(*User) error

func userValidationFuncs(user *User, funcs ...userValidationFunc) error {
	for _, fn := range funcs {
		err := fn(user)
		if err != nil {
			return err
		}
	}
	return nil
}
