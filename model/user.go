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

// FindByID is used to find user by its id
// it will return the user from db and error if there is an error
// if there is no user found it will return error of type ErrNotFound
func (userService *UserService) FindByID(ID string) (*User, error) {
	// define the user
	user := new(User)

	// fetch user by id from db
	query := userService.db.Where(User{
		Base: Base{
			ID: uuid.FromStringOrNil(ID),
		},
	})

	// get user record
	err := getRecord(query, &user)

	// return result
	return user, err
}

// FindByEmail is used to find user by its email address
// it will return the user from db and error if there is an error
// if there is no user found it will return error of type ErrNotFound
func (userService *UserService) FindByEmail(email string) (*User, error) {
	user := new(User)
	query := userService.db.Where(&User{
		Email: email,
	})
	err := getRecord(query, user)
	return user, err
}

func getRecord(query *gorm.DB, destination interface{}) error {
	switch err := query.First(destination).Error; err {
	case nil:
		return nil
	case gorm.ErrRecordNotFound:
		destination = nil
		return ErrNotFound
	default:
		destination = nil
		return err
	}
}

// FindAndDeleteByID is used to delete user by its id
//
// it will first find the user and then delete it
//
// if the user is not found it will return ErrNotFound
//
// it returns the deleted user and the error if existed
func (userService *UserService) FindAndDeleteByID(userID string) (*User, error) {
	user, err := userService.FindByID(userID)
	if err != nil {
		return nil, err
	}
	err = userService.db.Delete(&user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

// FindAndUpdateByID is used to update user by its id
// it will return the error if there is something wrong while updating user
//
// if the user is updated correctly it will return the updated user and nil error
//
// if there is no user found it will return error of type ErrNotFound
func (userService *UserService) FindAndUpdateByID(userID string, updates map[string]interface{}) (*User, error) {
	user, err := userService.FindByID(userID)
	if err != nil {
		return nil, err
	}
	err = userService.db.Model(&user).Updates(updates).Error
	if err != nil {
		return nil, err
	}
	return user, err
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

// AutoMigrate is used to auto migrate user table into the database
func (userService *UserService) AutoMigrate() error {
	return userService.db.AutoMigrate(&User{})
}

// ResetUserDB is used to drop user table and create new one
func (userService *UserService) ResetUserDB() error {
	if err := userService.db.Migrator().DropTable(&User{}); err != nil {
		return err
	}

	// auto migrate user
	return userService.AutoMigrate()
}
