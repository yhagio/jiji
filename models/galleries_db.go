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
	Update(gallery *Gallery) error
	Delete(id uint) error
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

func (gg *galleryGorm) Create(gallery *Gallery) error {
	return gg.db.Create(gallery).Error
}

func (gg *galleryGorm) Update(gallery *Gallery) error {
	return gg.db.Save(gallery).Error
}

func (gg *galleryGorm) Delete(id uint) error {
	gallery := &Gallery{Model: gorm.Model{ID: id}}
	return gg.db.Delete(gallery).Error
}
