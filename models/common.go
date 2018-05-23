package models

import (
	"errors"

	"github.com/jinzhu/gorm"
)

var (
	ErrNotFound  = errors.New("models: resource not found")
	ErrInvalidID = errors.New("models: ID provided was invalid")
)

func First(db *gorm.DB, dst interface{}) error {
	err := db.First(dst).Error
	if err == gorm.ErrRecordNotFound {
		return ErrNotFound
	}
	return err
}