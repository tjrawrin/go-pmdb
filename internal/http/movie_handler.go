package http

import (
	"encoding/json"
	"io"
	"io/ioutil"
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

// Routes creates a REST router for the movie handler.
func (m *MovieHandler) Routes() chi.Router {
	r := chi.NewRouter()

	// Load middleware specific to this router.
	// r.Use()

	r.Get("/", m.Index)
	r.Post("/", m.Create)
	r.Get("/{id}", m.Show)
	r.Put("/{id}", m.Update)
	r.Delete("/{id}", m.Delete)

	return r
}

// Index responds to a request for a list of movies.
func (m *MovieHandler) Index(w http.ResponseWriter, r *http.Request) {
	// Call GetMovies to retrieve all movies from the database.
	if movies, err := m.MovieService.GetMovies(); err != nil {
		// Render a JSON response and set status code.
		render.JSON(w, http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
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
func (m *MovieHandler) Create(w http.ResponseWriter, r *http.Request) {
	// Read the request body (limited to 1048576 bytes).
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	defer r.Body.Close()
	if err != nil {
		// Render a JSON response and set status code.
		render.JSON(w, http.StatusUnprocessableEntity, map[string]string{"error": "Unprocessable Entity"})
		return
	}

	// Create a temporary movie struct to unmarshal the request body into.
	var movie *service.Movie
	err = json.Unmarshal(body, &movie)
	if err != nil {
		// Render a JSON response and set status code.
		render.JSON(w, http.StatusUnprocessableEntity, map[string]string{"error": "Unprocessable Entity"})
		return
	}

	// Call the CreateMovie to add the new movie to the database.
	id, err := m.MovieService.CreateMovie(movie)
	if err != nil {
		// Render a JSON response and set status code.
		render.JSON(w, http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
		return
	}

	// Call GetMovie to get the movie from the database.
	if movie, err := m.MovieService.GetMovie(id); err != nil {
		// Render a JSON response and set status code.
		render.JSON(w, http.StatusNotFound, map[string]string{"error": "Not Found"})
	} else {
		// Render a JSON response and set status code.
		render.JSON(w, http.StatusCreated, movie)
	}
}

// Show responds to a request for a single movie.
func (m *MovieHandler) Show(w http.ResponseWriter, r *http.Request) {
	// Parse the id param from the URL and convert it into an int64.
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		// Render a JSON response and set status code.
		render.JSON(w, http.StatusNotFound, map[string]string{"error": "Not Found"})
		return
	}

	// Call GetMovie to get the movie from the database.
	if movie, err := m.MovieService.GetMovie(id); err != nil {
		// Render a JSON response and set status code.
		render.JSON(w, http.StatusNotFound, map[string]string{"error": "Not Found"})
	} else {
		// Render a JSON response and set status code.
		render.JSON(w, http.StatusOK, movie)
	}
}

// Update responds to a request for updating a movie.
func (m *MovieHandler) Update(w http.ResponseWriter, r *http.Request) {
	// Parse the id param from the URL and convert it into an int64.
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		// Render a JSON response and set status code.
		render.JSON(w, http.StatusNotFound, map[string]string{"error": "Not Found"})
		return
	}

	// Call GetMovie to get the movie from the database.
	if _, err := m.MovieService.GetMovie(id); err != nil {
		// Render a JSON response and set status code.
		render.JSON(w, http.StatusNotFound, map[string]string{"error": "Not Found"})
		return
	}

	// Read the request body (limited to 1048576 bytes).
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	defer r.Body.Close()
	if err != nil {
		// Render a JSON response and set status code.
		render.JSON(w, http.StatusUnprocessableEntity, map[string]string{"error": "Unprocessable Entity"})
		return
	}

	// Create a temporary movie struct to unmarshal the request body into.
	var movie *service.Movie
	err = json.Unmarshal(body, &movie)
	if err != nil {
		// Render a JSON response and set status code.
		render.JSON(w, http.StatusUnprocessableEntity, map[string]string{"error": "Unprocessable Entity"})
		return
	}

	// Call UpdateMovie to update the movie in the database.
	err = m.MovieService.UpdateMovie(id, movie)
	if err != nil {
		// Render a JSON response and set status code.
		render.JSON(w, http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
		return
	}

	// Call GetMovie to get the movie from the database.
	if movie, err := m.MovieService.GetMovie(id); err != nil {
		// Render a JSON response and set status code.
		render.JSON(w, http.StatusNotFound, map[string]string{"error": "Not Found"})
	} else {
		// Render a JSON response and set status code.
		render.JSON(w, http.StatusCreated, movie)
	}
}

// Delete responds to a request for removing a movie.
func (m *MovieHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// Parse the id param from the URL and convert it into an int64.
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		// Render a JSON response and set status code.
		render.JSON(w, http.StatusNotFound, map[string]string{"error": "Not Found"})
		return
	}

	// Call GetMovie to get the movie from the database.
	if _, err := m.MovieService.GetMovie(id); err != nil {
		// Render a JSON response and set status code.
		render.JSON(w, http.StatusNotFound, map[string]string{"error": "Not Found"})
		return
	}

	// Call DeleteMovie to remove the movie from the database.
	if err = m.MovieService.DeleteMovie(id); err != nil {
		// Render a JSON response and set status code.
		render.JSON(w, http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
	} else {
		// Render a JSON response and set status code.
		render.JSON(w, http.StatusOK, map[string]string{})
	}
}
