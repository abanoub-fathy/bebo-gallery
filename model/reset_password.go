package model

import (
	"github.com/abanoub-fathy/bebo-gallery/pkg/hash"
	"github.com/abanoub-fathy/bebo-gallery/pkg/rand"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type pwReset struct {
	Base
	UserID    uuid.UUID `gorm:"not null;unique;index"`
	Token     string    `gorm:"-"`
	TokenHash string    `gorm:"not null;"`
}

type pwResetDB interface {
	GetByToken(token string) (*pwReset, error)
	Create(p *pwReset) error
	Delete(id uuid.UUID) error
}

type pwResetValidator struct {
	pwResetDB
	hasher *hash.Hasher
}

type pwResetValidationFn func(p *pwReset) error

func runPwResetValidationFns(p *pwReset, fns ...pwResetValidationFn) error {
	for _, fn := range fns {
		if err := fn(p); err != nil {
			return err
		}
	}

	return nil
}

func newPwResetValidator(db pwResetDB, hasher *hash.Hasher) *pwResetValidator {
	return &pwResetValidator{
		pwResetDB: db,
		hasher:    hasher,
	}
}

func (pv *pwResetValidator) GetByToken(token string) (*pwReset, error) {
	p := &pwReset{
		Token: token,
	}
	if err := runPwResetValidationFns(p, pv.setTokenHash); err != nil {
		return nil, err
	}
	return pv.pwResetDB.GetByToken(p.TokenHash)
}

func (pv *pwResetValidator) Create(p *pwReset) error {
	err := runPwResetValidationFns(p,
		pv.requireUserID,
		pv.setToken,
		pv.setTokenHash,
	)
	if err != nil {
		return err
	}

	return pv.pwResetDB.Create(p)
}

func (pv *pwResetValidator) Delete(id uuid.UUID) error {
	p := &pwReset{
		Base: Base{
			ID: id,
		},
	}
	if err := runPwResetValidationFns(p, pv.validateID); err != nil {
		return err
	}
	return pv.pwResetDB.Delete(id)
}

func (pv *pwResetValidator) requireUserID(p *pwReset) error {
	if p.UserID.String() == ZeroID {
		return ErrUserIDRequired
	}

	return nil
}

func (pv *pwResetValidator) setToken(p *pwReset) error {
	token, err := rand.GenerateRememberToken()
	if err != nil {
		return err
	}
	p.Token = token
	return nil
}

func (pv *pwResetValidator) setTokenHash(p *pwReset) error {
	p.TokenHash = pv.hasher.HashByHMAC(p.Token)
	return nil
}

func (pv *pwResetValidator) validateID(p *pwReset) error {
	if p.ID.String() == ZeroID {
		return ErrInvalidID
	}
	return nil
}

type pwResetGorm struct {
	db *gorm.DB
}

// make sure that pwResetGorm implements pwResetDB
var _ pwResetDB = (*pwResetGorm)(nil)

func (pg *pwResetGorm) GetByToken(tokenHash string) (*pwReset, error) {
	p := new(pwReset)
	query := pg.db.Where(pwReset{
		TokenHash: tokenHash,
	})
	err := getRecord(query, &p)
	if err != nil {
		return nil, err
	}
	return p, err
}

func (pg *pwResetGorm) Create(p *pwReset) error {
	return pg.db.Create(&p).Error
}

func (pg *pwResetGorm) Delete(id uuid.UUID) error {
	return pg.db.Delete(pwReset{
		Base: Base{
			ID: id,
		},
	}).Error
}
