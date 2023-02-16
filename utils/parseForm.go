package utils

import (
	"net/http"

	"github.com/gorilla/schema"
)

func ParseForm(r *http.Request, data interface{}) error {
	// create schema decoder
	decoder := schema.NewDecoder()
	err := r.ParseForm()
	if err != nil {
		return err
	}

	err = decoder.Decode(data, r.PostForm)
	if err != nil {
		return err
	}

	return nil
}

// Must is used to panic an error if exist
func Must(err error) {
	if err != nil {
		panic(err)
	}
}
