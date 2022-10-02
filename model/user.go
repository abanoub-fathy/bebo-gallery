package model

import (
	"errors"

	uuid "github.com/satori/go.uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	// ErrNotFound is returned when the resource can not be found from DB
	ErrNotFound = errors.New("model: resource not found")
)

type User struct {
	Base
	FirstName string `gorm:"not null"`
	LastName  string `gorm:"not null"`
	Email     string `gorm:"not null;unique;index"`
}

type UserService struct {
	db *gorm.DB
}

// NewUserService creates a new userService to
// interact with users
func NewUserService(DB_URI string) (*UserService, error) {
	// connect to DB
	db, err := gorm.Open(postgres.Open(DB_URI), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// return userService
	return &UserService{
		db: db,
	}, nil
}

// CreateUser is used to create new user in our database
func (userService *UserService) CreateUser(user *User) error {
	return userService.db.Create(&user).Error
}

func (userService *UserService) FindByID(ID string) (*User, error) {
	// define the user
	user := new(User)

	// fetch user by id from db
	err := userService.db.Where(User{
		Base: Base{
			ID: uuid.FromStringOrNil(ID),
		},
	}).First(&user).Error

	// switch the error type
	switch err {
	case nil:
		return user, nil
	case gorm.ErrRecordNotFound:
		// if user is not found we will return nil for the user and Not Found error
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// Close used to close userService database connection
func (userService *UserService) Close() error {
	sqlDB, err := userService.db.DB()
	if err != nil {
		return err
	}
	if err := sqlDB.Close(); err != nil {
		return err
	}
	return nil
}

// ResetUserDB is used to drop user table and create new one
func (userService *UserService) ResetUserDB() error {
	if err := userService.db.Migrator().DropTable(&User{}); err != nil {
		return err
	}
	if err := userService.db.AutoMigrate(&User{}); err != nil {
		return err
	}
	return nil
}
