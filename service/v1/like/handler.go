package like

import (
	"encoding/json"
	"net/http"

	"github.com/teejays/clog"

	"github.com/teejays/matchapi/lib/rest"
	"github.com/teejays/matchapi/service/v1/user"
)

// HandleGetIncomingLikes ...
// Example Request: curl -v localhost:8080/{userid}/v1/like/incoming
func HandleGetIncomingLikes(w http.ResponseWriter, r *http.Request) {

	// Get the userID from the request
	userID, err := user.Authenticate(r)
	if err != nil {
		clog.Error(err.Error())
		http.Error(w, "Could not authenticate the user", http.StatusUnauthorized)
		return
	}

	// convert the id to the user.UserID type alias

	clog.Debugf("userID: %d", userID)

	// Get the incoming likes for the user
	likes, err := GetIncomingLikesByUserID(userID)
	if err != nil {
		clog.Error(err.Error())
		http.Error(w, rest.CleanAPIErrMessage, http.StatusInternalServerError)
		return
	}

	clog.Debugf("userID Incoming Likes: %v", likes)

	// Json marshal the response
	resp, err := json.Marshal(likes)
	if err != nil {
		clog.Error(err.Error())
		http.Error(w, rest.CleanAPIErrMessage, http.StatusInternalServerError)
		return
	}

	// Write the HTTP response
	_, err = w.Write(resp)
	if err != nil {
		clog.Error(err.Error())
		http.Error(w, rest.CleanAPIErrMessage, http.StatusInternalServerError)
		return
	}

	clog.Info("Request succesfully processed")

}
