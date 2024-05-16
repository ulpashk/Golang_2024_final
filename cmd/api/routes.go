package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	// "expvar"
)

func (app *application) routes() http.Handler {
	// Initialize a new httprouter router instance.
	router := httprouter.New()

	// Convert the notFoundResponse() helper to a http.Handler using the
	// http.HandlerFunc() adapter, and then set it as the custom error handler for 404
	// Not Found responses.
	router.NotFound = http.HandlerFunc(app.notFoundResponse)

	// Likewise, convert the methodNotAllowedResponse() helper to a http.Handler and set
	// it as the custom error handler for 405 Method Not Allowed responses.
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	// Register the relevant methods, URL patterns and handler functions for our
	// endpoints using the HandlerFunc() method. Note that http.MethodGet and
	// http.MethodPost are constants which equate to the strings "GET" and "POST"
	// respectively.
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	router.HandlerFunc(http.MethodGet, "/v1/groups", app.requirePermission("read", app.listGroupsHandler))
	router.HandlerFunc(http.MethodPost, "/v1/groups", app.requirePermission("read", app.createGroupHandler))
	router.HandlerFunc(http.MethodGet, "/v1/groups/:id", app.requirePermission("read", app.showGroupHandler))
	router.HandlerFunc(http.MethodPut, "/v1/groups/:id", app.requirePermission("read", app.updateGroupHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/groups/:id", app.requirePermission("write", app.deleteGroupHandler))

	router.HandlerFunc(http.MethodGet, "/v1/albums", app.requirePermission("read", app.listAlbumsHandler))
	router.HandlerFunc(http.MethodPost, "/v1/albums", app.requirePermission("read", app.createAlbumHandler))
	router.HandlerFunc(http.MethodGet, "/v1/albums/:id", app.requirePermission("read", app.showAlbumHandler))
	router.HandlerFunc(http.MethodPut, "/v1/albums/:id", app.requirePermission("read", app.updateAlbumHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/albums/:id", app.requirePermission("write", app.deleteAlbumHandler))

	router.HandlerFunc(http.MethodGet, "/v1/songs", app.requirePermission("read", app.listSongsHandler))
	router.HandlerFunc(http.MethodPost, "/v1/songs", app.requirePermission("read", app.createSongHandler))
	router.HandlerFunc(http.MethodGet, "/v1/songs/:id", app.requirePermission("read", app.showSongHandler))
	router.HandlerFunc(http.MethodPut, "/v1/songs/:id", app.requirePermission("read", app.updateSongHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/songs/:id", app.requirePermission("write", app.deleteSongHandler))

	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)
	router.HandlerFunc(http.MethodPost, "/v1/tokens/login", app.createAuthenticationTokenHandler)

	// router.Handler(http.MethodGet, "/debug/vars", expvar.Handler())

	// Return the httprouter instance.
	// return router
	return app.authenticate(router)
	//return app.recoverPanic(app.enableCORS(app.rateLimit(app.authenticate(router))))

}
