package main

import (
	"fmt"
	"net/http"
	"errors"

	"goproject/internal/data"
	"goproject/internal/validator"
)


func (app *application) createGroupHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name    		string `json:"name"`
		Num_of_members  int    `json:"num_of_members"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	group := &data.Group{
		Name:    input.Name,
		Num_of_members:   input.Num_of_members,
	}

	v := validator.New()

	if data.ValidateGroup(v, group); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Groups.Insert(group)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/groups/%d", group.Id))

	err = app.writeJSON(w, http.StatusCreated, envelope{"group": group}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}


func (app *application) showGroupHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}	

	group, err := app.models.Groups.Get(id)
	if err != nil {
		switch {
			case errors.Is(err, data.ErrRecordNotFound):
				app.notFoundResponse(w, r)
			default:
				app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"group": group}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateGroupHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)

	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	group, err := app.models.Groups.Get(id)
	
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
		Name    		string `json:"name"`
		Num_of_members  int    `json:"num_of_members"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	group.Name = input.Name
	group.Num_of_members = input.Num_of_members

	v := validator.New()

	if data.ValidateGroup(v, group); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Groups.Update(group)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"group": group}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteGroupHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)

	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Groups.Delete(id)

	if err != nil {
		switch {
			case errors.Is(err, data.ErrRecordNotFound):
				app.notFoundResponse(w, r)
			default:
				app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "group successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listGroupsHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Name    		string `json:"name"`
		Num_of_members  int    `json:"num_of_members"`
		data.Filters
	}

	v := validator.New()

	qs := r.URL.Query()

	input.Name = app.readString(qs, "name", "")


	input.Num_of_members = app.readInt(qs, "num_of_members", 1, v)
	 
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 30, v)

	input.Filters.Sort = app.readString(qs, "sort", "group_id")

	input.Filters.SortSafelist = []string{"group_id", "name", "num_of_members", "-group_id", "-name", "-num_of_members"}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	groups, metadata, err := app.models.Groups.GetAll(input.Name, input.Num_of_members, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	
	err = app.writeJSON(w, http.StatusOK, envelope{"metadata": metadata, "groups": groups}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
	