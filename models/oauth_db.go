package models

import "github.com/jinzhu/gorm"

type oauthGorm struct {
	db *gorm.DB
}

var _ OAuthDB = &oauthGorm{}

func (og *oauthGorm) Find(userID uint, service string) (*OAuth, error) {
	var oauth OAuth
	db := og.db.Where("user_id = ?", userID).Where("service = ?", service)
	err := First(db, &oauth)
	if err != nil {
		return nil, err
	}
	return &oauth, nil
}

func (og *oauthGorm) Create(oauth *OAuth) error {
	return og.db.Create(oauth).Error
}

func (og *oauthGorm) Delete(id uint) error {
	oauth := &OAuth{Model: gorm.Model{ID: id}}
	// NOTE: we want to delete oauth record permanently, not soft delete here
	// so, use Unscoped() from gorm
	return og.db.Unscoped().Delete(oauth).Error
}
