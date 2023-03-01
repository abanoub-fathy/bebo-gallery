package controllers

import (
	"fmt"
	"net/http"

	"github.com/abanoub-fathy/bebo-gallery/model"
	"github.com/abanoub-fathy/bebo-gallery/pkg/context"
	"github.com/abanoub-fathy/bebo-gallery/utils"
	"github.com/abanoub-fathy/bebo-gallery/views"
	"github.com/gorilla/mux"
)

type Gallery struct {
	ShowGalleryView   *views.View
	CreateGalleryView *views.View
	GalleryService    model.GalleryService
}

// NewGallery return a pointer to Gallery type which can be used
// as a receiver to call the handler functions
func NewGallery(galleryService model.GalleryService) *Gallery {
	return &Gallery{
		ShowGalleryView:   views.NewView("base", "gallery/gallery"),
		CreateGalleryView: views.NewView("base", "gallery/new"),
		GalleryService:    galleryService,
	}
}

func (g *Gallery) ViewGallery(w http.ResponseWriter, r *http.Request) {
	// get gallery id
	galleryID := mux.Vars(r)["galleryID"]

	// fetch gallery by id
	gallery, err := g.GalleryService.FindByID(galleryID)
	if err != nil {
		// redirect user to not found
		http.Redirect(w, r, "/notFound", http.StatusPermanentRedirect)
		return
	}

	// render the gallery
	err = g.ShowGalleryView.Render(w, views.Params{
		Data: gallery,
	})
	if err != nil {
		fmt.Println("err while rendering gallery", err)
	}
}

type createGalleryForm struct {
	Title string `schema:"title"`
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

	// get user from ctx
	user := context.UserValue(r.Context())

	gallery := &model.Gallery{
		Title:  form.Title,
		UserID: user.ID,
	}

	err := g.GalleryService.CreateGallery(gallery)
	if err != nil {
		params.SetAlert(err)
		g.CreateGalleryView.Render(w, params)
		return
	}

	w.Write([]byte("created successfully"))
}
