package like

import (
	"encoding/json"
	"net/http"

	"github.com/teejays/clog"

	"github.com/teejays/matchapi/v1/user"
)

var cleanHTTPRespErrorMessage = "There was an error processing the request. Please see the application logs"

// HandleGetIncomingLikes ...
func HandleGetIncomingLikes(w http.ResponseWriter, r *http.Request) {

	clog.Infof("Request received for %s", "HandleGetIncomingLikes")

	// Get the userID from the request
	userID, err := user.GetUserIDMuxVar(r)
	if err != nil {
		clog.Error(err.Error())
		http.Error(w, "userID is not a number", http.StatusBadRequest)
		return
	}

	clog.Debugf("userID: %d", userID)

	// Get the incoming likes for the user
	likes, err := GetIncomingLikesByUserID(userID)
	if err != nil {
		clog.Error(err.Error())
		http.Error(w, cleanHTTPRespErrorMessage, http.StatusInternalServerError)
		return
	}

	clog.Debugf("userID Incoming Likes: %v", likes)

	// Json marshal the response
	resp, err := json.Marshal(likes)
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
