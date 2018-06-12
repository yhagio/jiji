package models

import (
	"github.com/jinzhu/gorm"
)

type Gallery struct {
	gorm.Model
	UserId uint     `gorm:"not_null; index`
	Title  string   `gorm:"not_null`
	Images []string `gorm:"-"`
}

type GalleryService interface {
	GalleryDB
}

func NewGalleryService(db *gorm.DB) GalleryService {
	return &galleryService{
		GalleryDB: &galleryValidator{
			GalleryDB: &galleryGorm{
				db: db,
			},
		},
	}
}

var _ GalleryService = &galleryService{}

type galleryService struct {
	GalleryDB
}
