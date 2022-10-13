package controllers

import (
	"fmt"
	"net/http"

	"github.com/abanoub-fathy/bebo-gallery/model"
	"github.com/abanoub-fathy/bebo-gallery/utils"
	"github.com/abanoub-fathy/bebo-gallery/views"
)

type User struct {
	View        *views.View
	UserService *model.UserService
}

// NewUser return a pointer to User type which can be used
// as a receiver to call the handler functions
func NewUser(userService *model.UserService) *User {
	return &User{
		View:        views.NewView("base", "user/new"),
		UserService: userService,
	}
}

// RenderUserSignUpForm
func (u *User) RenderUserSignUpForm(w http.ResponseWriter, r *http.Request) {
	if err := u.View.Render(w, nil); err != nil {
		panic(err)
	}
}

type SignUpForm struct {
	FirstName string `schema:"firstName,required"`
	LastName  string `schema:"lastName,required"`
	Email     string `schema:"email,required"`
	Password  string `schema:"password,required"`
}

// CreateNewUser will create a new user
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
	}

	if err := u.UserService.CreateUser(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "id=", user.ID, "email=", user.Email, "firstName=", user.FirstName, "lastName=", user.LastName)
}
