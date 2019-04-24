package likes

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/teejays/clog"

	"github.com/teejays/matchapi/users"
)

var cleanHTTPRespErrorMessage = "There was an error processing the request. Please see the application logs"

// ErrHTTPMethodNotSupported is returned when a handler receives a request with an unexpected HTTP method
var ErrHTTPMethodNotSupported = errors.New("HTTP method is not supported")

// HandleGetIncomingLikes ...
func HandleGetIncomingLikes(w http.ResponseWriter, r *http.Request) {
	clog.Infof("Request received for %s", "HandleGetIncomingLikes")

	var vars = mux.Vars(r)

	userIDString := vars["userid"]
	if userIDString == "" {
		http.Error(w, "could not find userID in the route", http.StatusBadRequest)
	}
	userID, err := strconv.Atoi(userIDString)
	if err != nil {
		clog.Error(err.Error())
		http.Error(w, "userID is not a number", http.StatusBadRequest)
	}

	clog.Debugf("userID: %d", userID)

	likes, err := GetIncomingLikesByUserID(users.UserID(userID))
	if err != nil {
		clog.Error(err.Error())
		http.Error(w, cleanHTTPRespErrorMessage, http.StatusInternalServerError)
	}

	clog.Debugf("userID Incoming Likes: %v", likes)

	// Json marshal the response
	resp, err := json.Marshal(likes)
	if err != nil {
		clog.Error(err.Error())
		http.Error(w, cleanHTTPRespErrorMessage, http.StatusInternalServerError)
	}

	w.Write(resp)
	clog.Info("Request succesfully processed")
}
