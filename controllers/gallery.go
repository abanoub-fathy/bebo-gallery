package controllers

import (
	"net/http"

	"github.com/abanoub-fathy/bebo-gallery/model"
	"github.com/abanoub-fathy/bebo-gallery/utils"
	"github.com/abanoub-fathy/bebo-gallery/views"
)

type Gallery struct {
	CreateGalleryView *views.View
	GalleryService    model.GalleryService
}

// NewGallery return a pointer to Gallery type which can be used
// as a receiver to call the handler functions
func NewGallery(galleryService model.GalleryService) *Gallery {
	return &Gallery{
		CreateGalleryView: views.NewView("base", "gallery/new"),
		GalleryService:    galleryService,
	}
}

type createGalleryForm struct {
	Title string `schema:"title,required"`
}

func (g *Gallery) CreateNewGallery(w http.ResponseWriter, r *http.Request) {
	// define view params data
	params := views.Params{}

	// define createGalleryForm
	var form createGalleryForm

	// Parse the form
	if err := utils.ParseForm(r, &form); err != nil {
		// set the alert
		params.SetAlert(err)

		// render the create gallery view with params
		g.CreateGalleryView.Render(w, params)
		return
	}

	gallery := &model.Gallery{
		Title: form.Title,
	}

	err := g.GalleryService.CreateGallery(gallery)
	if err != nil {
		params.SetAlert(err)
		g.CreateGalleryView.Render(w, params)
		return
	}

	w.Write([]byte("created successfully"))
}
