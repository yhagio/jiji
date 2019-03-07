package models

import "github.com/jinzhu/gorm"

type Services struct {
	Gallery GalleryService
	User    UserService
	Image   ImageService
	OAuth   OAuthService
	db      *gorm.DB
}

type ServicesConfig func(*Services) error

func NewServices(configs ...ServicesConfig) (*Services, error) {
	var services Services
	// For each ServicesConfig function...
	for _, config := range configs {
		// Run the function passing in a pointer to our Services
		// object and catching any errors
		if err := config(&services); err != nil {
			return nil, err
		}
	}
	// Then finally return the result
	return &services, nil
}

func (services *Services) Close() error {
	return services.db.Close()
}

// For development, testing only
// Recreate tables
func (services *Services) DestructiveReset() error {
	err := services.db.DropTableIfExists(&User{}, &Gallery{}, &OAuth{}, &passwordReset{}).Error
	if err != nil {
		return err
	}
	return services.AutoMigrate()
}

// Auto-migrate tables
func (services *Services) AutoMigrate() error {
	err := services.db.AutoMigrate(&User{}, &Gallery{}, &OAuth{}, &passwordReset{}).Error
	if err != nil {
		return err
	}
	return nil
}

func WithGorm(dialect, connectionInfo string) ServicesConfig {
	return func(s *Services) error {
		db, err := gorm.Open(dialect, connectionInfo)
		if err != nil {
			return err
		}
		s.db = db
		return nil
	}
}

func WithLogMode(mode bool) ServicesConfig {
	return func(s *Services) error {
		s.db.LogMode(mode)
		return nil
	}
}

/////////////////////////
// Create each service //
/////////////////////////
func WithUser(pepper, hmacKey string) ServicesConfig {
	return func(s *Services) error {
		s.User = NewUserService(s.db, pepper, hmacKey)
		return nil
	}
}

func WithGallery() ServicesConfig {
	return func(s *Services) error {
		s.Gallery = NewGalleryService(s.db)
		return nil
	}
}

func WithImage() ServicesConfig {
	return func(s *Services) error {
		s.Image = NewImageService()
		return nil
	}
}

func WithOAuth() ServicesConfig {
	return func(s *Services) error {
		// TODO later (optional)
		// Error handling in case s.db is nil
		s.OAuth = NewOAuthService(s.db)
		return nil
	}
}
