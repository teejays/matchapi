package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/teejays/clog"

	"github.com/teejays/matchapi/likes"
	"github.com/teejays/matchapi/users"
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

	v1 := r.PathPrefix("/v1").Subrouter()

	// add user id info to it
	v1 = v1.PathPrefix("/{userid}").Subrouter()

	v1.HandleFunc("/likes/incoming", likes.HandleGetIncomingLikes).
		Methods("GET")
	v1.HandleFunc("/user", users.HandleGetUser).
		Methods("GET")

	http.Handle("/", r)

	return http.ListenAndServe(fmt.Sprintf(":%d", listenPort), nil)
}
