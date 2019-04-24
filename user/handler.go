package user

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/teejays/clog"
)

var cleanHTTPRespErrorMessage = "There was an error processing the request. Please see the application logs"

// GetUserIDMuxVar extracts the userid param out of the request route
func GetUserIDMuxVar(r *http.Request) (UserID, error) {
	var id UserID
	var err error

	var vars = mux.Vars(r)

	idStr := vars["userid"]
	if idStr == "" {
		return id, fmt.Errorf("could not find userID in the route")
	}

	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		return id, fmt.Errorf("could not convert userID to a number: %v", err)
	}

	return UserID(idInt), nil
}

// HandleGetUser ...
func HandleGetUser(w http.ResponseWriter, r *http.Request) {

	clog.Infof("Request received for %s", "HandleGetUser")

	// Get the userID from the request
	userID, err := GetUserIDMuxVar(r)
	if err != nil {
		clog.Error(err.Error())
		http.Error(w, "userID is not a number", http.StatusBadRequest)
		return
	}

	clog.Debugf("userID: %d", userID)

	// Get the incoming likes for the user
	user, err := GetUserByID(userID)
	if err != nil {
		clog.Error(err.Error())
		http.Error(w, cleanHTTPRespErrorMessage, http.StatusInternalServerError)
		return
	}

	clog.Debugf("UserID fetched: %v", user)

	// Json marshal the response
	resp, err := json.Marshal(user)
	if err != nil {
		clog.Error(err.Error())
		http.Error(w, cleanHTTPRespErrorMessage, http.StatusInternalServerError)
		return
	}

	// Write the HTTP response
	_, err = w.Write(resp)
	if err != nil {
		clog.Error(err.Error())
		http.Error(w, cleanHTTPRespErrorMessage, http.StatusInternalServerError)
		return
	}

	clog.Info("Request succesfully processed")

}