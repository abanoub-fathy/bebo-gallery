package controllers

import "github.com/abanoub-fathy/bebo-gallery/views"

// Static type contains the static views
type Static struct {
	Home     *views.View
	Contact  *views.View
	NotFound *views.View
}

// NewStatic is constructor func for creating new static controller with
// all static pages hard coded inside it
func NewStatic() *Static {
	return &Static{
		Home:     views.NewView("base", "static/home"),
		Contact:  views.NewView("base", "static/contact"),
		NotFound: views.NewView("base", "static/notFound"),
	}
}
