package controllers

import "github.com/abanoub-fathy/bebo-gallery/views"

type Static struct {
	Home    *views.View
	Contact *views.View
}

func NewStatic() *Static {
	return &Static{
		Home:    views.NewView("base", "views/static/home.gohtml"),
		Contact: views.NewView("base", "views/static/contact.gohtml"),
	}
}
