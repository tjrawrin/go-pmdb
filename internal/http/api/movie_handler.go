package api

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"../../render"
	"../../service"
	"../../sqlite"
	"github.com/go-chi/chi"
)

// MovieHandler ...
type MovieHandler struct {
	MovieService *sqlite.MovieService
}

// Routes creates a REST router for the movie handler.
func (h *MovieHandler) Routes() chi.Router {
	r := chi.NewRouter()

	// Load middleware specific to this router.
	// r.Use()

	r.Get("/", h.index)
	r.Post("/", h.create)
	r.Get("/{id}", h.show)
	r.Put("/{id}", h.update)
	r.Delete("/{id}", h.delete)

	return r
}

// Index responds to a request for a list of movies.
func (h *MovieHandler) index(w http.ResponseWriter, r *http.Request) {
	// Call GetMovies to retrieve all movies from the database.
	if movies, err := h.MovieService.GetMovies(); err != nil {
		// Render a JSON response and set status code.
		render.JSON(w, http.StatusInternalServerError,
			map[string]string{
				"error":   "Internal Server Error",
				"message": err.Error(),
			})
		log.Println("Error:", err)
	} else {
		// If the movies slice does not return nil. Respond with the movies,
		// otherwise respond with an empty slice.
		if *movies != nil {
			// Render a JSON response and set status code.
			render.JSON(w, http.StatusOK, movies)
		} else {
			// Render a JSON response and set status code.
			render.JSON(w, http.StatusOK, []string{})
		}
	}
}

// Create responds to a request for adding a movie.
func (h *MovieHandler) create(w http.ResponseWriter, r *http.Request) {
	// Read the request body (limited to 1048576 bytes).
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	defer r.Body.Close()
	if err != nil {
		// Render a JSON response and set status code.
		render.JSON(w, http.StatusUnprocessableEntity,
			map[string]string{
				"error":   "Unprocessable Entity",
				"message": err.Error(),
			})
		log.Println("Error:", err)
		return
	}

	// Create a temporary movie struct to unmarshal the request body into.
	var movie *service.Movie
	err = json.Unmarshal(body, &movie)
	if err != nil {
		// Render a JSON response and set status code.
		render.JSON(w, http.StatusUnprocessableEntity,
			map[string]string{
				"error":   "Unprocessable Entity",
				"message": err.Error(),
			})
		log.Println("Error:", err)
		return
	}

	// Call the CreateMovie to add the new movie to the database.
	id, err := h.MovieService.CreateMovie(movie)
	if err != nil {
		// Render a JSON response and set status code.
		render.JSON(w, http.StatusInternalServerError,
			map[string]string{
				"error":   "Internal Server Error",
				"message": err.Error(),
			})
		log.Println("Error:", err)
		return
	}

	// Call GetMovie to get the movie from the database.
	if movie, err := h.MovieService.GetMovie(id); err != nil {
		// Render a JSON response and set status code.
		render.JSON(w, http.StatusNotFound,
			map[string]string{
				"error":   "Not Found",
				"message": err.Error(),
			})
		log.Println("Error:", err)
	} else {
		// Render a JSON response and set status code.
		render.JSON(w, http.StatusCreated, movie)
	}
}

// Show responds to a request for a single movie.
func (h *MovieHandler) show(w http.ResponseWriter, r *http.Request) {
	// Parse the id param from the URL and convert it into an int64.
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		// Render a JSON response and set status code.
		render.JSON(w, http.StatusNotFound,
			map[string]string{
				"error":   "Not Found",
				"message": err.Error(),
			})
		log.Println("Error:", err)
		return
	}

	// Call GetMovie to get the movie from the database.
	if movie, err := h.MovieService.GetMovie(id); err != nil {
		// Render a JSON response and set status code.
		render.JSON(w, http.StatusNotFound,
			map[string]string{
				"error":   "Not Found",
				"message": err.Error(),
			})
		log.Println("Error:", err)
	} else {
		// Render a JSON response and set status code.
		render.JSON(w, http.StatusOK, movie)
	}
}

// Update responds to a request for updating a movie.
func (h *MovieHandler) update(w http.ResponseWriter, r *http.Request) {
	// Parse the id param from the URL and convert it into an int64.
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		// Render a JSON response and set status code.
		render.JSON(w, http.StatusNotFound,
			map[string]string{
				"error":   "Not Found",
				"message": err.Error(),
			})
		log.Println("Error:", err)
		return
	}

	// Call GetMovie to get the movie from the database.
	if _, err := h.MovieService.GetMovie(id); err != nil {
		// Render a JSON response and set status code.
		render.JSON(w, http.StatusNotFound,
			map[string]string{
				"error":   "Not Found",
				"message": err.Error(),
			})
		log.Println("Error:", err)
		return
	}

	// Read the request body (limited to 1048576 bytes).
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	defer r.Body.Close()
	if err != nil {
		// Render a JSON response and set status code.
		render.JSON(w, http.StatusUnprocessableEntity,
			map[string]string{
				"error":   "Unprocessable Entity",
				"message": err.Error(),
			})
		log.Println("Error:", err)
		return
	}

	// Create a temporary movie struct to unmarshal the request body into.
	var movie *service.Movie
	err = json.Unmarshal(body, &movie)
	if err != nil {
		// Render a JSON response and set status code.
		render.JSON(w, http.StatusUnprocessableEntity,
			map[string]string{
				"error":   "Unprocessable Entity",
				"message": err.Error(),
			})
		log.Println("Error:", err)
		return
	}

	// Call UpdateMovie to update the movie in the database.
	err = h.MovieService.UpdateMovie(id, movie)
	if err != nil {
		// Render a JSON response and set status code.
		render.JSON(w, http.StatusInternalServerError,
			map[string]string{
				"error":   "Internal Server Error",
				"message": err.Error(),
			})
		log.Println("Error:", err)
		return
	}

	// Call GetMovie to get the movie from the database.
	if movie, err := h.MovieService.GetMovie(id); err != nil {
		// Render a JSON response and set status code.
		render.JSON(w, http.StatusNotFound,
			map[string]string{
				"error":   "Not Found",
				"message": err.Error(),
			})
		log.Println("Error:", err)
	} else {
		// Render a JSON response and set status code.
		render.JSON(w, http.StatusCreated, movie)
	}
}

// Delete responds to a request for removing a movie.
func (h *MovieHandler) delete(w http.ResponseWriter, r *http.Request) {
	// Parse the id param from the URL and convert it into an int64.
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		// Render a JSON response and set status code.
		render.JSON(w, http.StatusNotFound,
			map[string]string{
				"error":   "Not Found",
				"message": err.Error(),
			})
		log.Println("Error:", err)
		return
	}

	// Call GetMovie to get the movie from the database.
	if _, err := h.MovieService.GetMovie(id); err != nil {
		// Render a JSON response and set status code.
		render.JSON(w, http.StatusNotFound,
			map[string]string{
				"error":   "Not Found",
				"message": err.Error(),
			})
		log.Println("Error:", err)
		return
	}

	// Call DeleteMovie to remove the movie from the database.
	if err = h.MovieService.DeleteMovie(id); err != nil {
		// Render a JSON response and set status code.
		render.JSON(w, http.StatusInternalServerError,
			map[string]string{
				"error":   "Internal Server Error",
				"message": err.Error(),
			})
		log.Println("Error:", err)
	} else {
		// Render a JSON response and set status code.
		render.JSON(w, http.StatusOK, map[string]string{})
	}
}
