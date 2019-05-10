package handler

import (
	"encoding/json"
	"net/http"

	"github.com/teejays/clog"

	"github.com/teejays/matchapi/lib/auth"
	"github.com/teejays/matchapi/lib/rest"
	"github.com/teejays/matchapi/service/like/v1"
	likeV2 "github.com/teejays/matchapi/service/like/v2"
)

// HandleGetIncomingLikes ...
// Example Request: curl -v localhost:8080/{userid}/v1/like/incoming
func HandleGetIncomingLikes(w http.ResponseWriter, r *http.Request) {

	// Get the userID from the request
	userID, err := auth.GetUserIdFromRequest(r)
	if err != nil {
		clog.Error(err.Error())
		http.Error(w, "Could not authenticate the auth", http.StatusUnauthorized)
		return
	}

	// convert the id to the user.UserID type alias

	clog.Debugf("userID: %d", userID)

	// Get the incoming likes for the user
	likesV2, err := likeV2.GetIncomingLikesByUserID(userID)
	if err != nil {
		clog.Error(err.Error())
		http.Error(w, rest.CleanAPIErrMessage, http.StatusInternalServerError)
		return
	}

	// To make the API compatible
	var likes []like.Like
	for _, lv2 := range likesV2 {
		l := like.ConvertV2IncomimgLikeToV1(lv2)
		likes = append(likes, l)
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
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(resp)
	if err != nil {
		clog.Error(err.Error())
		http.Error(w, rest.CleanAPIErrMessage, http.StatusInternalServerError)
		return
	}

	clog.Info("Request succesfully processed")

}
