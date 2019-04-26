package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/teejays/clog"

	"github.com/teejays/matchapi/db"
	"github.com/teejays/matchapi/handler/v1"
	handlerV2 "github.com/teejays/matchapi/handler/v2"
	"github.com/teejays/matchapi/lib/rest"
)

const listenPort = 8080 // we should probably move this to a config file, env variable or command-line flag

func main() {
	var err error

	// Lower the log level
	clog.LogLevel = 1

	// Initialize the database
	err = db.InitDB()
	if err != nil {
		clog.FatalErr(err)
	}

	// Initialize the webserver
	err = initServer()
	if err != nil {
		clog.FatalErr(err)
	}
}

func initServer() error {

	// Register the routes and handler so we can start directing the
	// requests to the right places
	registerHandlers()

	// Start the server
	clog.Infof("Listenining on: %d", listenPort)
	return http.ListenAndServe(fmt.Sprintf(":%d", listenPort), nil)
}

func registerHandlers() {

	// Create a new gorilla.mux router
	r := mux.NewRouter()

	// Unauthenticated Routes: We are going to goahead and deal with pseudo-authenticated
	// routes but first, let's create routes that do no need any authentication
	// - Unauthenticated V1:
	rv1 := r.PathPrefix("/v1").Subrouter()
	rv1.HandleFunc("/user", handler.HandleCreateUser).Methods(http.MethodPost)

	// Authenticated Routes: Create a route path that takes userid as the first param
	// We are going to use it as a prxoy for authentication
	a := r.PathPrefix("/{userid}").Subrouter()

	// - Authenticated V1; Create a path that takes v1 as prefix
	av1 := a.PathPrefix("/v1").Subrouter()
	av1.HandleFunc("/user", handler.HandleGetUser).Methods(http.MethodGet)
	av1.HandleFunc("/user", handler.HandleUpdateUserProfile).Methods(http.MethodPut)
	av1.HandleFunc("/like/incoming", handler.HandleGetIncomingLikes).Methods("GET")

	// - Authenticated V2; Create a path that takes v2 as prefix
	av2 := a.PathPrefix("/v2").Subrouter()
	av2.HandleFunc("/like/incoming", handlerV2.HandleGetIncomingLikes).Methods("GET")
	av2.HandleFunc("/like", handlerV2.HandlePostLike).Methods("POST")

	// Add a simple middleware function so we can log the requests
	r.Use(rest.LoggerMiddleware)

	// Register the router as the handler in the standard net/http package
	http.Handle("/", r)

}
