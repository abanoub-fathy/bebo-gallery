package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/abanoub-fathy/bebo-gallery/controllers"
	"github.com/abanoub-fathy/bebo-gallery/model"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	// Database URI
	var DB_URI = os.Getenv("DATABASE_URI")

	// create new service
	service, err := model.NewService(DB_URI)
	if err != nil {
		panic(err)
	}

	// TODO: we need to (close and migrate) from the service db top level
	defer service.UserService.Close()
	// must(userService.ResetUserDB())
	must(service.UserService.AutoMigrate())

	// create new user controller
	userController := controllers.NewUser(service.UserService)

	// create StaticController
	staticController := controllers.NewStatic()

	// set router
	r := mux.NewRouter()
	r.Handle("/", staticController.Home).Methods("GET")
	r.Handle("/contact", staticController.Contact).Methods("GET")
	r.Handle("/signup", userController.SignUpView).Methods("GET")
	r.HandleFunc("/new", userController.CreateNewUser).Methods("POST")
	r.Handle("/login", userController.LogInView).Methods("GET")
	r.HandleFunc("/login", userController.Login).Methods("POST")
	r.HandleFunc("/cookie", userController.CookieTest).Methods("GET")
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
