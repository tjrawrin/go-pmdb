package http

import (
	"net/http"

	"github.com/go-chi/chi"
)

// PageHandler ...
type PageHandler struct{}

// Routes creates a REST router for the page handler.
func (p *PageHandler) Routes() chi.Router {
	r := chi.NewRouter()

	// Load middleware specific to this router.
	// r.Use()

	r.Get("/*", p.index)

	return r
}

// Index responds to a request for the site index page.
func (p *PageHandler) index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	w.Write([]byte("<h1>Page#Index</h1>"))
}
