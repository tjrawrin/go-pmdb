package http

import (
	"log"
	"net/http"
	"strconv"

	"../render"
	"../service"
	"../sqlite"
	"github.com/go-chi/chi"
)

// MovieHandler ...
type MovieHandler struct {
	MovieService *sqlite.MovieService
}

// Routes creates a REST router for the page handler.
func (h *MovieHandler) Routes() chi.Router {
	r := chi.NewRouter()

	// Load middleware specific to this router.
	// r.Use()

	r.Get("/", h.index)
	r.Get("/new", h.new)
	r.Post("/", h.create)
	r.Get("/{id}", h.show)
	r.Get("/{id}/edit", h.edit)
	r.Put("/{id}", h.update)
	r.Post("/{id}", h.delete)

	return r
}

// Index responds to a request for a list of movies.
func (h *MovieHandler) index(w http.ResponseWriter, r *http.Request) {
	// Call GetMovies to retrieve all movies from the database.
	if movies, err := h.MovieService.GetMovies(); err != nil {
		// Render an error response and set status code.
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error:", err)
	} else {
		// Render a HTML response and set status code.
		render.HTML(w, http.StatusOK, "movie/index.html", movies)
	}
}

// New responds to a request for entering details for a movie.
func (h *MovieHandler) new(w http.ResponseWriter, r *http.Request) {
	// Render a HTML response and set status code.
	render.HTML(w, http.StatusOK, "movie/new.html", nil)
}

// Create responds to a request for adding a movie.
func (h *MovieHandler) create(w http.ResponseWriter, r *http.Request) {
	// Parse the page form values.
	err := r.ParseForm()
	if err != nil {
		// Render an error response and set status code.
		http.Error(w, "Unprocessable Entity", http.StatusUnprocessableEntity)
		log.Println("Error:", err)
		return
	}

	// Create a temporary movie struct to unmarshal the request body into.
	movie := &service.Movie{
		Title:  r.FormValue("title"),
		ImdbID: r.FormValue("imdb_id"),
	}

	// Call the CreateMovie to add the new movie to the database.
	id, err := h.MovieService.CreateMovie(movie)
	if err != nil {
		// Render an error response and set status code.
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error:", err)
		return
	}

	// Call GetMovie to get the movie from the database.
	if _, err := h.MovieService.GetMovie(id); err != nil {
		// Render an error response and set status code.
		http.Error(w, "Not Found", http.StatusNotFound)
		log.Println("Error:", err)
	} else {
		http.Redirect(w, r, "/movies/"+strconv.FormatInt(id, 10), http.StatusCreated)
		return
	}
}

// Show responds to a request for a single movie.
func (h *MovieHandler) show(w http.ResponseWriter, r *http.Request) {
	// Parse the id param from the URL and convert it into an int64.
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		// Render an error response and set status code.
		http.Error(w, "Not Found", http.StatusNotFound)
		log.Println("Error:", err)
		return
	}

	// Call GetMovie to get the movie from the database.
	if movie, err := h.MovieService.GetMovie(id); err != nil {
		// Render an error response and set status code.
		http.Error(w, "Not Found", http.StatusNotFound)
		log.Println("Error:", err)
	} else {
		// Render a HTML response and set status code.
		render.HTML(w, http.StatusOK, "movie/show.html", movie)
	}
}

// Edit responds to a request for entering details for a movie.
func (h *MovieHandler) edit(w http.ResponseWriter, r *http.Request) {
	// Parse the id param from the URL and convert it into an int64.
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		// Render an error response and set status code.
		http.Error(w, "Not Found", http.StatusNotFound)
		log.Println("Error:", err)
		return
	}

	// Call GetMovie to get the movie from the database.
	if movie, err := h.MovieService.GetMovie(id); err != nil {
		// Render an error response and set status code.
		http.Error(w, "Not Found", http.StatusNotFound)
		log.Println("Error:", err)
	} else {
		// Render a HTML response and set status code.
		render.HTML(w, http.StatusOK, "movie/edit.html", movie)
	}
}

// Update responds to a request for updating a movie.
func (h *MovieHandler) update(w http.ResponseWriter, r *http.Request) {
	// Parse the id param from the URL and convert it into an int64.
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		// Render an error response and set status code.
		http.Error(w, "Not Found", http.StatusNotFound)
		log.Println("Error:", err)
		return
	}

	// Call GetMovie to get the movie from the database.
	if _, err := h.MovieService.GetMovie(id); err != nil {
		// Render an error response and set status code.
		http.Error(w, "Not Found", http.StatusNotFound)
		log.Println("Error:", err)
		return
	}

	// Parse the page form values.
	err = r.ParseForm()
	if err != nil {
		// Render an error response and set status code.
		http.Error(w, "Unprocessable Entity", http.StatusUnprocessableEntity)
		log.Println("Error:", err)
		return
	}

	// Create a temporary movie struct to unmarshal the request body into.
	movie := &service.Movie{
		Title:  r.FormValue("title"),
		ImdbID: r.FormValue("imdb_id"),
	}

	// Call UpdateMovie to update the movie in the database.
	err = h.MovieService.UpdateMovie(id, movie)
	if err != nil {
		// Render an error response and set status code.
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error:", err)
		return
	}

	// Call GetMovie to get the movie from the database.
	if _, err := h.MovieService.GetMovie(id); err != nil {
		// Render an error response and set status code.
		http.Error(w, "Not Found", http.StatusNotFound)
		log.Println("Error:", err)
	} else {
		http.Redirect(w, r, "/movies/"+strconv.FormatInt(id, 10), http.StatusCreated) // TODO(tim): FIX THIS
		return
	}
}

// Delete responds to a request for removing a movie.
func (h *MovieHandler) delete(w http.ResponseWriter, r *http.Request) {
	// Parse the id param from the URL and convert it into an int64.
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		// Render an error response and set status code.
		http.Error(w, "Not Found", http.StatusNotFound)
		log.Println("Error:", err)
		return
	}

	// Call GetMovie to get the movie from the database.
	if _, err := h.MovieService.GetMovie(id); err != nil {
		// Render an error response and set status code.
		http.Error(w, "Not Found", http.StatusNotFound)
		log.Println("Error:", err)
		return
	}

	// Call DeleteMovie to remove the movie from the database.
	if err = h.MovieService.DeleteMovie(id); err != nil {
		// Render an error response and set status code.
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error:", err)
	} else {
		http.Redirect(w, r, "/movies", http.StatusSeeOther) // TODO(tim): FIX THIS
		return
	}
}
