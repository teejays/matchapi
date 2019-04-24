package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/teejays/clog"

	"github.com/teejays/matchapi/like/v1"
	"github.com/teejays/matchapi/user/v1"
)

const listenPort = 8080 // we should probably move this to a config file, env variable or command-line flag

func main() {
	err := startServer()
	if err != nil {
		clog.FatalErr(err)
	}
}

func startServer() error {
	r := mux.NewRouter()

	// add user id as a route variable to the path, as a proxy for auth (for now)
	r = r.PathPrefix("/{userid}").Subrouter()

	// Setup V1 routes
	v1 := r.PathPrefix("/v1").Subrouter()

	v1.HandleFunc("/likes/incoming", like.HandleGetIncomingLikes).
		Methods("GET")
	v1.HandleFunc("/user", user.HandleGetUser).
		Methods("GET")
	v1.HandleFunc("/user", user.HandleUpdateUserProfile).
		Methods("POST")

	v1.HandleFunc("/likes/incoming", like.HandleGetIncomingLikes).
		Methods("GET")

	// Register the handler
	http.Handle("/", r)

	return http.ListenAndServe(fmt.Sprintf(":%d", listenPort), nil)
}
