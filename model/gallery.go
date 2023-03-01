package model

import (
	"strings"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

const (
	ErrUserIDRequired       publicError = "model: user id is required"
	ErrGalleryTitleRequired publicError = "model: gallery title is required"
	ErrInvalidID            publicError = "model: not valid id"

	ZeroID = "00000000-0000-0000-0000-000000000000"
)

// Gallery is the container for images we will add
type Gallery struct {
	Base
	Title  string    `gorm:"not_null"`
	UserID uuid.UUID `gorm:"not_null;index"`
}

// galleryValidationFunc is a type for gallery validation
// functions.
//
// these functions receives refernce to gallery and return error
type galleryValidationFn func(*Gallery) error

func runGalleryValidationFns(gallery *Gallery, fns ...galleryValidationFn) error {
	for _, fn := range fns {
		if err := fn(gallery); err != nil {
			return err
		}
	}
	return nil
}

type GalleryService interface {
	GalleryDB
}

// GalleryDB has all methods needed to implemnt and
// use the Gallery database methods
type GalleryDB interface {
	// CreateGallery is used to create a new gallery into the DB
	CreateGallery(gallery *Gallery) error

	// GetGalleryByID is used to get specific gallery by its id
	FindByID(ID string) (*Gallery, error)
}

type galleryService struct {
	GalleryDB
}

type galleryValidator struct {
	GalleryDB
}

func (gv *galleryValidator) validateGalleryUserID(g *Gallery) error {
	if g.UserID.String() == ZeroID {
		return ErrUserIDRequired
	}
	return nil
}

func (gv *galleryValidator) validateGalleryTitle(g *Gallery) error {
	if strings.TrimSpace(g.Title) == "" {
		return ErrGalleryTitleRequired
	}
	return nil
}

func (gv *galleryValidator) CreateGallery(gallery *Gallery) error {
	err := runGalleryValidationFns(gallery,
		gv.validateGalleryTitle,
		gv.validateGalleryUserID,
	)
	if err != nil {
		return err
	}

	return gv.GalleryDB.CreateGallery(gallery)
}

func (gv *galleryValidator) FindByID(ID string) (*Gallery, error) {
	parsedUUID := uuid.FromStringOrNil(ID)
	if parsedUUID.String() == ZeroID {
		return nil, ErrInvalidID
	}

	return gv.GalleryDB.FindByID(ID)
}

// NewGalleryService is used to return GalleryService
// with its layers first layer is the validator the second
// is the gorm layer
func NewGalleryService(db *gorm.DB) GalleryService {
	return &galleryService{
		GalleryDB: &galleryValidator{
			GalleryDB: &galleryGorm{
				db: db,
			},
		},
	}
}

// galleryGorm is the type that will implements the
// the GalleryDB for gorm
type galleryGorm struct {
	db *gorm.DB
}

// NewGalleryGorm is used to create a gallery gorm
// that implements the GalleryDB interface
func NewGalleryGorm(db *gorm.DB) *galleryGorm {
	return &galleryGorm{
		db: db,
	}
}

// making sure that galleryGorm implemnts the GalleryDB
var _ GalleryDB = (*galleryGorm)(nil)

func (gg *galleryGorm) CreateGallery(gallery *Gallery) error {
	return gg.db.Create(&gallery).Error
}

func (gg *galleryGorm) FindByID(ID string) (*Gallery, error) {
	gallery := new(Gallery)
	query := gg.db.Where(Gallery{
		Base: Base{
			ID: uuid.FromStringOrNil(ID),
		},
	})
	err := getRecord(query, &gallery)
	return gallery, err
}
