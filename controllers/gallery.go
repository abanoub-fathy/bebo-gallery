package controllers

import (
	"fmt"
	"net/http"

	"github.com/abanoub-fathy/bebo-gallery/model"
	"github.com/abanoub-fathy/bebo-gallery/pkg/context"
	"github.com/abanoub-fathy/bebo-gallery/utils"
	"github.com/abanoub-fathy/bebo-gallery/views"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

const (
	ViewGalleryEndpoint       = "view_gallery_endpoint"
	ViewGalleriesEndpoint     = "view_galleries_endpoint"
	ViewCreateGalleryEndpoint = "view_create_gallery_end_point"
	EditGalleryPageEndpoint   = "edit_gallery_page_end_point"
)

const (
	PARSE_FORM_MAX_MEMORY = 30 << 20
)

type Gallery struct {
	ShowGalleryView       *views.View
	ShowUserGalleriesView *views.View
	CreateGalleryView     *views.View
	EditGalleryView       *views.View
	GalleryService        model.GalleryService
	ImageService          model.ImageService
	router                *mux.Router
}

// NewGallery return a pointer to Gallery type which can be used
// as a receiver to call the handler functions
func NewGallery(galleryService model.GalleryService, imageService model.ImageService, muxRouter *mux.Router) *Gallery {
	return &Gallery{
		ShowGalleryView:       views.NewView("base", "gallery/gallery"),
		ShowUserGalleriesView: views.NewView("base", "gallery/user_galleries"),
		CreateGalleryView:     views.NewView("base", "gallery/new"),
		EditGalleryView:       views.NewView("base", "gallery/edit"),
		GalleryService:        galleryService,
		ImageService:          imageService,
		router:                muxRouter,
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

	// fetch gallery images
	gallery.Images, _ = g.ImageService.GetImagesByGalleryID(gallery.ID)

	// render the gallery
	err = g.ShowGalleryView.Render(w, r, views.Params{
		Data: gallery,
	})
	if err != nil {
		fmt.Println("err while rendering gallery", err)
	}
}

func (g *Gallery) ShowUserGalleriesPage(w http.ResponseWriter, r *http.Request) {
	// get user from conext
	user := context.UserValue(r.Context())

	// get the galleries of the user
	galleries, err := g.GalleryService.FindByUserID(user.ID)
	if err != nil {
		http.Error(w, "could not get galleries by used id", http.StatusInternalServerError)
		return
	}

	// render user galleries page
	params := views.Params{
		Data: galleries,
	}

	if err = g.ShowUserGalleriesView.Render(w, r, params); err != nil {
		http.Error(w, "could not show your galleries", http.StatusInternalServerError)
		return
	}
}

func (g *Gallery) EditGalleryPage(w http.ResponseWriter, r *http.Request) {
	// get gallery id
	galleryID := mux.Vars(r)["galleryID"]

	// fetch gallery by id
	gallery, err := g.GalleryService.FindByID(galleryID)
	if err != nil {
		// redirect user to not found
		http.Redirect(w, r, "/notFound", http.StatusPermanentRedirect)
		return
	}

	// get user from conext
	user := context.UserValue(r.Context())

	// check that the user own the gallery
	if !uuid.Equal(user.ID, gallery.UserID) {
		// redirect user to not found
		http.Redirect(w, r, "/notFound", http.StatusPermanentRedirect)
		return
	}

	// fetch gallery images
	gallery.Images, _ = g.ImageService.GetImagesByGalleryID(gallery.ID)

	// render the gallery
	err = g.EditGalleryView.Render(w, r, views.Params{
		Data: gallery,
	})
	if err != nil {
		fmt.Println("err while rendering gallery", err)
	}
}

func (g *Gallery) EditGallery(w http.ResponseWriter, r *http.Request) {
	// get gallery id variable
	galleryID := mux.Vars(r)["galleryID"]

	// fetch gallery by id
	gallery, err := g.GalleryService.FindByID(galleryID)
	if err != nil {
		// redirect user to not found
		http.Redirect(w, r, "/notFound", http.StatusPermanentRedirect)
		return
	}

	// get user from conext
	user := context.UserValue(r.Context())

	// check that the user own the gallery
	if !uuid.Equal(user.ID, gallery.UserID) {
		// redirect user to not found
		http.Redirect(w, r, "/notFound", http.StatusPermanentRedirect)
		return
	}

	// define view params data
	params := views.Params{}

	// define createGalleryForm
	var form createGalleryForm

	// Parse the form
	if err := utils.ParseForm(r, &form); err != nil {
		params.SetAlert(err)
		g.EditGalleryView.Render(w, r, params)
		return
	}

	// update the gallery
	gallery.Title = form.Title

	err = g.GalleryService.Update(gallery)
	if err != nil {
		params.SetAlert(err)
		params.Data = gallery
		g.EditGalleryView.Render(w, r, params)
		return
	}

	// redirect user to show gallery page
	url, err := g.router.Get(ViewGalleryEndpoint).URL("galleryID", gallery.ID.String())
	if err != nil {
		http.Redirect(w, r, "/", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, url.String(), http.StatusFound)
}

func (g *Gallery) UploadImage(w http.ResponseWriter, r *http.Request) {
	// get gallery id variable
	galleryID := mux.Vars(r)["galleryID"]

	// fetch gallery by id
	gallery, err := g.GalleryService.FindByID(galleryID)
	if err != nil {
		// redirect user to not found
		http.Redirect(w, r, "/notFound", http.StatusPermanentRedirect)
		return
	}

	// get user from conext
	user := context.UserValue(r.Context())

	// check that the user own the gallery
	if !uuid.Equal(user.ID, gallery.UserID) {
		// redirect user to not found
		http.Redirect(w, r, "/notFound", http.StatusPermanentRedirect)
		return
	}

	// define view params data
	params := views.Params{
		Data: gallery,
	}

	// parse multipart gallery
	if err = r.ParseMultipartForm(PARSE_FORM_MAX_MEMORY); err != nil {
		params.SetAlert(err)
		g.EditGalleryView.Render(w, r, params)
		return
	}

	fileHeaders := r.MultipartForm.File["images"]

	for _, f := range fileHeaders {
		// open the file
		file, err := f.Open()
		if err != nil {
			params.SetAlert(err)
			g.EditGalleryView.Render(w, r, params)
			return
		}
		defer file.Close()

		err = g.ImageService.CreateImage(file, gallery.ID, f.Filename)
		if err != nil {
			params.SetAlert(err)
			g.EditGalleryView.Render(w, r, params)
			return
		}
	}

	// redirect user to show gallery page
	url, err := g.router.Get(EditGalleryPageEndpoint).URL("galleryID", gallery.ID.String())
	if err != nil {
		http.Redirect(w, r, "/", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, url.String(), http.StatusFound)
}

func (g *Gallery) DeleteImage(w http.ResponseWriter, r *http.Request) {
	// get gallery id variable
	galleryID := mux.Vars(r)["galleryID"]

	// fetch gallery by id
	gallery, err := g.GalleryService.FindByID(galleryID)
	if err != nil {
		// redirect user to not found
		http.Redirect(w, r, "/notFound", http.StatusPermanentRedirect)
		return
	}

	// get user from conext
	user := context.UserValue(r.Context())

	// check that the user own the gallery
	if !uuid.Equal(user.ID, gallery.UserID) {
		// redirect user to not found
		http.Redirect(w, r, "/notFound", http.StatusPermanentRedirect)
		return
	}

	// get fileName
	fileName := mux.Vars(r)["fileName"]

	// define image
	image := model.Image{
		GalleryID: galleryID,
		FileName:  fileName,
	}

	// delete the image
	if err = g.ImageService.DeleteImage(&image); err != nil {
		params := views.Params{
			Data: gallery,
		}
		params.SetAlert(err)
		g.EditGalleryView.Render(w, r, params)
		return
	}

	// redirect user to show gallery page
	url, err := g.router.Get(EditGalleryPageEndpoint).URL("galleryID", gallery.ID.String())
	if err != nil {
		http.Redirect(w, r, "/", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, url.String(), http.StatusFound)
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
		g.CreateGalleryView.Render(w, r, params)
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
		g.CreateGalleryView.Render(w, r, params)
		return
	}

	url, err := g.router.Get(ViewGalleriesEndpoint).URL()
	if err != nil {
		http.Redirect(w, r, "/", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, url.String(), http.StatusFound)
}

func (g *Gallery) DeleteGallery(w http.ResponseWriter, r *http.Request) {
	// find gallery by id
	galleryID := mux.Vars(r)["galleryID"]

	// fetch the gallery by ID
	gallery, err := g.GalleryService.FindByID(galleryID)
	if err != nil {
		switch err {
		case model.ErrNotFound:
			// redirect user to the 404 page
			http.Redirect(w, r, "/notFound", http.StatusPermanentRedirect)
		default:
			http.Error(w, "could not get gallery", http.StatusInternalServerError)
		}
		return
	}

	// get the user from context
	user := context.UserValue(r.Context())

	// check that the user own the gallery
	if gallery.UserID != user.ID {
		// redirect user to the 404 page
		http.Redirect(w, r, "/notFound", http.StatusPermanentRedirect)
		return
	}

	// define view params
	params := views.Params{}

	// delete gallery
	if err := g.GalleryService.Delete(gallery); err != nil {
		params.SetAlert(err)
		params.Data = gallery
		g.EditGalleryView.Render(w, r, params)
		return
	}

	// return the gallery
	url, err := g.router.Get(ViewGalleriesEndpoint).URL()
	if err != nil {
		http.Redirect(w, r, "/", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, url.String(), http.StatusFound)
}
