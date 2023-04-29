package main

import (
	"fmt"
	"net/http"

	"github.com/abanoub-fathy/bebo-gallery/config"
	"github.com/abanoub-fathy/bebo-gallery/controllers"
	"github.com/abanoub-fathy/bebo-gallery/middlewares"
	"github.com/abanoub-fathy/bebo-gallery/model"
	"github.com/abanoub-fathy/bebo-gallery/pkg/email"
	"github.com/abanoub-fathy/bebo-gallery/utils"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
)

func main() {
	// create new email client
	emailClient := email.NewClient(config.AppConfig.EmailAPIKey)

	// create new service
	service, err := model.NewService(config.AppConfig.DatabaseURI)
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

	// serve static assets
	assetsServerHandler := http.FileServer(http.Dir("./views/assets/"))
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", assetsServerHandler))

	// file server
	fileServerHandler := http.FileServer(http.Dir("./images"))
	r.PathPrefix("/images/").Handler(http.StripPrefix("/images/", fileServerHandler))

	// create StaticController
	staticController := controllers.NewStatic()

	// static routes
	r.Handle("/", staticController.Home).Methods("GET")
	r.Handle("/contact", staticController.Contact).Methods("GET")
	r.NotFoundHandler = staticController.NotFound

	// create new user controller
	userController := controllers.NewUser(service.UserService, r, emailClient)

	// user routes
	r.HandleFunc("/signup", userController.NewUser).Methods("GET")
	r.HandleFunc("/new", userController.CreateNewUser).Methods("POST")
	r.Handle("/login", userController.LogInView).Methods("GET")
	r.HandleFunc("/login", userController.Login).Methods("POST")
	r.HandleFunc("/logout", requireUserMiddleWare.ApplyFunc(userController.Logout)).Methods("POST")

	// create gallery controllers
	galleryController := controllers.NewGallery(service.GalleryService, service.ImageService, r)

	// gallery routes
	r.Handle("/galleries/new", requireUserMiddleWare.Apply(galleryController.CreateGalleryView)).Methods("GET").Name(controllers.ViewCreateGalleryEndpoint)
	r.HandleFunc("/galleries/{galleryID}", galleryController.ViewGallery).Methods("GET").Name(controllers.ViewGalleryEndpoint)
	r.HandleFunc("/galleries", requireUserMiddleWare.ApplyFunc(galleryController.CreateNewGallery)).Methods("POST")
	r.HandleFunc("/galleries", requireUserMiddleWare.ApplyFunc(galleryController.ShowUserGalleriesPage)).Methods("GET").Name(controllers.ViewGalleriesEndpoint)
	r.HandleFunc("/galleries/{galleryID}/edit", requireUserMiddleWare.ApplyFunc(galleryController.EditGalleryPage)).Methods("GET").Name(controllers.EditGalleryPageEndpoint)
	r.HandleFunc("/galleries/{galleryID}/edit", requireUserMiddleWare.ApplyFunc(galleryController.EditGallery)).Methods("POST")
	r.HandleFunc("/galleries/{galleryID}/images", requireUserMiddleWare.ApplyFunc(galleryController.UploadImage)).Methods("POST")
	r.HandleFunc("/galleries/{galleryID}/images/{fileName}/delete", requireUserMiddleWare.ApplyFunc(galleryController.DeleteImage)).Methods("POST")
	r.HandleFunc("/galleries/{galleryID}/delete", requireUserMiddleWare.ApplyFunc(galleryController.DeleteGallery)).Methods("POST")

	// CSRF Protection
	CSRF := csrf.Protect([]byte(config.AppConfig.CSRFKey), csrf.Secure(config.AppConfig.IsProductionEnv))

	// start the app
	fmt.Printf("🚀🚀 Server is working on http://localhost:%v\n", config.AppConfig.Port)
	utils.Must(http.ListenAndServe(fmt.Sprintf(":%v", config.AppConfig.Port), CSRF(userMiddleWare.UserInCtxApply(r))))
}
