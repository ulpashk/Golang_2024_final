package data

import (
	"database/sql"
	"errors"
)

// Define a custom ErrRecordNotFound error. We'll return this from our Get() method when
// looking up a movie that doesn't exist in our database.
var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict = errors.New("edit conflict")
)

// Create a Models struct which wraps the MovieModel. We'll add other models to this,
// like a UserModel and PermissionModel, as our build progresses.
type Models struct {
	// Set the Movies field to be an interface containing the methods that both the
	// 'real' model and mock model need to support.
	Groups interface{
		Insert(group *Group) error
		Get(id int64) (*Group, error)
		GetAll(name string, num_of_members int, filters Filters) ([]*Group, Metadata, error)
		Update(group *Group) error
		Delete(id int64) error
	}
	Albums interface{
		Insert(album *Album) error
		Get(id int64) (*Album, error)
		GetAll(title string, genre string, tracks int, filters Filters) ([]*Album, Metadata, error)
		Update(album *Album) error
		Delete(id int64) error
	}
	Songs interface{
		Insert(song *Song) error
		Get(id int64) (*Song, error)
		GetAll(title string, length int, filters Filters) ([]*Song, Metadata, error)
		Update(song *Song) error
		Delete(id int64) error
	}
	
	Users UserModel
	Tokens TokenModel
	Permissions PermissionModel
}

// Create a helper function which returns a Models instance containing the mock models
// only.
func NewModels(db *sql.DB) Models {
	return Models{
		Groups: GroupModel{DB: db},
		Albums: AlbumModel{DB: db},
		Songs: SongModel{DB: db},
		Permissions: PermissionModel{DB: db},
		Tokens: TokenModel{DB: db}, 
		Users: UserModel{DB: db},
	}
}
