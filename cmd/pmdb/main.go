package main

import (
	"log"

	"../../internal/http"
	"../../internal/sqlite"
)

func main() {
	// Start database.
	db, err := sqlite.Start("./web/data/pmdb.db")
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}

	// Create services.
	movieService := &sqlite.MovieService{DB: db}

	// Init handlers and attach services to handlers if necessary.
	movieHandler := &http.MovieHandler{MovieService: movieService}
	pageHandler := &http.PageHandler{}

	// Attach handlers to router.
	router := &http.Router{
		MovieHandler: movieHandler,
		PageHandler:  pageHandler,
	}

	// Create a server.
	srv := &http.Server{Router: router.Router()}

	// Run the server.
	log.Fatal(srv.Run())
}
