package controllers

import (
	"log"
	"net/http"
	"time"

	"github.com/abanoub-fathy/bebo-gallery/model"
	"github.com/abanoub-fathy/bebo-gallery/pkg/context"
	"github.com/abanoub-fathy/bebo-gallery/utils"
	"github.com/abanoub-fathy/bebo-gallery/views"
	"github.com/gorilla/mux"
)

const DEFAULT_TOKEN_VALID_DURATION = time.Hour * 120

type User struct {
	SignUpView  *views.View
	LogInView   *views.View
	UserService model.UserService
	router      *mux.Router
}

// NewUser return a pointer to User type which can be used
// as a receiver to call the handler functions
func NewUser(userService model.UserService, muxRouter *mux.Router) *User {
	return &User{
		SignUpView:  views.NewView("base", "user/new"),
		LogInView:   views.NewView("base", "user/login"),
		router:      muxRouter,
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
	// define view params data
	params := views.Params{}

	// define signUpForm
	var form SignUpForm

	// Parse the form
	if err := utils.ParseForm(r, &form); err != nil {
		params.SetAlert(err)
		u.SignUpView.Render(w, r, params)
		return
	}

	// create a new user
	user := &model.User{
		FirstName: form.FirstName,
		LastName:  form.LastName,
		Email:     form.Email,
		Password:  form.Password,
	}

	if err := u.UserService.CreateUser(user); err != nil {
		params.SetAlert(err)
		log.Println(err)
		u.SignUpView.Render(w, r, params)
		return
	}

	// set remember token
	setRemeberTokenToCookie(w, user, DEFAULT_TOKEN_VALID_DURATION)

	// redirect  user to create galleries page
	url, err := u.router.Get(ViewCreateGalleryEndpoint).URL()
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, url.String(), http.StatusFound)
}

// Logut will generate new token to the user and set the cookie in the http
// with valid duration of 0.
//
// this is the same as making our current cookie invlaid
func (u *User) Logout(w http.ResponseWriter, r *http.Request) {
	// get user from conext
	user := context.UserValue(r.Context())

	// set remember token to user
	if err := u.UserService.SaveNewRemeberToken(user); err != nil {
		log.Println("Error while save new remeber toke to user")
	}

	// set remeber token in the cookie and make expire after 0 second
	setRemeberTokenToCookie(w, user, 0)

	// redirect user
	http.Redirect(w, r, "/", http.StatusFound)
}

type LoginForm struct {
	Email    string `schema:"email,required"`
	Password string `schema:"password,required"`
}

// Login is a handler func that will receive data from login Form
//
// and log user in
func (u *User) Login(w http.ResponseWriter, r *http.Request) {
	// define params
	params := views.Params{}

	// define loginForm
	form := LoginForm{}

	// Parse the form
	if err := utils.ParseForm(r, &form); err != nil {
		params.SetAlert(err)
		u.LogInView.Render(w, r, params)
		return
	}

	// authenticate user
	user, err := u.UserService.AuthenticateUser(form.Email, form.Password)
	if err != nil {
		switch err {
		case model.ErrEmailNotValidFormat, model.ErrPasswordNotCorrect:
			params.SetAlert(err)
		case model.ErrNotFound:
			params.SetAlertWithErrMsg("Email address is not found")
		default:
			params.SetAlert(err)
		}
		// render login page with alert
		u.LogInView.Render(w, r, params)
		return
	}

	// set remember token to user
	if err := u.UserService.SaveNewRemeberToken(user); err != nil {
		params.SetAlert(err)
		u.LogInView.Render(w, r, params)
		return
	}

	// set remeber token in the cookie
	setRemeberTokenToCookie(w, user, DEFAULT_TOKEN_VALID_DURATION)

	// redirect to galleries page
	url, err := u.router.Get(ViewGalleriesEndpoint).URL()
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, url.String(), http.StatusFound)
}

// setRemeberTokenToCookie is used to set cookie for user in the response writer
func setRemeberTokenToCookie(w http.ResponseWriter, user *model.User, validDuration time.Duration) {
	// create cookie to store user token
	cookie := &http.Cookie{
		Name:     "token",
		Value:    user.RememberToken,
		HttpOnly: true,
		Expires:  time.Now().Add(validDuration),
	}
	// set cookie in the response writer header
	http.SetCookie(w, cookie)
}
