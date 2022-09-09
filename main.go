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
	r.NotFoundHandler = http.HandlerFunc(NotFound)

	// start the app
	fmt.Println("🚀🚀 Server is working on http://localhost:3000")
	must(http.ListenAndServe(":3000", r))
}

// NotFound is the handlerFunc for not found page
func NotFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "<h1>404 Not Found</h1>")
}

// must is used to panic an error if exist
func must(err error) {
	if err != nil {
		panic(err)
	}
}
