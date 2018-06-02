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
	Create(user *Gallery) error
}

func (gg *galleryGorm) Create(gallery *Gallery) error {
	return nil
}
