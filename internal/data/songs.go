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
	query :=  fmt.Sprintf(`
		SELECT count(*) OVER(), song_id, title, length, album_id
		FROM songs
		WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '')
		AND (length = $2 OR $2 = 1)
		ORDER BY %s %s, song_id
		LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{title, length, filters.limit(), filters.offset()}

	rows, err := s.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}

	defer rows.Close()

	totalRecords := 0

	songs := []*Song{}


	for rows.Next() {

		var song Song

		err := rows.Scan(
			&totalRecords, 
			&song.Id,
			&song.Title,
			&song.Length,
			&song.Album_id,
		)
		if err != nil {
			return nil, Metadata{}, err
		}

		songs = append(songs, &song)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return songs, metadata, nil
}

func (s SongModel) GetSongsByAlbum(id int64, title string, length int, filters Filters) ([]*Song, Metadata, error) {
	query :=  fmt.Sprintf(`
		SELECT count(*) OVER(), song_id, title, length, album_id
		FROM songs
		WHERE (album_id = $1)
		AND (to_tsvector('simple', title) @@ plainto_tsquery('simple', $2) OR $2 = '')
		AND (length = $3 OR $3 = 1)
		ORDER BY %s %s, song_id
		LIMIT $4 OFFSET $5`, filters.sortColumn(), filters.sortDirection())


		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
	
		args := []interface{}{id, title, length, filters.limit(), filters.offset()}
	
		rows, err := s.DB.QueryContext(ctx, query, args...)
		if err != nil {
			return nil, Metadata{}, err
		}
	
		defer rows.Close()
	
		totalRecords := 0
	
		songs := []*Song{}
	
	
		for rows.Next() {
	
			var song Song
	
			err := rows.Scan(
				&totalRecords, 
				&song.Id,
				&song.Title,
				&song.Length,
				&song.Album_id,
			)
			if err != nil {
				return nil, Metadata{}, err
			}
	
			songs = append(songs, &song)
		}
	
		if err = rows.Err(); err != nil {
			return nil, Metadata{}, err
		}
	
		metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)
	
		return songs, metadata, nil
}