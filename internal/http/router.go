package http

import (
	"./api"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// Router ...
type Router struct {
	APIMovieHandler *api.MovieHandler
	MovieHandler    *MovieHandler
	PageHandler     *PageHandler
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
	router.Mount("/movies", r.MovieHandler.Routes())

	// API (v1) routes
	router.Route("/api/v1", func(sr chi.Router) {
		sr.Mount("/movies", r.APIMovieHandler.Routes())
	})

	return router
}
