package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

// templates global vars
var (
	homeTemplate    *template.Template
	contactTemplate *template.Template
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", Home)
	r.HandleFunc("/contact", Contact)
	r.NotFoundHandler = http.HandlerFunc(NotFound)

	// parse the home template
	var err error
	if homeTemplate, err = template.ParseFiles("views/home.gohtml"); err != nil {
		panic(err)
	}

	// parse the contact template
	if contactTemplate, err = template.ParseFiles("views/contact.gohtml"); err != nil {
		panic(err)
	}

	// start the app
	fmt.Println("🚀🚀 Server is working on http://localhost:3000")
	if err := http.ListenAndServe(":3000", r); err != nil {
		panic(err)
	}
}

// Home is the handlerFunc for the home page
func Home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	homeTemplate.Execute(w, nil)
}

func Contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	contactTemplate.Execute(w, nil)
}

// NotFound is the handlerFunc for not found page
func NotFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "<h1>404 Not Found</h1>")
}
