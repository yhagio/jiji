package models

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
)

type Image struct {
	GalleryID uint
	Filename  string
}

func (i *Image) Path() string {
	temp := url.URL{
		Path: "/" + i.RelativePath(),
	}
	return temp.String()
}

func (i *Image) RelativePath() string {
	galleryID := fmt.Sprintf("%v", i.GalleryID)
	return filepath.ToSlash(filepath.Join("images", "galleries", galleryID, i.Filename))
}

func NewImageService() ImageService {
	return &imageService{}
}

type imageService struct {
}

type ImageService interface {
	Create(galleryID uint, r io.Reader, filename string) error
	GetAllByGalleryID(galleryID uint) ([]Image, error)
	Delete(i *Image) error
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

func (is *imageService) GetAllByGalleryID(galleryID uint) ([]Image, error) {
	imagePath := is.imagePath(galleryID)
	strings, err := filepath.Glob(filepath.Join(imagePath, "*"))
	if err != nil {
		return nil, err
	}

	imageSlice := make([]Image, len(strings))

	// Add "/" to all image file paths
	for i, imageStr := range strings {
		imageSlice[i] = Image{
			GalleryID: galleryID,
			Filename:  filepath.Base(imageStr),
		}
	}
	return imageSlice, nil
}

func (is *imageService) Delete(i *Image) error {
	return os.Remove(i.RelativePath())
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
