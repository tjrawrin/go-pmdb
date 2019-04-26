package http

import (
	"net/http"

	"../render"
	"github.com/go-chi/chi"
)

// PageHandler ...
type PageHandler struct{}

// Routes creates a REST router for the page handler.
func (h *PageHandler) Routes() chi.Router {
	r := chi.NewRouter()

	// Load middleware specific to this router.
	// r.Use()

	r.Get("/*", h.index)

	return r
}

// Index responds to a request for the site index page.
func (h *PageHandler) index(w http.ResponseWriter, r *http.Request) {
	render.HTML(w, http.StatusOK, "page/index.html", nil)
}
