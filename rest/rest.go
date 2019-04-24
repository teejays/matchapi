package rest

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// CleanAPIErrMessage is the message that should be used in at API response when we want to mask the
// internal error details from the client
var CleanAPIErrMessage = "There was an error processing the request. Please see the application logs"

// Authenticate should implement the authentication logic. It should should at the auth token
// and figure out what user context. Currently, this is not implemented and it only relies on
// and explicitly passed userID in the route.
func Authenticate(r *http.Request) (int, error) {
	return GetUserIDMuxVar(r)
}

// GetUserIDMuxVar extracts the userid param out of the request route
func GetUserIDMuxVar(r *http.Request) (int, error) {

	var vars = mux.Vars(r)

	idStr := vars["userid"]
	if idStr == "" {
		return -1, fmt.Errorf("could not find userID in the route")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return -1, fmt.Errorf("could not convert userID to a number: %v", err)
	}

	return id, nil
}
