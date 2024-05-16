package main

import (
	"errors"
	"fmt"
	"net/http"

	"goproject/internal/data"
	"goproject/internal/validator"
)

func (app *application) createAlbumHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title    string `json:"title"`
		Genre    string `json:"genre"`
		Tracks   int    `json:"tracks"`
		Group_id int    `json:"groupId"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	album := &data.Album{
		Title:    input.Title,
		Genre:    input.Genre,
		Tracks:   input.Tracks,
		Group_id: input.Group_id,
	}

	v := validator.New()

	if data.ValidateAlbum(v, album); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Albums.Insert(album)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/albums/%d", album.Id))

	err = app.writeJSON(w, http.StatusCreated, envelope{"album": album}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) showAlbumHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	album, err := app.models.Albums.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"album": album}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateAlbumHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)

	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	album, err := app.models.Albums.Get(id)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		Title    string `json:"title"`
		Genre    string `json:"genre"`
		Tracks   int    `json:"tracks"`
		Group_id int    `json:"groupId"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	album.Title = input.Title
	album.Genre = input.Genre
	album.Tracks = input.Tracks
	album.Group_id = input.Group_id

	v := validator.New()

	if data.ValidateAlbum(v, album); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Albums.Update(album)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"album": album}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteAlbumHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)

	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Albums.Delete(id)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "album successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listAlbumsHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Title  string `json:"title"`
		Genre  string `json:"genre"`
		Tracks int    `json:"tracks"`
		data.Filters
	}

	v := validator.New()

	qs := r.URL.Query()

	input.Title = app.readString(qs, "title", "")
	input.Genre = app.readString(qs, "genre", "")

	input.Tracks = app.readInt(qs, "tracks", 1, v)

	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 30, v)

	input.Filters.Sort = app.readString(qs, "sort", "album_id")

	input.Filters.SortSafelist = []string{"album_id", "title", "genre", "tracks", "-album_id", "-title", "-genre", "-tracks"}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	albums, metadata, err := app.models.Albums.GetAll(input.Title, input.Genre, input.Tracks, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"albums": albums, "metadata": metadata}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) getAlbumsByGroupHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	var input struct {
		Title  string 	`json:"title"`
		Genre  string 	`json:"genre"`
		Tracks int    	`json:"tracks"`
		data.Filters
	}

	v := validator.New()

	qs := r.URL.Query()
	
	// input.Group_id = id
	input.Title = app.readString(qs, "title", "")
	input.Genre = app.readString(qs, "genre", "")

	input.Tracks = app.readInt(qs, "tracks", 1, v)

	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 30, v)

	input.Filters.Sort = app.readString(qs, "sort", "album_id")

	input.Filters.SortSafelist = []string{"album_id", "title", "genre", "tracks", "-album_id", "-title", "-genre", "-tracks"}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	albums, metadata, err := app.models.Albums.GetAlbumsByGroup(id, input.Title, input.Genre, input.Tracks, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"albums": albums, "metadata": metadata}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}