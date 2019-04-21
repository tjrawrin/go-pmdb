package http

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// Router ...
type Router struct {
	MovieHandler *MovieHandler
	PageHandler  *PageHandler
}

// Router ...
func (r *Router) Router() *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.DefaultCompress)

	// Non-API routes
	router.Mount("/", r.PageHandler.Routes())

	// API (v1) routes
	router.Route("/api/v1", func(sr chi.Router) {
		sr.Mount("/movies", r.MovieHandler.Routes())
	})

	return router
}
