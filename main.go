package main

import (
	"fmt"
	"net/http"

	"github.com/abanoub-fathy/bebo-gallery/views"
	"github.com/gorilla/mux"
)

// templates global vars
var (
	homeView    *views.View
	contactView *views.View
)

func main() {
	// create template views
	homeView = views.NewView("base", "views/home.gohtml")
	contactView = views.NewView("base", "views/contact.gohtml")

	// set router
	r := mux.NewRouter()
	r.HandleFunc("/", Home)
	r.HandleFunc("/contact", Contact)
	r.NotFoundHandler = http.HandlerFunc(NotFound)

	// start the app
	fmt.Println("🚀🚀 Server is working on http://localhost:3000")
	must(http.ListenAndServe(":3000", r))
}

// Home is the handlerFunc for the home page
func Home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	must(homeView.Render(w, nil))

}

// Contact is the handlerFunc for the contact page
func Contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	must(contactView.Render(w, nil))
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
