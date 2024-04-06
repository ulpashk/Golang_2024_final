package model

import (
	"context"
	"database/sql"
	"log"
	"time"
)

type Song struct{
	Id 			string	`json:"id"`
	Title 		string	`json:"title"`
	Length 		int		`json:"length"`
	Album_id	int 	`json:"albumId"`
} 


type SongModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}


// Inserting a new song
func (s SongModel) Insert(song *Song) error {
	query := `
		INSERT INTO song(song_id, title, length, album_id)
		VALUES ($1, $2, $3, $4)
		RETURNING song_id, title, length, album_id;
		`
	args := []interface{}{song.Id, song.Title, song.Length, song.Album_id}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return s.DB.QueryRowContext(ctx, query, args...).Scan(&song.Id, &song.Title, &song.Length, &song.Album_id)
}


// Getting a song by id
func (s SongModel) Get(id int) (*Song, error) {
	// Retrieve a specific menu item based on its ID.
	query := `
		SELECT song_id, title, length, album_id
		FROM song
		WHERE song_id = $1;
		`
	var song Song
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := s.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&song.Id, &song.Title, &song.Length, &song.Album_id)
	if err != nil {
		return nil, err
	}
	return &song, nil
}


// Updating a song by id
func (s SongModel) Update(song *Song) error {
	query := `
		UPDATE song
		SET title = $1, length = $2, album_id = $3
		WHERE song_id = $4
		RETURNING song_id, title, length, album_id;
		`
	args := []interface{}{song.Title, song.Length, song.Album_id, song.Id}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return s.DB.QueryRowContext(ctx, query, args...).Scan(&song.Id, &song.Title, &song.Length, &song.Album_id)
}


//Deleting a song by id
func (s SongModel) Delete(id int) error {
	query := `
		DELETE FROM song
		WHERE song_id = $1
		`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := s.DB.ExecContext(ctx, query, id)
	return err
}