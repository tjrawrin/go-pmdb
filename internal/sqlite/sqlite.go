package sqlite

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3" // ...
)

// Start attempts to open a SQLite database or returns an error
// if opening fails. It also pings the database to test the connection
// or returns an error if a connection cannot be made. It then
// creates any tables or returns an error if table creation fails.
func Start(dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	err = initialize(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// initialize creates the database tables. This is non-destructive if
// data already exists.
func initialize(db *sql.DB) error {
	dbTx, err := db.Begin()
	if err != nil {
		dbTx.Rollback()
		return err
	}

	// Create the movies table.
	if err = moviesTable(dbTx); err != nil {
		dbTx.Rollback()
		return err
	}

	return dbTx.Commit()
}

// moviesTable defines and creates a new movies database table if
// one doesn't already exist.
func moviesTable(db *sql.Tx) error {
	stmt := `
		CREATE TABLE IF NOT EXISTS movies(
			id INTEGER PRIMARY KEY,
			title TEXT
		);
	`

	_, err := db.Exec(stmt)

	return err
}
