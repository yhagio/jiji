package models

import "github.com/jinzhu/gorm"

type Services struct {
	Gallery GalleryService
	User    UserService
	db      *gorm.DB
}

func NewServices(connectionInfo string) (*Services, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)

	return &Services{
		User:    NewUserService(db),
		Gallery: &galleryGorm{},
		db:      db,
	}, nil
}

func (services *Services) Close() error {
	return services.db.Close()
}

// For development, testing only
// Recreate tables
func (services *Services) DestructiveReset() error {
	err := services.db.DropTableIfExists(&User{}, &Gallery{}).Error
	if err != nil {
		return err
	}
	return services.AutoMigrate()
}

// Auto-migrate tables
func (services *Services) AutoMigrate() error {
	err := services.db.AutoMigrate(&User{}, &Gallery{}).Error
	if err != nil {
		return err
	}
	return nil
}
