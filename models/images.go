package models

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func NewImageService() ImageService {
	return &imageService{}
}

type imageService struct {
}

type ImageService interface {
	Create(galleryID uint, r io.Reader, filename string) error
	GetAllByGalleryID(galleryID uint) ([]string, error)
}

func (is *imageService) Create(galleryID uint, r io.Reader, filename string) error {
	imagePath, err := is.createImagePath(galleryID)
	if err != nil {
		return err
	}
	// Create a destination file
	destination, err := os.Create(filepath.Join(imagePath, filename))
	if err != nil {
		return err
	}

	// Copy the uploaded file data to thedestination file
	_, err = io.Copy(destination, r)
	if err != nil {
		return err
	}

	return nil
}

func (is *imageService) GetAllByGalleryID(galleryID uint) ([]string, error) {
	imagePath := is.imagePath(galleryID)
	strings, err := filepath.Glob(filepath.Join(imagePath, "*"))
	if err != nil {
		return nil, err
	}
	return strings, nil
}

///////////////////////////////////////////////////////////
// Private
///////////////////////////////////////////////////////////

func (is *imageService) createImagePath(galleryID uint) (string, error) {
	imagePath := is.imagePath(galleryID)
	err := os.MkdirAll(imagePath, 0755)
	if err != nil {
		return "", err
	}
	return imagePath, nil
}

func (is *imageService) imagePath(galleryID uint) string {
	return filepath.Join("images", "galleries", fmt.Sprintf("%v", galleryID))
}
