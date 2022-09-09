package main

import (
	"fmt"
	"net/http"

	"github.com/abanoub-fathy/bebo-gallery/controllers"
	"github.com/gorilla/mux"
)

func main() {
	// create new user controller
	userController := controllers.NewUser()

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
