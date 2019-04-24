package service

import "time"

// Movie is a struct containing information about a movie.
type Movie struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	ImdbID    string    `json:"imdbId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Movies is a slice of movie structs.
type Movies []*Movie

// MovieService contains function signatures for implementing a movie service.
type MovieService interface {
	GetMovies() (*Movies, error)
	GetMovie(id int64) (*Movie, error)
	CreateMovie(m *Movie) error
	UpdateMovie(id int64, m *Movie) error
	DeleteMovie(id int64) error
}
