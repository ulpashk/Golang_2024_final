package data

import (
	"goproject/internal/validator"
	"context"
	"database/sql"
	"errors"
	"time"
	"fmt"

	//"github.com/lib/pq"
)


type Song struct{
	Id 			int		`json:"id"`
	Title 		string	`json:"title"`
	Length 		int		`json:"length"`
	Album_id	int 	`json:"albumId"`
} 

func ValidateSong(v *validator.Validator, song *Song){
	v.Check(song.Title != "", "title", "must be provided")
	// v.Check(song.Id > 0, "id", "must be greater than 0")
	v.Check(song.Album_id != 0, "albumId", "must be greater than 0")
}


type SongModel struct {
	DB *sql.DB
}


func (s SongModel) Insert(song *Song) error {
	// Define the SQL query for inserting a new record in the movies table and returning
	// the system-generated data.
	query := `
		INSERT INTO songs(title, length, album_id)
		VALUES ($1, $2, $3)
		RETURNING song_id, title, length, album_id;`

	// Create an args slice containing the values for the placeholder parameters from
	// the movie struct. Declaring this slice immediately next to our SQL query helps to
	// make it nice and clear *what values are being used where* in the query.
	args := []interface{}{song.Title, song.Length, song.Album_id}

	// Use the QueryRow() method to execute the SQL query on our connection pool,
	// passing in the args slice as a variadic parameter and scanning the system-
	// generated id, created_at and version values into the movie struct.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return s.DB.QueryRowContext(ctx, query, args...).Scan(&song.Id, &song.Title, &song.Length, &song.Album_id)
}

// Add a placeholder method for fetching a specific record from the movies table.
func (s SongModel) Get(id int64) (*Song, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	// Retrieve a specific menu item based on its ID.
	query := `
		SELECT song_id, title, length, album_id
		FROM songs
		WHERE song_id = $1;`

	var song Song
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := s.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&song.Id, &song.Title, &song.Length, &song.Album_id)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &song, nil
}

// Add a placeholder method for updating a specific record in the movies table.
func (s SongModel) Update(song *Song) error {
	query := `
		UPDATE songs
		SET title = $1, length = $2, album_id = $3
		WHERE song_id = $4
		RETURNING song_id, title, length, album_id;
		`

	args := []interface{}{song.Title, song.Length, song.Album_id, song.Id}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return s.DB.QueryRowContext(ctx, query, args...).Scan(&song.Id, &song.Title, &song.Length, &song.Album_id)
}


func (s SongModel) Delete(id int64) error {
	// Return an ErrRecordNotFound error if the movie ID is less than 1.
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM songs
		WHERE song_id = $1
		`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := s.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	// Call the RowsAffected() method on the sql.Result object to get the number of rows
	// affected by the query.
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	// If no rows were affected, we know that the movies table didn't contain a record
	// with the provided ID at the moment we tried to delete it. In that case we
	// return an ErrRecordNotFound error.
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}


// Create a new GetAll() method which returns a slice of movies. Although we're not
// using them right now, we've set this up to accept the various filter parameters as
// arguments.
func (s SongModel) GetAll(title string, length int, filters Filters) ([]*Song, Metadata, error) {
	// Construct the SQL query to retrieve all movie records.
	query :=  fmt.Sprintf(`
		SELECT count(*) OVER(), song_id, title, length, album_id
		FROM songs
		WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '')
		AND (length = $2 OR $2 = 1)
		ORDER BY %s %s, song_id
		LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())
	
/*SELECT count(*) OVER(), song_id, title, length, album_id
		FROM song
		WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '')
		AND (length = $2 OR $2 = '{}')
		ORDER BY %s %s, id ASC
		LIMIT $3 OFFSET $4`
		WHERE (title ILIKE $1 OR $1 = '')
		*/

	// Create a context with a 3-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// As our SQL query now has quite a few placeholder parameters, let's collect the
	// values for the placeholders in a slice. Notice here how we call the limit() and
	// offset() methods on the Filters struct to get the appropriate values for the
	// LIMIT and OFFSET clauses.
	args := []interface{}{title, length, filters.limit(), filters.offset()}

	// Use QueryContext() to execute the query. This returns a sql.Rows resultset
	// containing the result.
	rows, err := s.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	// Importantly, defer a call to rows.Close() to ensure that the resultset is closed
	// before GetAll() returns.
	defer rows.Close()
	// Declare a totalRecords variable.
	totalRecords := 0
	// Initialize an empty slice to hold the movie data.
	songs := []*Song{}

	// Use rows.Next to iterate through the rows in the resultset.
	for rows.Next() {
		// Initialize an empty Movie struct to hold the data for an individual movie.
		var song Song
		// Scan the values from the row into the Movie struct. Again, note that we're
		// using the pq.Array() adapter on the genres field here.
		err := rows.Scan(
			&totalRecords, // Scan the count from the window function into totalRecords.
			&song.Id,
			&song.Title,
			&song.Length,
			&song.Album_id,
		)
		if err != nil {
			return nil, Metadata{}, err
		}
		// Add the Movie struct to the slice.
		songs = append(songs, &song)
	}
	// When the rows.Next() loop has finished, call rows.Err() to retrieve any error
	// that was encountered during the iteration.
	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	// Generate a Metadata struct, passing in the total record count and pagination
	// parameters from the client.
	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)
	// If everything went OK, then return the slice of movies.
	return songs, metadata, nil
}
