package data

import (
	"goproject/internal/validator"
	"context"
	"database/sql"
	"errors"
	"time"
	"fmt"

)


type Album struct{
	Id 			int		`json:"id"`
	Title 		string	`json:"title"`
	Genre    	string  `json:"genre"`
	Tracks 		int		`json:"tracks"`
	Group_id	int 	`json:"groupId"`
} 

func ValidateAlbum(v *validator.Validator, album *Album){
	v.Check(album.Title != "", "title", "must be provided")
	v.Check(album.Genre != "", "genre", "must be provided")
	v.Check(album.Tracks != 0, "tracks", "must be greater than 0")
	v.Check(album.Group_id != 0, "groupId", "must be greater than 0")
}

type AlbumModel struct {
	DB *sql.DB
}

func (a AlbumModel) Insert(album *Album) error {
	// Define the SQL query for inserting a new record in the movies table and returning
	// the system-generated data.
	query := `
		INSERT INTO album(title, genre, tracks, group_id)
		VALUES ($1, $2, $3, $4)
		RETURNING album_id, title, genre, tracks, group_id;`

	// Create an args slice containing the values for the placeholder parameters from
	// the movie struct. Declaring this slice immediately next to our SQL query helps to
	// make it nice and clear *what values are being used where* in the query.
	args := []interface{}{album.Title, album.Genre, album.Tracks, album.Group_id}

	// Use the QueryRow() method to execute the SQL query on our connection pool,
	// passing in the args slice as a variadic parameter and scanning the system-
	// generated id, created_at and version values into the movie struct.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return a.DB.QueryRowContext(ctx, query, args...).Scan(&album.Id, &album.Title, &album.Genre, &album.Tracks, &album.Group_id)
}

// Add a placeholder method for fetching a specific record from the movies table.
func (a AlbumModel) Get(id int64) (*Album, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	// Retrieve a specific menu item based on its ID.
	query := `
		SELECT album_id, title, genre, tracks, group_id
		FROM album
		WHERE album_id = $1;`

	var album Album
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := a.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&album.Id, &album.Title, &album.Genre, &album.Tracks, &album.Group_id)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &album, nil
}

// Add a placeholder method for updating a specific record in the movies table.
func (a AlbumModel) Update(album *Album) error {
	query := `
		UPDATE album
		SET title = $1, genre = $2, tracks = $3, group_id = $4
		WHERE album_id = $5
		RETURNING album_id, title, genre, tracks, group_id;
		`

	args := []interface{}{album.Title, album.Genre, album.Tracks, album.Group_id, album.Id}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return a.DB.QueryRowContext(ctx, query, args...).Scan(&album.Id, &album.Title, &album.Genre, &album.Tracks, &album.Group_id)
}


func (a AlbumModel) Delete(id int64) error {
	// Return an ErrRecordNotFound error if the movie ID is less than 1.
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM album
		WHERE album_id = $1
		`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := a.DB.ExecContext(ctx, query, id)
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
func (a AlbumModel) GetAll(title string, genre string, tracks int, filters Filters) ([]*Album, Metadata, error) {
	// Construct the SQL query to retrieve all movie records.
	query :=  fmt.Sprintf(`
		SELECT count(*) OVER(), album_id, title, genre, tracks, group_id
		FROM album
		WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '')
		AND (to_tsvector('simple', genre) @@ plainto_tsquery('simple', $2) OR $2 = '')
		AND (tracks = $3 OR $3 = 1)
		ORDER BY %s %s, group_id
		LIMIT $4 OFFSET $5`, filters.sortColumn(), filters.sortDirection())

	// Create a context with a 3-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// As our SQL query now has quite a few placeholder parameters, let's collect the
	// values for the placeholders in a slice. Notice here how we call the limit() and
	// offset() methods on the Filters struct to get the appropriate values for the
	// LIMIT and OFFSET clauses.
	args := []interface{}{title, genre, tracks, filters.limit(), filters.offset()}

	// Use QueryContext() to execute the query. This returns a sql.Rows resultset
	// containing the result.
	rows, err := a.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	// Importantly, defer a call to rows.Close() to ensure that the resultset is closed
	// before GetAll() returns.
	defer rows.Close()
	// Declare a totalRecords variable.
	totalRecords := 0
	// Initialize an empty slice to hold the movie data.
	albums := []*Album{}

	// Use rows.Next to iterate through the rows in the resultset.
	for rows.Next() {
		// Initialize an empty Movie struct to hold the data for an individual movie.
		var album Album
		// Scan the values from the row into the Movie struct. Again, note that we're
		// using the pq.Array() adapter on the genres field here.
		err := rows.Scan(
			&totalRecords, // Scan the count from the window function into totalRecords.
			&album.Id,
			&album.Title,
			&album.Genre,
			&album.Tracks,
			&album.Group_id,
		)
		if err != nil {
			return nil, Metadata{}, err
		}
		// Add the Movie struct to the slice.
		albums = append(albums, &album)
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
	return albums, metadata, nil
}
