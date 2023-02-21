package controllers

import (
	"github.com/abanoub-fathy/bebo-gallery/model"
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
