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

	// ignore unknown keys
	decoder.IgnoreUnknownKeys(true)

	err = decoder.Decode(data, r.PostForm)
	if err != nil {
		return err
	}

	return nil
}
