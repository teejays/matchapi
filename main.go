package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/teejays/clog"

	"github.com/teejays/matchapi/lib/rest"
	"github.com/teejays/matchapi/service/v1/like"
	"github.com/teejays/matchapi/service/v1/user"
	likeV2 "github.com/teejays/matchapi/service/v2/like"
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
	v1.HandleFunc("/user", user.HandleGetUser).
		Methods("GET")
	v1.HandleFunc("/user", user.HandleUpdateUserProfile).
		Methods("POST")
	v1.HandleFunc("/like/incoming", like.HandleGetIncomingLikes).
		Methods("GET")

	// Setup V1 routes
	v2 := r.PathPrefix("/v2").Subrouter()
	v2.HandleFunc("/like/incoming", likeV2.HandleGetIncomingLikes).
		Methods("GET")

	// Register the handler
	r.Use(rest.LoggerMiddleware)
	http.Handle("/", r)

	return http.ListenAndServe(fmt.Sprintf(":%d", listenPort), nil)
}
