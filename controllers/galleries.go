package controllers

import (
	"jiji/middlewares"
	"jiji/models"
	"jiji/views"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

const (
	ShowGallery      = "show_gallery"
	msgSuccessUpdate = "Successfully updated gallery"
)

type Galleries struct {
	New      *views.View
	ShowView *views.View
	EditView *views.View
	gs       models.GalleryService
	r        *mux.Router
}

type GalleryForm struct {
	Title string `schema:"title"`
}

func NewGalleries(gs models.GalleryService, r *mux.Router) *Galleries {
	return &Galleries{
		New:      views.NewView("bootstrap", "galleries/new"),
		ShowView: views.NewView("bootstrap", "galleries/show"),
		EditView: views.NewView("bootstrap", "galleries/edit"),
		gs:       gs,
		r:        r,
	}
}

// POST /galleries
func (g *Galleries) Create(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var galleryForm GalleryForm

	err := parseForm(r, &galleryForm)
	if err != nil {
		vd.SetAlert(err)
		g.New.Render(w, vd)
		return
	}

	user := middlewares.LookUpUserFromContext(r.Context())

	gallery := models.Gallery{
		Title:  galleryForm.Title,
		UserId: user.ID,
	}

	err = g.gs.Create(&gallery)
	if err != nil {
		vd.SetAlert(err)
		g.New.Render(w, vd)
		return
	}

	url, err := g.r.Get(ShowGallery).URL("id", strconv.Itoa(int(gallery.ID)))
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	http.Redirect(w, r, url.Path, http.StatusFound)
}

// GET /galleries/:id
func (g *Galleries) Show(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.getGalleryById(w, r)
	if err != nil {
		return // At this point err is already handled
	}

	var vd views.Data
	vd.Yield = gallery
	g.ShowView.Render(w, vd)
}

// GET /galleries/:id/edit
func (g *Galleries) Edit(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.getGalleryById(w, r)
	if err != nil {
		return // At this point err is already handled
	}

	// A user needs logged in to access this page, so we can assume that
	// the RequireUser middleware has run and set the user for us in the request context.
	user := middlewares.LookUpUserFromContext(r.Context())
	if gallery.UserId != user.ID {
		http.Error(w, "You do not have permission to edit this gallery", http.StatusForbidden)
		return
	}

	var vd views.Data
	vd.Yield = gallery
	g.EditView.Render(w, vd)
}

func (g *Galleries) Update(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.getGalleryById(w, r)
	if err != nil {
		return // At this point err is already handled
	}

	// A user needs logged in to access this page, so we can assume that
	// the RequireUser middleware has run and set the user for us in the request context.
	user := middlewares.LookUpUserFromContext(r.Context())
	if gallery.UserId != user.ID {
		http.Error(w, "You do not have permission to edit this gallery", http.StatusForbidden)
		return
	}

	var vd views.Data
	var galleryForm GalleryForm

	err = parseForm(r, &galleryForm)
	if err != nil {
		vd.SetAlert(err)
		g.New.Render(w, vd)
		return
	}

	gallery.Title = galleryForm.Title

	err = g.gs.Update(gallery)
	if err != nil {
		vd.SetAlert(err)
	} else {
		vd.Alert = &views.Alert{
			Level:   views.AlertLvlSuccess,
			Message: msgSuccessUpdate,
		}
	}
	vd.Yield = gallery
	g.EditView.Render(w, vd)
}

// ------ Helper ------
func (g *Galleries) getGalleryById(w http.ResponseWriter, r *http.Request) (*models.Gallery, error) {
	// Get :id from url id param, converted from string to int
	vars := mux.Vars(r)
	idParam := vars["id"]
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid gallery ID", http.StatusNotFound)
		return nil, err
	}

	gallery, err := g.gs.GetOneById(uint(id))
	if err != nil {
		switch err {
		case models.ErrNotFound:
			http.Error(w, "Gallery not found", http.StatusNotFound)
		default:
			http.Error(w, "Whoops! Something went wrong.", http.StatusInternalServerError)
		}
		return nil, err
	}

	return gallery, nil
}
