package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/teejays/clog"

	"github.com/teejays/matchapi/db"
	"github.com/teejays/matchapi/handler/v1"
	handlerV2 "github.com/teejays/matchapi/handler/v2"
	"github.com/teejays/matchapi/lib/rest"
)

// listenPort is the port at which the server will listen. It can be
// passed as a command line flag using `--port <port>`. It defaults to 8080.
var listenPort = flag.Int("port", 8080, "port at which the server should listen")

// verbose mode can be turned on through the command line flag by passing `--verbose`.
// Turning it on increases the log level of the logging library we're using
var verbose = flag.Bool("verbose", false, "verbose mode")

func main() {
	var err error

	// Initialize the flags so we can start using them in our logic
	flag.Parse()

	// If verbose mode is off, we should lower the log level
	// The log level is by default set to the noisiest in teejays/clog
	if !*verbose {
		clog.LogLevel = 1
	}

	// Initialize the database: Consult the README for more details
	err = db.InitDB()
	if err != nil {
		clog.FatalErr(err)
	}

	// Initialize & start the webserver
	err = initServer(*listenPort)
	if err != nil {
		clog.FatalErr(err)
	}
}

// initServer setups and star the webserver
func initServer(port int) error {

	// Register the routes and handler so we can start directing the
	// requests to the right places
	registerHandlers()

	// Start the server
	clog.Infof("Listenining on: %d", port)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

// registerHandlers setups the routes and middleware for the webserver
func registerHandlers() {

	// Create a new gorilla.mux router
	r := mux.NewRouter()

	// Set up the routes

	// 1. Unauthenticated Routes: We are going to goahead and deal with pseudo-authenticated
	// routes but first, let's create routes that do no need any authentication
	// - Unauthenticated V1:
	rv1 := r.PathPrefix("/v1").Subrouter()
	rv1.HandleFunc("/user", handler.HandleCreateUser).Methods(http.MethodPost)

	// 2. Authenticated Routes: Create a route path that takes userid as the first param
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
