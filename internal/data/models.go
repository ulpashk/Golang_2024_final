package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict = errors.New("edit conflict")
)

type Models struct {
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
		GetAlbumsByGroup(groupId int64, title string, genre string, tracks int, filters Filters) ([]*Album, Metadata, error)
	}
	Songs interface{
		Insert(song *Song) error
		Get(id int64) (*Song, error)
		GetAll(title string, length int, filters Filters) ([]*Song, Metadata, error)
		Update(song *Song) error
		Delete(id int64) error
		GetSongsByAlbum(id int64, title string, length int, filters Filters) ([]*Song, Metadata, error)
	}
	
	Users UserModel
	Tokens TokenModel
	Permissions PermissionModel
}

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
