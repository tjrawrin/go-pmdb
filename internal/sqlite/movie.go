package sqlite

import (
	"database/sql"

	"../service"
)

// MovieService represents a SQLite implementation of a MovieService.
type MovieService struct {
	DB *sql.DB
}

// GetMovies returns all movies from the database.
func (s *MovieService) GetMovies() (*service.Movies, error) {
	rows, err := s.DB.Query(`SELECT * FROM movies;`)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	var movies service.Movies
	for rows.Next() {
		var movie service.Movie
		if err := rows.Scan(&movie.ID, &movie.Title); err != nil {
			return nil, err
		}
		movies = append(movies, &movie)
	}

	return &movies, nil
}

// GetMovie returns a single movie from the database.
func (s *MovieService) GetMovie(id int64) (*service.Movie, error) {
	row := s.DB.QueryRow(`SELECT id, title FROM movies WHERE id = $1;`, id)
	var movie service.Movie
	if err := row.Scan(&movie.ID, &movie.Title); err != nil {
		return nil, err
	}

	return &movie, nil
}

// CreateMovie adds a new movie to the database.
func (s *MovieService) CreateMovie(movie *service.Movie) (int64, error) {
	res, err := s.DB.Exec(`INSERT INTO movies (title) VALUES ($1);`, movie.Title)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

// UpdateMovie updates an existing movie in the database.
func (s *MovieService) UpdateMovie(id int64, movie *service.Movie) error {
	_, err := s.DB.Exec(`UPDATE movies SET id = $1, title = $2 WHERE id = $1;`, id, movie.Title)
	if err != nil {
		return err
	}

	return nil
}

// DeleteMovie removes an existing movie from the database.
func (s *MovieService) DeleteMovie(id int64) error {
	_, err := s.DB.Exec(`DELETE FROM movies WHERE id = $1;`, id)
	if err != nil {
		return err
	}

	return nil
}
