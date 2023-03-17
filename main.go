package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/abanoub-fathy/bebo-gallery/controllers"
	"github.com/abanoub-fathy/bebo-gallery/middlewares"
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
	utils.Must(err)

	// defer closing the services
	defer service.Close()

	// migrate all the models to the DB
	utils.Must(service.AutoMigrate())

	// creat middleware
	requireUserMiddleWare := middlewares.RequireUser{
		Service: service,
	}

	userMiddleWare := middlewares.UserMiddleware{
		Service: service,
	}

	// set router
	r := mux.NewRouter()

	// create StaticController
	staticController := controllers.NewStatic()

	// static routes
	r.Handle("/", staticController.Home).Methods("GET")
	r.Handle("/contact", staticController.Contact).Methods("GET")
	r.NotFoundHandler = staticController.NotFound

	// create new user controller
	userController := controllers.NewUser(service.UserService, r)

	// user routes
	r.Handle("/signup", userController.SignUpView).Methods("GET")
	r.HandleFunc("/new", userController.CreateNewUser).Methods("POST")
	r.Handle("/login", userController.LogInView).Methods("GET")
	r.HandleFunc("/login", userController.Login).Methods("POST")

	// create gallery controllers
	galleryController := controllers.NewGallery(service.GalleryService, service.ImageService, r)

	// gallery routes
	r.Handle("/galleries/new", requireUserMiddleWare.Apply(galleryController.CreateGalleryView)).Methods("GET").Name(controllers.ViewCreateGalleryEndpoint)
	r.HandleFunc("/galleries/{galleryID}", galleryController.ViewGallery).Methods("GET").Name(controllers.ViewGalleryEndpoint)
	r.HandleFunc("/galleries", requireUserMiddleWare.ApplyFunc(galleryController.CreateNewGallery)).Methods("POST")
	r.HandleFunc("/galleries", requireUserMiddleWare.ApplyFunc(galleryController.ShowUserGalleriesPage)).Methods("GET").Name(controllers.ViewGalleriesEndpoint)
	r.HandleFunc("/galleries/{galleryID}/edit", requireUserMiddleWare.ApplyFunc(galleryController.EditGalleryPage)).Methods("GET")
	r.HandleFunc("/galleries/{galleryID}/edit", requireUserMiddleWare.ApplyFunc(galleryController.EditGallery)).Methods("POST")
	r.HandleFunc("/galleries/{galleryID}/images", requireUserMiddleWare.ApplyFunc(galleryController.UploadImage)).Methods("POST")
	r.HandleFunc("/galleries/{galleryID}/delete", requireUserMiddleWare.ApplyFunc(galleryController.DeleteGallery)).Methods("POST")

	// start the app
	fmt.Println("🚀🚀 Server is working on http://localhost:3000")
	utils.Must(http.ListenAndServe(":3000", userMiddleWare.UserInCtxApply(r)))
}
