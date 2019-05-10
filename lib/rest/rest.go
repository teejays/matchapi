package rest

import (
	"fmt"
	"net/http"
	"strconv"
	"context"

	"github.com/gorilla/mux"
	"github.com/teejays/clog"

	"github.com/teejays/matchapi/lib/auth"
	"github.com/teejays/matchapi/lib/pk"
)

// CleanAPIErrMessage is the message that should be used in at API response when we want to mask the
// internal error details from the client
var CleanAPIErrMessage = "There was an error processing the request. Please see the application logs"

// GetUserIDMuxVar extracts the userid param out of the request route
func GetUserIDMuxVar(r *http.Request) (pk.ID, error) {

	var vars = mux.Vars(r)

	idStr := vars["userid"]
	if idStr == "" {
		return -1, fmt.Errorf("could not find userID in the route")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return -1, fmt.Errorf("could not convert userID to a number: %v", err)
	}

	return pk.ID(id), nil
}

// AddContextMiddleware is a http.Handler middleware function that logs any request received
func AddContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Call the next handler, with request.WithContext()
		next.ServeHTTP(w, r.WithContext(context.Background()))
	})
}

// AuthenticateMiddleware authenticates the requests
func AuthenticateMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the Authenticated Request
		ar, err := auth.AuthenticateRequest(r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, ar)
	})
}

// LoggerMiddleware is a http.Handler middleware function that logs any request received
func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log the request
		clog.Debugf("Server: HTTP request received for %s %s", r.Method, r.URL.Path)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}
