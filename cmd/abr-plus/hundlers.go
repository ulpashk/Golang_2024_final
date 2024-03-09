package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ulpashk/Golang_2024/pkg/abr-plus/model"
)


func (app *application) respondWithError(w http.ResponseWriter, code int, message string) {
	app.respondWithJSON(w, code, map[string]string{"error": message})
}


func (app *application) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)

	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}



// Creating a new song
func (app *application) createSongHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Id       string `json:"id"`
		Title    string `json:"title"`
		Length   int    `json:"length"`
		Album_id int    `json:"albumId"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	song := &model.Song{
		Id:       input.Id,
		Title:    input.Title,
		Length:   input.Length,
		Album_id: input.Album_id,
	}

	err = app.models.Songs.Insert(song)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusCreated, song)
}



// Getting a specific song by id
func (app *application) getSongHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["songId"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid song ID")
		return
	}

	song, err := app.models.Songs.Get(id)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "404 Not Found")
		return
	}

	app.respondWithJSON(w, http.StatusOK, song)
}


// Updating a song by id
func (app *application) updateSongHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["songId"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid song ID")
		return
	}

	song, err := app.models.Songs.Get(id)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "404 Not Found")
		return
	}

	var input struct {
		Id       *string `json:"id"`
		Title    *string `json:"title"`
		Length   *int    `json:"length"`
		Album_id *int    `json:"albumId"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if input.Title != nil {
		song.Title = *input.Title
	}

	if input.Length != nil {
		song.Length = *input.Length
	}

	if input.Album_id != nil {
		song.Album_id = *input.Album_id
	}

	err = app.models.Songs.Update(song)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusOK, song)
}


// Deleting a song by id
func (app *application) deleteSongHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["songId"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid song ID")
		return
	}

	err = app.models.Songs.Delete(id)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}


func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(dst)
	if err != nil {
		return err
	}

	return nil
}