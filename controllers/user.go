package controllers

import (
	"fmt"
	"net/http"

	"github.com/abanoub-fathy/bebo-gallery/views"
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

// CreateNewUser will create a new user
func (u *User) CreateNewUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Trying to create new user...")
}
