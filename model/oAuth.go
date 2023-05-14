package model

import (
	uuid "github.com/satori/go.uuid"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

const (
	ErrProviderIsEmpty      publicError = "model: provider is empty"
	ErrProviderNotSupported publicError = "model: provider not supported"
)

type OAuth struct {
	Base
	UserID   uuid.UUID `gorm:"not null;uniqueIndex:user_provider_index"`
	Provider string    `gorm:"not null;uniqueIndex:user_provider_index"`
	oauth2.Token
}

type OAuthDB interface {
	Find(userID string, provider string) (*OAuth, error)
	Create(o *OAuth) error
	Delete(id string) error
}

type OAuthService interface {
	OAuthDB
}

// NewOAuthService will create OAuthService to deal with consisting
// of two layers
//
// OAuthValidator -> OAuthGorm
func NewOAuthService(db *gorm.DB) OAuthService {
	oAuthGorm := NewOAuthGorm(db)
	return NewOAuthValidator(oAuthGorm)
}

type oAuthValidator struct {
	OAuthDB
}

func NewOAuthValidator(nextLayer OAuthDB) *oAuthValidator {
	return &oAuthValidator{
		OAuthDB: nextLayer,
	}
}

type oAuthValidationFn func(oAuth *OAuth) error

func runOAuthValidationFns(oAuth *OAuth, fns ...oAuthValidationFn) error {
	for _, fn := range fns {
		if err := fn(oAuth); err != nil {
			return err
		}
	}
	return nil
}

func (ov *oAuthValidator) validateOAuthID(oAuth *OAuth) error {
	if oAuth.ID.String() == ZeroID {
		return ErrInvalidID
	}
	if _, err := uuid.FromString(oAuth.ID.String()); err != nil {
		return ErrInvalidID
	}
	return nil
}

func (ov *oAuthValidator) validateUserID(oAuth *OAuth) error {
	if oAuth.UserID.String() == ZeroID {
		return ErrInvalidID
	}
	if _, err := uuid.FromString(oAuth.UserID.String()); err != nil {
		return ErrInvalidID
	}
	return nil
}

func (ov *oAuthValidator) validateProvider(oAuth *OAuth) error {
	if oAuth.Provider == "" {
		return ErrProviderIsEmpty
	} else if oAuth.Provider != "dropbox" {
		return ErrProviderNotSupported
	}
	return nil
}

func (ov *oAuthValidator) Create(oAuth *OAuth) error {
	err := runOAuthValidationFns(oAuth,
		ov.validateProvider,
		ov.validateOAuthID,
		ov.validateUserID,
	)
	if err != nil {
		return err
	}
	return ov.OAuthDB.Create(oAuth)
}

func (ov *oAuthValidator) Find(userID, provider string) (*OAuth, error) {
	oAuth := new(OAuth)
	oAuth.UserID = uuid.FromStringOrNil(userID)
	oAuth.Provider = provider

	err := runOAuthValidationFns(oAuth,
		ov.validateProvider,
		ov.validateUserID,
	)
	if err != nil {
		return nil, err
	}
	return ov.OAuthDB.Find(userID, provider)
}

func (ov *oAuthValidator) Delete(id string) error {
	err := runOAuthValidationFns(&OAuth{
		Base: Base{
			ID: uuid.FromStringOrNil(id),
		},
	})
	if err != nil {
		return err
	}
	return ov.OAuthDB.Delete(id)
}

type oAuthGorm struct {
	db *gorm.DB
}

func NewOAuthGorm(db *gorm.DB) *oAuthGorm {
	return &oAuthGorm{
		db: db,
	}
}

var _ OAuthDB = (*oAuthGorm)(nil)

func (og *oAuthGorm) Create(o *OAuth) error {
	return og.db.Create(&o).Error
}

func (og *oAuthGorm) Find(userID, provider string) (*OAuth, error) {
	o := new(OAuth)
	query := og.db.Where(OAuth{
		UserID:   uuid.FromStringOrNil(userID),
		Provider: provider,
	})
	err := getRecord(query, &o)
	if err != nil {
		return nil, err
	}
	return o, err
}

func (og *oAuthGorm) Delete(id string) error {
	return og.db.Delete(&OAuth{
		Base: Base{
			ID: uuid.FromStringOrNil(id),
		},
	}).Error
}
