package main

import (
	"fmt"
	"net/http"
	"errors"

	"goproject/internal/data"
	"goproject/internal/validator"
)


func (app *application) createSongHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		// Id       int 	`json:"id"`
		Title    string `json:"title"`
		Length   int    `json:"length"`
		Album_id int    `json:"albumId"`
	}

	// Initialize a new json.Decoder instance which reads from the request body, and
	// then use the Decode() method to decode the body contents into the input struct.
	// Importantly, notice that when we call Decode() we pass a *pointer* to the input
	// struct as the target decode destination. If there was an error during decoding,
	// we also use our generic errorResponse() helper to send the client a 400 Bad
	// Request response containing the error message.
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Copy the values from the input struct to a new Song struct.

	song := &data.Song{
		// Id:       input.Id,
		Title:    input.Title,
		Length:   input.Length,
		Album_id: input.Album_id,
	}

	// Initialize a new Validator instance.
	v := validator.New()

	// Call the ValidateSong() function and return a response containing the errors if
	// any of the checks fail.
	if data.ValidateSong(v, song); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}


	// Call the Insert() method on our songs model, passing in a pointer to the
	// validated song struct. This will create a record in the database and update the
	// song struct with the system-generated information.
	err = app.models.Songs.Insert(song)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// When sending a HTTP response, we want to include a Location header to let the
	// client know which URL they can find the newly-created resource at. We make an
	// empty http.Header map and then use the Set() method to add a new Location header,
	// interpolating the system-generated ID for our new movie in the URL.
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/songs/%d", song.Id))
	// Write a JSON response with a 201 Created status code, the song data in the
	// response body, and the Location header.

	err = app.writeJSON(w, http.StatusCreated, envelope{"song": song}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}


func (app *application) showSongHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}	

	song, err := app.models.Songs.Get(id)
	if err != nil {
		switch {
			case errors.Is(err, data.ErrRecordNotFound):
				app.notFoundResponse(w, r)
			default:
				app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"song": song}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateSongHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)

	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	song, err := app.models.Songs.Get(id)
	
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
		Length   int    `json:"length"`
		Album_id int    `json:"albumId"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	song.Title = input.Title
	song.Length = input.Length
	song.Album_id = input.Album_id


	v := validator.New()

	if data.ValidateSong(v, song); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Songs.Update(song)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"song": song}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteSongHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)

	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Songs.Delete(id)

	if err != nil {
		switch {
			case errors.Is(err, data.ErrRecordNotFound):
				app.notFoundResponse(w, r)
			default:
				app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "song successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listSongsHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Title    string `json:"title"`
		Length   int    `json:"length"`
		data.Filters
	}

	v := validator.New()

	qs := r.URL.Query()

	input.Title = app.readString(qs, "title", "")


	input.Length = app.readInt(qs, "length", 1, v)
	 
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 30, v)

	input.Filters.Sort = app.readString(qs, "sort", "song_id")

	input.Filters.SortSafelist = []string{"song_id", "title", "length", "-song_id", "-title", "-length"}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	songs, metadata, err := app.models.Songs.GetAll(input.Title, input.Length, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	
	err = app.writeJSON(w, http.StatusOK, envelope{"songs": songs, "metadata": metadata}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
	