package model

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Service type contains all the types of the
// services used in our application
type Service struct {
	db *gorm.DB
	GalleryService
	UserService
	ImageService
}

// NewService is used to create service struct
// that will used to wrap all the services
// we used in our application inside it
func NewService(DB_URI string) (*Service, error) {
	// open db connection to be used in all services
	db, err := gorm.Open(postgres.Open(DB_URI), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	service := &Service{
		db:             db,
		GalleryService: NewGalleryService(db),
		UserService:    NewUserService(db),
		ImageService:   NewImageService(),
	}

	return service, nil
}

// Close should be used to close the db connection
func (s *Service) Close() error {
	sqlDB, err := s.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// ResetDB is used to reset all the tables inside
// the DB it will Drop all the the tables Created
// and it will also re-build them again
//
// if you want to Drop the tables and then create
// new fresh tables with no data inside them
// then call this method
func (s *Service) ResetDB() error {
	if err := s.db.Migrator().DropTable(&User{}, &Gallery{}, &pwReset{}); err != nil {
		return err
	}
	return s.AutoMigrate()
}

// AutoMigrate should be used to auto migrate
// all models to the database
func (s *Service) AutoMigrate() error {
	return s.db.AutoMigrate(&User{}, &Gallery{}, &pwReset{})
}
