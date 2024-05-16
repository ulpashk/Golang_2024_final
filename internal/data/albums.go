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

	query := `
		INSERT INTO album(title, genre, tracks, group_id)
		VALUES ($1, $2, $3, $4)
		RETURNING album_id, title, genre, tracks, group_id;`

	args := []interface{}{album.Title, album.Genre, album.Tracks, album.Group_id}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return a.DB.QueryRowContext(ctx, query, args...).Scan(&album.Id, &album.Title, &album.Genre, &album.Tracks, &album.Group_id)
}


func (a AlbumModel) Get(id int64) (*Album, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

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

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}


func (a AlbumModel) GetAll(title string, genre string, tracks int, filters Filters) ([]*Album, Metadata, error) {
	query :=  fmt.Sprintf(`
		SELECT count(*) OVER(), album_id, title, genre, tracks, group_id
		FROM album
		WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '')
		AND (to_tsvector('simple', genre) @@ plainto_tsquery('simple', $2) OR $2 = '')
		AND (tracks = $3 OR $3 = 1)
		ORDER BY %s %s, group_id
		LIMIT $4 OFFSET $5`, filters.sortColumn(), filters.sortDirection())


	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{title, genre, tracks, filters.limit(), filters.offset()}

	rows, err := a.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}

	defer rows.Close()

	totalRecords := 0

	albums := []*Album{}

	for rows.Next() {

		var album Album

		err := rows.Scan(
			&totalRecords, 
			&album.Id,
			&album.Title,
			&album.Genre,
			&album.Tracks,
			&album.Group_id,
		)
		if err != nil {
			return nil, Metadata{}, err
		}

		albums = append(albums, &album)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)
	return albums, metadata, nil
}


func (a AlbumModel) GetAlbumsByGroup(id int64, title string, genre string, tracks int, filters Filters) ([]*Album, Metadata, error) {
	query :=  fmt.Sprintf(`
		SELECT count(*) OVER(), album_id, title, genre, tracks, group_id
		FROM album
		WHERE (group_id = $1)
		AND (to_tsvector('simple', title) @@ plainto_tsquery('simple', $2) OR $2 = '')
		AND (to_tsvector('simple', genre) @@ plainto_tsquery('simple', $3) OR $3 = '')
		AND (tracks = $4 OR $4 = 1)
		ORDER BY %s %s, album_id
		LIMIT $5 OFFSET $6`, filters.sortColumn(), filters.sortDirection())


	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{id, title, genre, tracks, filters.limit(), filters.offset()}

	rows, err := a.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}

	defer rows.Close()

	totalRecords := 0

	albums := []*Album{}

	for rows.Next() {

		var album Album

		err := rows.Scan(
			&totalRecords, 
			&album.Id,
			&album.Title,
			&album.Genre,
			&album.Tracks,
			&album.Group_id,
		)
		if err != nil {
			return nil, Metadata{}, err
		}

		albums = append(albums, &album)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)
	return albums, metadata, nil
}