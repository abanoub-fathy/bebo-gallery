package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", Hello)
	r.NotFoundHandler = http.HandlerFunc(NotFound)

	// start the app
	fmt.Println("🚀🚀 Server is working on http://localhost:3000")
	if err := http.ListenAndServe(":3000", r); err != nil {
		panic(err)
	}
}

// Hello is the handlerFunc for the home page
func Hello(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, "<h1>Hello in bebo gallery</h1>")
}

// NotFound is the handlerFunc for not found page
func NotFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "<h1>404 Not Found</h1>")
}
