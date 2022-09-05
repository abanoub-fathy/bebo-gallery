package controllers

import (
	"fmt"
	"net/http"

	"github.com/abanoub-fathy/bebo-gallery/views"
	"github.com/gorilla/schema"
)

type User struct {
	View *views.View
}

// NewUser return a pointer to User type which can be used
// as a receiver to call the handler functions
func NewUser() *User {
	return &User{
		View: views.NewView("base", "views/user/new.gohtml"),
	}
}

// RenderUserSignUpForm
func (u *User) RenderUserSignUpForm(w http.ResponseWriter, r *http.Request) {
	if err := u.View.Render(w, nil); err != nil {
		panic(err)
	}
}

type SignUpForm struct {
	Email    string `schema:"email,required"`
	Password string `schema:"password,required"`
}

// CreateNewUser will create a new user
func (u *User) CreateNewUser(w http.ResponseWriter, r *http.Request) {
	// Parse the form
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}

	// create schema decoder
	var decoder = schema.NewDecoder()

	// define signUpForm
	var form SignUpForm

	// decode the form
	err = decoder.Decode(&form, r.PostForm)
	if err != nil {
		panic(err)
	}

	fmt.Fprintln(w, "email=", form.Email, "password=", form.Password)
}
