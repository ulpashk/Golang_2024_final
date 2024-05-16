package main

import (
	"net/http"
	"github.com/julienschmidt/httprouter"

)

func (app *application) routes() http.Handler {

	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)

	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/api/v1/healthcheck", app.healthcheckHandler)

	router.HandlerFunc(http.MethodGet, "/api/v1/groups", app.requirePermission("read", app.listGroupsHandler))
	router.HandlerFunc(http.MethodPost, "/api/v1/groups", app.requirePermission("read", app.createGroupHandler))
	router.HandlerFunc(http.MethodGet, "/api/v1/groups/:id", app.requirePermission("read", app.showGroupHandler))
	router.HandlerFunc(http.MethodPut, "/api/v1/groups/:id", app.requirePermission("read", app.updateGroupHandler))
	router.HandlerFunc(http.MethodDelete, "/api/v1/groups/:id", app.requirePermission("write", app.deleteGroupHandler))

	router.HandlerFunc(http.MethodGet, "/api/v1/albums", app.requirePermission("read", app.listAlbumsHandler))
	router.HandlerFunc(http.MethodPost, "/api/v1/albums", app.requirePermission("read", app.createAlbumHandler))
	router.HandlerFunc(http.MethodGet, "/api/v1/albums/:id", app.requirePermission("read", app.showAlbumHandler))
	router.HandlerFunc(http.MethodPut, "/api/v1/albums/:id", app.requirePermission("read", app.updateAlbumHandler))
	router.HandlerFunc(http.MethodDelete, "/api/v1/albums/:id", app.requirePermission("write", app.deleteAlbumHandler))
	router.HandlerFunc(http.MethodGet, "/api/v1/groups/:id/albums", app.requirePermission("read", app.getAlbumsByGroupHandler))

	router.HandlerFunc(http.MethodGet, "/api/v1/songs", app.requirePermission("read", app.listSongsHandler))
	router.HandlerFunc(http.MethodPost, "/api/v1/songs", app.requirePermission("read", app.createSongHandler))
	router.HandlerFunc(http.MethodGet, "/api/v1/songs/:id", app.requirePermission("read", app.showSongHandler))
	router.HandlerFunc(http.MethodPut, "/api/v1/songs/:id", app.requirePermission("read", app.updateSongHandler))
	router.HandlerFunc(http.MethodDelete, "/api/v1/songs/:id", app.requirePermission("write", app.deleteSongHandler))
	router.HandlerFunc(http.MethodGet, "/api/v1/albums/:id/songs", app.requirePermission("read", app.getSongsByAlbum))

	router.HandlerFunc(http.MethodPost, "/api/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/api/v1/users/activated", app.activateUserHandler)
	router.HandlerFunc(http.MethodPost, "/api/v1/tokens/login", app.createAuthenticationTokenHandler)

	return app.authenticate(router)

}