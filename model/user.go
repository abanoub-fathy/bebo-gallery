package model

import (
	"errors"
	"regexp"
	"strings"

	"github.com/abanoub-fathy/bebo-gallery/config"
	"github.com/abanoub-fathy/bebo-gallery/pkg/hash"
	"github.com/abanoub-fathy/bebo-gallery/pkg/rand"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User is a tyype represent our user model
// in the database it is used for user accounts.
type User struct {
	Base
	FirstName        string `gorm:"not null"`
	LastName         string `gorm:"not null"`
	Email            string `gorm:"not null;unique;index"`
	Password         string `gorm:"-"`
	PasswordHash     string `gorm:"not null"`
	RememberToken    string `gorm:"-"`
	RemeberTokenHash string `gorm:"unique;index"`
	Galleries        []Gallery
}

// UserDB is used to interact with the users database.
//
// For pretty much all single user queries:
// If the user is found, we will return a nil error
// If the user is not found, we will return ErrNotFound
// If there is another error, we will return an error with
// more information about what went wrong. This may not be
// an error generated by the models package.
//
// For single user queries, any error but ErrNotFound should
// probably result in a 500 error until we make "public"
// facing errors.
type UserDB interface {
	// Methods for querying for single users
	FindByID(ID string) (*User, error)
	FindByEmail(email string) (*User, error)
	FindUserByRememberToken(token string) (*User, error)

	// Methods for altering users
	CreateUser(user *User) error
	FindAndUpdateByID(userID string, updates map[string]interface{}) (*User, error)
	FindAndDeleteByID(userID string) (*User, error)
	Save(user *User) error
	SaveNewRemeberToken(user *User) error
}

// userGorm represents our database interaction layer
// and implements the UserDB interface fully.
type userGorm struct {
	db     *gorm.DB
	hasher *hash.Hasher
}

// newUserGorm creates a new userGorm
// that implements the the UserDB interface
func newUserGorm(db *gorm.DB) *userGorm {
	// return userGorm object
	return &userGorm{
		db:     db,
		hasher: hash.NewHasher(config.AppConfig.HashSecretKey),
	}
}

var _ UserDB = &userGorm{}

// UserService is an interface that contains
// smethods to interact with user model
type UserService interface {
	// AuthenticateUser is used to check the user email vs password
	// if it is correct you will get the user and nil error
	// otherwise you will get an error
	//
	// error can be ErrEmailNotValidFormat, ErrNotFound, ErrPasswordNotCorrect
	// or other generic error during authenticate user
	AuthenticateUser(email, password string) (*User, error)

	UserDB
}

// userService struct is an implementation for UserService
// interface type.
type userService struct {
	UserDB
}

var _ UserService = &userService{}

// NewUserService creates a new userService to
// interact with users
func NewUserService(db *gorm.DB) UserService {
	// create new userGorm
	userGorm := newUserGorm(db)

	// create userValidator
	userValidator := newUserValidator(userGorm, hash.NewHasher(config.AppConfig.HashSecretKey))

	// set the userGorm to UserDB in the UserService
	userService := &userService{
		UserDB: userValidator,
	}

	// return
	return userService
}

type userValidator struct {
	UserDB
	hasher     *hash.Hasher
	emailRegex *regexp.Regexp
}

func newUserValidator(userDB UserDB, hasher *hash.Hasher) *userValidator {
	return &userValidator{
		UserDB:     userDB,
		hasher:     hasher,
		emailRegex: regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,16}$`),
	}
}

var _ UserDB = &userValidator{}

type userValidationFunc func(*User) error

func runUserValidationFuncs(user *User, fns ...userValidationFunc) error {
	for _, fn := range fns {
		err := fn(user)
		if err != nil {
			return err
		}
	}

	return nil
}

// CreateUser is used to create new user in our database
//
// in this layer this method will validate email and password
//
// also the method will create a hashed password
// then pass the user to the next UserDB layer
func (uv *userValidator) CreateUser(user *User) error {
	err := runUserValidationFuncs(user,
		uv.NormalizeEmail,
		uv.ValidateEmail,
		uv.EmailIsNotTaken,
		uv.RequirePassword,
		uv.ValidatePassword(8),
		uv.HashUserPassword,
		uv.GenerateNewRemeberToken,
		uv.CheckRemeberTokenLength,
		uv.HashUserRememberToken,
		uv.RequireRemeberTokenHash,
	)

	if err != nil {
		return err
	}

	// call the next DB layer
	return uv.UserDB.CreateUser(user)
}

// ValidatePassword is used to check for valid user password
// if the password is less than minLength chars it will trow an error
func (uv *userValidator) ValidatePassword(minLength uint) userValidationFunc {
	return func(user *User) error {
		if len(user.Password) < int(minLength) {
			return ErrPasswordTooShort
		}
		return nil
	}
}

func (uv *userValidator) RequirePassword(user *User) error {
	if user.Password == "" {
		return ErrPasswordRequired
	}
	return nil
}

func (uv *userValidator) HashUserPassword(user *User) error {
	if user.Password == "" {
		return errors.New("user password is required")
	}

	// hash the user password
	password, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(password)
	user.Password = ""
	return nil
}

func (uv *userValidator) GenerateNewRemeberToken(user *User) error {
	// generate new random remember token
	token, err := rand.GenerateRememberToken()
	if err != nil {
		return err
	}

	user.RememberToken = token
	return nil
}

func (uv *userValidator) HashUserRememberToken(user *User) error {
	if user.RememberToken == "" {
		return errors.New("token not present to hash")
	}
	// hash the token
	hashedToken := uv.hasher.HashByHMAC(user.RememberToken)
	user.RemeberTokenHash = hashedToken
	return nil
}

func (uv *userValidator) CheckRemeberTokenLength(user *User) error {
	n, err := rand.NBytes(user.RememberToken)
	if err != nil {
		return err
	}
	if n < 32 {
		return ErrRememberTooShort
	}
	return nil
}

func (uv *userValidator) RequireRemeberTokenHash(user *User) error {
	if user.RemeberTokenHash == "" {
		return ErrRememberTokenHashRequired
	}
	return nil
}

// SaveNewRemeberToken is used to generate and set new remember token
// for the given user in the argument
//
// if there is no error the method will return nil error
func (uv *userValidator) SaveNewRemeberToken(user *User) error {
	// generate and set token
	runUserValidationFuncs(user,
		uv.GenerateNewRemeberToken,
		uv.CheckRemeberTokenLength,
		uv.HashUserRememberToken,
		uv.RequireRemeberTokenHash,
	)

	// return to the next UserDB layer
	return uv.UserDB.SaveNewRemeberToken(user)
}

// FindUserByRememberToken will hash the token and
// pass the hashed token to the next UserDB layer
func (uv *userValidator) FindUserByRememberToken(token string) (*User, error) {
	user := &User{RememberToken: token}
	err := runUserValidationFuncs(user, uv.HashUserRememberToken)
	if err != nil {
		return nil, err
	}

	// call the next UserDB layer
	return uv.UserDB.FindUserByRememberToken(user.RemeberTokenHash)
}

func (uv *userValidator) FindAndUpdateByID(userID string, updates map[string]interface{}) (*User, error) {
	user := &User{}
	if _, emailUpdate := updates["email"]; emailUpdate {
		err := runUserValidationFuncs(user, uv.NormalizeEmail, uv.ValidateEmail, uv.EmailIsNotTaken)
		if err != nil {
			return nil, err
		}
	}

	if _, passwordUpdate := updates["password"]; passwordUpdate {
		err := runUserValidationFuncs(user, uv.ValidatePassword(8), uv.HashUserPassword)
		if err != nil {
			return nil, err
		}

		updates["PasswordHash"] = user.PasswordHash
		delete(updates, "password")
	}

	return uv.UserDB.FindAndUpdateByID(userID, updates)
}

// ValidateEmail validate that the email address is correct
// first it is going to check if it is empty or not
//
// if the email is not valid because it is not in the
// email format it is going to ruturn ErrEmailNotValidFormat
func (uv *userValidator) ValidateEmail(user *User) error {
	validEmail := uv.emailRegex.MatchString(user.Email)
	if !validEmail {
		return ErrEmailNotValidFormat
	}
	return nil
}

// NormalizeEmail is used to trim the space in the email
// address and also convert all chars to lowercase
//
// it is usually used before ValidateEmail method
func (uv *userValidator) NormalizeEmail(user *User) error {
	user.Email = strings.TrimSpace(user.Email)
	user.Email = strings.ToLower(user.Email)
	return nil
}

// EmailIsNotTaken is used to check if the email address
// is not taken by other users
func (uv *userValidator) EmailIsNotTaken(user *User) error {
	// call FindByEmail
	existingUser, err := uv.FindByEmail(user.Email)

	// if the user not found
	if err == ErrNotFound {
		return nil
	}

	// return other errors
	if err != nil {
		return err
	}

	// if the user exists

	// check if another user try to use existed email
	if existingUser.ID.String() != user.ID.String() {
		return ErrEmailIsTaken
	}

	return nil
}

// FindByEmail validation method is used
// to validate and normalize email address
func (uv *userValidator) FindByEmail(email string) (*User, error) {
	user := &User{Email: email}
	err := runUserValidationFuncs(user, uv.NormalizeEmail, uv.ValidateEmail)
	if err != nil {
		return nil, err
	}
	return uv.UserDB.FindByEmail(user.Email)
}

// AuthenticateUser is used to return user by email and password
//
// if the user is found the method will return the user object and nil error
func (userService *userService) AuthenticateUser(email, password string) (*User, error) {
	// find user by email
	user, err := userService.FindByEmail(email)
	if err != nil {
		return nil, err
	}

	// compare user password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		switch err {
		case bcrypt.ErrMismatchedHashAndPassword:
			return nil, ErrPasswordNotCorrect
		default:
			return nil, err
		}
	}

	// return the user and nil error
	return user, nil
}

// CreateUser is used to save user in the DB
func (ug *userGorm) CreateUser(user *User) error {
	return ug.db.Create(&user).Error
}

// SaveNewRemeberToken is used to save new remeber token
// to user
func (ug *userGorm) SaveNewRemeberToken(user *User) error {
	return ug.Save(user)
}

// FindByID is used to find user by its id
// it will return the user from db and error if there is an error
// if there is no user found it will return error of type ErrNotFound
func (ug *userGorm) FindByID(ID string) (*User, error) {
	// define the user
	user := new(User)

	// fetch user by id from db
	query := ug.db.Where(User{
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
func (ug *userGorm) FindByEmail(email string) (*User, error) {
	user := new(User)
	query := ug.db.Where(&User{
		Email: email,
	})
	err := getRecord(query, user)
	return user, err
}

// FindAndDeleteByID is used to delete user by its id
//
// it will first find the user and then delete it
//
// if the user is not found it will return ErrNotFound
//
// it returns the deleted user and the error if existed
func (ug *userGorm) FindAndDeleteByID(userID string) (*User, error) {
	user, err := ug.FindByID(userID)
	if err != nil {
		return nil, err
	}
	err = ug.db.Delete(&user).Error
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
func (ug *userGorm) FindAndUpdateByID(userID string, updates map[string]interface{}) (*User, error) {
	user, err := ug.FindByID(userID)
	if err != nil {
		return nil, err
	}
	err = ug.db.Model(&user).Updates(updates).Error
	if err != nil {
		return nil, err
	}
	return user, err
}

// Save is used to update existing user
func (ug *userGorm) Save(user *User) error {
	return ug.db.Save(&user).Error
}

// FindUserByRememberToken is used to find user by remember token
//
// in this layer the method expect to receive the hashed token
func (ug *userGorm) FindUserByRememberToken(hashedToken string) (*User, error) {
	// define user
	user := new(User)

	// make query
	query := ug.db.Where(&User{
		RemeberTokenHash: hashedToken,
	})
	// get the user
	err := getRecord(query, user)

	// return user
	return user, err
}
