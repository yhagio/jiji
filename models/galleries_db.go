package models

import (
	"github.com/jinzhu/gorm"
)

// galleryGorm represents our database interaction layer
// and implements the UserDB interface fully.
type galleryGorm struct {
	db *gorm.DB
}

var _ GalleryDB = &galleryGorm{}

type GalleryDB interface {
	GetOneById(id uint) (*Gallery, error)
	Create(gallery *Gallery) error
}

func (gg *galleryGorm) Create(gallery *Gallery) error {
	return gg.db.Create(gallery).Error
}

func (gg *galleryGorm) GetOneById(id uint) (*Gallery, error) {
	var gallery Gallery
	db := gg.db.Where("id = ?", id)
	err := First(db, &gallery)
	if err != nil {
		return nil, err
	}
	return &gallery, nil
}
