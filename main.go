package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/abanoub-fathy/bebo-gallery/controllers"
	"github.com/abanoub-fathy/bebo-gallery/model"
	"github.com/abanoub-fathy/bebo-gallery/utils"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// load env file
	utils.Must(godotenv.Load())

	// Database URI
	var DB_URI = os.Getenv("DATABASE_URI")

	// create new service
	service, err := model.NewService(DB_URI)
	if err != nil {
		panic(err)
	}

	// defer closing the services
	defer service.Close()

	// migrate all the models to the DB
	utils.Must(service.AutoMigrate())

	// create new user controller
	userController := controllers.NewUser(service.UserService)

	// create StaticController
	staticController := controllers.NewStatic()

	// create gallery controllers
	galleryController := controllers.NewGallery(service.GalleryService)

	// set router
	r := mux.NewRouter()

	// static routes
	r.Handle("/", staticController.Home).Methods("GET")
	r.Handle("/contact", staticController.Contact).Methods("GET")
	r.NotFoundHandler = staticController.NotFound

	// user routes
	r.Handle("/signup", userController.SignUpView).Methods("GET")
	r.HandleFunc("/new", userController.CreateNewUser).Methods("POST")
	r.Handle("/login", userController.LogInView).Methods("GET")
	r.HandleFunc("/login", userController.Login).Methods("POST")
	r.HandleFunc("/cookie", userController.CookieTest).Methods("GET")

	// gallery routes
	r.Handle("/galleries/new", galleryController.CreateGalleryView).Methods("GET")

	// start the app
	fmt.Println("🚀🚀 Server is working on http://localhost:3000")
	utils.Must(http.ListenAndServe(":3000", r))
}
