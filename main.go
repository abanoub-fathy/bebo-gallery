package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/abanoub-fathy/bebo-gallery/config"
	"github.com/abanoub-fathy/bebo-gallery/controllers"
	"github.com/abanoub-fathy/bebo-gallery/middlewares"
	"github.com/abanoub-fathy/bebo-gallery/model"
	"github.com/abanoub-fathy/bebo-gallery/pkg/email"
	"github.com/abanoub-fathy/bebo-gallery/utils"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
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

	// oAuth
	oAuthConfig := oauth2.Config{
		ClientID:     config.AppConfig.OAuthAppKey,
		ClientSecret: config.AppConfig.OAuthSecretKey,
		Endpoint: oauth2.Endpoint{
			AuthURL:  config.AppConfig.AuthURL,
			TokenURL: config.AppConfig.TokenURL,
		},
		RedirectURL: "http://localhost:3000/oauth/dropbox/callback",
	}

	r.HandleFunc("/oauth/dropbox/connect", func(w http.ResponseWriter, r *http.Request) {
		state := csrf.Token(r)

		// create a cookie with the state
		coockie := &http.Cookie{
			Name:     "oauth_state",
			Value:    state,
			Path:     "/",
			HttpOnly: true,
		}
		// setting the cookie
		http.SetCookie(w, coockie)

		// generate and redirect to authURL
		url := oAuthConfig.AuthCodeURL(state)
		http.Redirect(w, r, url, http.StatusFound)
	})

	r.HandleFunc("/oauth/dropbox/callback", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		code := query.Get("code")
		state := query.Get("state")

		// get state from request cookie
		cookie, err := r.Cookie("oauth_state")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		} else if cookie == nil || cookie.Value != state {
			http.Error(w, "invalid state", http.StatusBadRequest)
			return
		}

		// expire the state cookie
		cookie.Value = ""
		cookie.Expires = time.Now()
		http.SetCookie(w, cookie)

		// exchange the code
		token, err := oAuthConfig.Exchange(r.Context(), code)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "%+v", token)

		// w.Header().Set("Content-Type", "application/json")
		// json.NewEncoder(w).Encode(token)
	})

	// user routes
	r.HandleFunc("/signup", userController.NewUser).Methods("GET")
	r.HandleFunc("/new", userController.CreateNewUser).Methods("POST")
	r.Handle("/login", userController.LogInView).Methods("GET")
	r.HandleFunc("/login", userController.Login).Methods("POST")
	r.HandleFunc("/password/forget", userController.ForgetPasswordPage).Methods("GET")
	r.HandleFunc("/password/forget", userController.ForgetPassword).Methods("POST")
	r.HandleFunc("/password/reset", userController.ResetPasswordPage).Methods("GET")
	r.HandleFunc("/password/reset", userController.ResetPassword).Methods("POST")
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
