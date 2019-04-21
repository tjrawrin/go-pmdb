package http

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

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
	// Set the header Content-Type.
	w.Header().Set("Content-Type", "application/json")

	// Call GetMovies to retrieve all movies from the database.
	if movies, err := m.MovieService.GetMovies(); err != nil {
		// Set the status to 500.
		w.WriteHeader(http.StatusInternalServerError)
		// Respond with error json message.
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal Server Error"})
	} else {
		// Set the status to 200.
		w.WriteHeader(http.StatusOK)

		// If the movies slice does not return nil. Respond with the movies,
		// otherwise respond with an empty slice.
		if *movies != nil {
			json.NewEncoder(w).Encode(movies)
		} else {
			json.NewEncoder(w).Encode([]map[string]string{})
		}
	}
}

// Create responds to a request for adding a movie.
func (m *MovieHandler) Create(w http.ResponseWriter, r *http.Request) {
	// Set the header Content-Type.
	w.Header().Set("Content-Type", "application/json")

	// Read the request body (limited to 1048576 bytes).
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	defer r.Body.Close()
	if err != nil {
		// Set the status to 422.
		w.WriteHeader(http.StatusUnprocessableEntity)
		// Respond with error json message.
		json.NewEncoder(w).Encode(map[string]string{"error": "Unprocessable Entity"})
		return
	}

	// Create a temporary movie struct to unmarshal the request body into.
	var movie *service.Movie
	err = json.Unmarshal(body, &movie)
	if err != nil {
		// Set the status to 422.
		w.WriteHeader(http.StatusUnprocessableEntity)
		// Respond with error json message.
		json.NewEncoder(w).Encode(map[string]string{"error": "Unprocessable Entity"})
		return
	}

	// Call the CreateMovie to add the new movie to the database.
	id, err := m.MovieService.CreateMovie(movie)
	if err != nil {
		// Set the status to 500.
		w.WriteHeader(http.StatusInternalServerError)
		// Respond with error json message.
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal Server Error"})
		return
	}

	// Call GetMovie to get the movie from the database.
	if movie, err := m.MovieService.GetMovie(id); err != nil {
		// Set the status to 404.
		w.WriteHeader(http.StatusNotFound)
		// Respond with error json message.
		json.NewEncoder(w).Encode(map[string]string{"error": "Not Found"})
	} else {
		// Set the status to 201.
		w.WriteHeader(http.StatusCreated)
		// Respond with the movie.
		json.NewEncoder(w).Encode(movie)
	}
}

// Show responds to a request for a single movie.
func (m *MovieHandler) Show(w http.ResponseWriter, r *http.Request) {
	// Set the header Content-Type.
	w.Header().Set("Content-Type", "application/json")

	// Parse the id param from the URL and convert it into an int64.
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		// Set the status to 404.
		w.WriteHeader(http.StatusNotFound)
		// Respond with error json message.
		json.NewEncoder(w).Encode(map[string]string{"error": "Not Found"})
		return
	}

	// Call GetMovie to get the movie from the database.
	if movie, err := m.MovieService.GetMovie(id); err != nil {
		// Set the status to 404.
		w.WriteHeader(http.StatusNotFound)
		// Respond with error json message.
		json.NewEncoder(w).Encode(map[string]string{"error": "Not Found"})
	} else {
		// Set the status to 200.
		w.WriteHeader(http.StatusOK)
		// Respond with the movie.
		json.NewEncoder(w).Encode(movie)
	}
}

// Update responds to a request for updating a movie.
func (m *MovieHandler) Update(w http.ResponseWriter, r *http.Request) {
	// Set the header Content-Type.
	w.Header().Set("Content-Type", "application/json")

	// Parse the id param from the URL and convert it into an int64.
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		// Set the status to 404.
		w.WriteHeader(http.StatusNotFound)
		// Respond with error json message.
		json.NewEncoder(w).Encode(map[string]string{"error": "Not Found"})
		return
	}

	// Call GetMovie to get the movie from the database.
	if _, err := m.MovieService.GetMovie(id); err != nil {
		// Set the status to 404.
		w.WriteHeader(http.StatusNotFound)
		// Respond with error json message.
		json.NewEncoder(w).Encode(map[string]string{"error": "Not Found"})
		return
	}

	// Read the request body (limited to 1048576 bytes).
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	defer r.Body.Close()
	if err != nil {
		// Set the status to 422.
		w.WriteHeader(http.StatusUnprocessableEntity)
		// Respond with error json message.
		json.NewEncoder(w).Encode(map[string]string{"error": "Unprocessable Entity"})
		return
	}

	// Create a temporary movie struct to unmarshal the request body into.
	var movie *service.Movie
	err = json.Unmarshal(body, &movie)
	if err != nil {
		// Set the status to 422.
		w.WriteHeader(http.StatusUnprocessableEntity)
		// Respond with error json message.
		json.NewEncoder(w).Encode(map[string]string{"error": "Unprocessable Entity"})
		return
	}

	// Call UpdateMovie to update the movie in the database.
	err = m.MovieService.UpdateMovie(id, movie)
	if err != nil {
		// Set the status to 500.
		w.WriteHeader(http.StatusInternalServerError)
		// Respond with error json message.
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal Server Error"})
		return
	}

	// Call GetMovie to get the movie from the database.
	if movie, err := m.MovieService.GetMovie(id); err != nil {
		// Set the status to 404.
		w.WriteHeader(http.StatusNotFound)
		// Respond with error json message.
		json.NewEncoder(w).Encode(map[string]string{"error": "Not Found"})
	} else {
		// Set the status to 201.
		w.WriteHeader(http.StatusCreated)
		// Respond with the movie.
		json.NewEncoder(w).Encode(movie)
	}
}

// Delete responds to a request for removing a movie.
func (m *MovieHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// Set the header Content-Type.
	w.Header().Set("Content-Type", "application/json")

	// Parse the id param from the URL and convert it into an int64.
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		// Set the status to 404.
		w.WriteHeader(http.StatusNotFound)
		// Respond with error json message.
		json.NewEncoder(w).Encode(map[string]string{"error": "Not Found"})
		return
	}

	// Call GetMovie to get the movie from the database.
	if _, err := m.MovieService.GetMovie(id); err != nil {
		// Set the status to 404.
		w.WriteHeader(http.StatusNotFound)
		// Respond with error json message.
		json.NewEncoder(w).Encode(map[string]string{"error": "Not Found"})
		return
	}

	// Call DeleteMovie to remove the movie from the database.
	if err = m.MovieService.DeleteMovie(id); err != nil {
		// Set the status to 500.
		w.WriteHeader(http.StatusInternalServerError)
		// Respond with error json message.
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal Server Error"})
	} else {
		// Set the status to 200.
		w.WriteHeader(http.StatusOK)
		// Respond with an empty object.
		json.NewEncoder(w).Encode(map[string]string{})
	}
}
