package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/teejays/clog"

	"github.com/teejays/matchapi/like"
	"github.com/teejays/matchapi/user"
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

	v1.HandleFunc("/likes/incoming", like.HandleGetIncomingLikes).
		Methods("GET")
	v1.HandleFunc("/user", user.HandleGetUser).
		Methods("GET")
	v1.HandleFunc("/user", user.HandleUpdateUserProfile).
		Methods("POST")

	http.Handle("/", r)

	return http.ListenAndServe(fmt.Sprintf(":%d", listenPort), nil)
}
