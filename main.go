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
	homeView = views.NewView("views/home.gohtml")
	contactView = views.NewView("views/contact.gohtml")

	// set router
	r := mux.NewRouter()
	r.HandleFunc("/", Home)
	r.HandleFunc("/contact", Contact)
	r.NotFoundHandler = http.HandlerFunc(NotFound)

	// start the app
	fmt.Println("🚀🚀 Server is working on http://localhost:3000")
	if err := http.ListenAndServe(":3000", r); err != nil {
		panic(err)
	}
}

// Home is the handlerFunc for the home page
func Home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	homeView.Template.Execute(w, nil)
}

func Contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	contactView.Template.Execute(w, nil)
}

// NotFound is the handlerFunc for not found page
func NotFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "<h1>404 Not Found</h1>")
}
