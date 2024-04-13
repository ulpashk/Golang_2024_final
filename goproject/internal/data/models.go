package data

import (
	"database/sql"
	"errors"
)

// Define a custom ErrRecordNotFound error. We'll return this from our Get() method when
// looking up a movie that doesn't exist in our database.
var (
	ErrRecordNotFound = errors.New("record not found")
)

// Create a Models struct which wraps the MovieModel. We'll add other models to this,
// like a UserModel and PermissionModel, as our build progresses.
type Models struct {
	// Set the Movies field to be an interface containing the methods that both the
	// 'real' model and mock model need to support.
	Songs interface{
		Insert(song *Song) error
		Get(id int64) (*Song, error)
		GetAll(title string, length int, filters Filters) ([]*Song, Metadata, error)
		Update(song *Song) error
		Delete(id int64) error
	}
}

// Create a helper function which returns a Models instance containing the mock models
// only.
func NewModels(db *sql.DB) Models {
	return Models{
		Songs: SongModel{DB: db},
	}
}
