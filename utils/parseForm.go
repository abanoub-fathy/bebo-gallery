package utils

import (
	"net/http"
	"net/url"

	"github.com/gorilla/schema"
)

func ParseForm(r *http.Request, dst interface{}) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}
	return parseValues(r.PostForm, dst)
}

func ParseURLParams(r *http.Request, dst interface{}) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}
	return parseValues(r.Form, dst)
}

func parseValues(values url.Values, dst interface{}) error {
	// create schema decoder
	decoder := schema.NewDecoder()

	// ignore unknown keys
	decoder.IgnoreUnknownKeys(true)

	// decode the values in the destination
	err := decoder.Decode(dst, values)
	if err != nil {
		return err
	}

	return nil
}
