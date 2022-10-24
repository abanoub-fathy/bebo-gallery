package controllers

import (
	"fmt"
	"net/http"

	"github.com/abanoub-fathy/bebo-gallery/model"
	"github.com/abanoub-fathy/bebo-gallery/utils"
	"github.com/abanoub-fathy/bebo-gallery/views"
)

type User struct {
	SignUpView  *views.View
	LogInView   *views.View
	UserService *model.UserService
}

// NewUser return a pointer to User type which can be used
// as a receiver to call the handler functions
func NewUser(userService *model.UserService) *User {
	return &User{
		SignUpView:  views.NewView("base", "user/new"),
		LogInView:   views.NewView("base", "user/login"),
		UserService: userService,
	}
}

type SignUpForm struct {
	FirstName string `schema:"firstName,required"`
	LastName  string `schema:"lastName,required"`
	Email     string `schema:"email,required"`
	Password  string `schema:"password,required"`
}

// CreateNewUser is a handler func that will receive data from sigup Form
// and create a new user
//
// and save it to the database
func (u *User) CreateNewUser(w http.ResponseWriter, r *http.Request) {
	// define signUpForm
	var form SignUpForm

	// Parse the form
	if err := utils.ParseForm(r, &form); err != nil {
		panic(err)
	}

	// create a new user
	user := &model.User{
		FirstName: form.FirstName,
		LastName:  form.LastName,
		Email:     form.Email,
		Password:  form.Password,
	}

	if err := u.UserService.CreateUser(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "id=", user.ID, "email=", user.Email, "firstName=", user.FirstName, "lastName=", user.LastName)
}

type LoginForm struct {
	Email    string `schema:"email,required"`
	Password string `schema:"password,required"`
}

// Login is a handler func that will receive data from login Form
//
// and log user in
func (u *User) Login(w http.ResponseWriter, r *http.Request) {
	// define loginForm
	form := LoginForm{}

	// Parse the form
	if err := utils.ParseForm(r, &form); err != nil {
		panic(err)
	}

	// authenticate user
	user, err := u.UserService.AuthenticateUser(form.Email, form.Password)
	if err != nil {
		switch err {
		case model.ErrNotValidEmail:
			fmt.Fprintln(w, "not valid email address")
		case model.ErrNotFound:
			fmt.Fprintln(w, "user email not found")
		case model.ErrPasswordNotCorrect:
			fmt.Fprintln(w, "password is incorrect")
		default:
			fmt.Fprintln(w, err)
		}
	}

	fmt.Fprintln(w, user)

}
