package model

import (
	"strings"

	"github.com/abanoub-fathy/bebo-gallery/views"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// public Error is the type that
// implements the PublicError interface
type publicError string

// make sure that public error will implement the
// PublicError interface
var _ views.PublicError = (*publicError)(nil)

const (
	// ErrNotFound is returned when the resource can not be found from DB
	ErrNotFound publicError = "model: resource not found"

	// ErrPasswordNotCorrect is returned when the password is wrong for user
	ErrPasswordNotCorrect publicError = "model: password is incorrect"

	// ErrEmailNotValidFormat is not correct is used to tell that the email address is not valid
	ErrEmailNotValidFormat publicError = "email address is not valid"

	// ErrEmailIsTaken
	ErrEmailIsTaken publicError = "email address is already taken"

	// ErrPasswordTooShort is returned on password is less than 8 chars
	ErrPasswordTooShort publicError = "password should be at least 8 chars"

	//ErrPasswordRequired
	ErrPasswordRequired publicError = "password can not be empty"

	// ErrRememberRequired is returned when a create or update
	// is attempted without a user remember token hash
	ErrRememberTokenHashRequired publicError = "model: remember token hash is required"

	// ErrRememberTooShort is returned when a remember token is
	// not at least 32 bytes
	ErrRememberTooShort publicError = "model: remember token must be at least 32 bytes"
)

func (e publicError) Error() string {
	return string(e)
}

func (e publicError) PublicErrMsg() string {
	errMsg := strings.Replace(string(e), "model: ", "", 1)
	parts := strings.Split(errMsg, " ")
	if len(parts) > 0 {
		parts[0] = cases.Title(language.English, cases.NoLower).String(parts[0])
	}
	return strings.Join(parts, " ")
}
