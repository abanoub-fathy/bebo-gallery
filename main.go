package main

import (
	"fmt"
	"net/http"

	"github.com/abanoub-fathy/bebo-gallery/controllers"
	"github.com/abanoub-fathy/bebo-gallery/model"
	"github.com/gorilla/mux"
)

func main() {
	// Database URI
	const DB_URI = "postgresql://postgres:popTop123@localhost:5432/bebo-gallery?sslmode=disable"

	// create new UserService
	userService, err := model.NewUserService(DB_URI)
	if err != nil {
		panic(err)
	}
	defer userService.Close()
	// must(userService.ResetUserDB())
	must(userService.AutoMigrate())

	// create new user controller
	userController := controllers.NewUser(userService)

	// create StaticController
	staticController := controllers.NewStatic()

	// set router
	r := mux.NewRouter()
	r.Handle("/", staticController.Home).Methods("GET")
	r.Handle("/contact", staticController.Contact).Methods("GET")
	r.HandleFunc("/signup", userController.RenderUserSignUpForm).Methods("GET")
	r.HandleFunc("/new", userController.CreateNewUser).Methods("POST")
	r.NotFoundHandler = staticController.NotFound

	// start the app
	fmt.Println("🚀🚀 Server is working on http://localhost:3000")
	must(http.ListenAndServe(":3000", r))
}

// must is used to panic an error if exist
func must(err error) {
	if err != nil {
		panic(err)
	}
}
