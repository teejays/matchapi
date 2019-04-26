package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/teejays/clog"

	"github.com/teejays/matchapi/lib/auth"
	"github.com/teejays/matchapi/lib/rest"
	"github.com/teejays/matchapi/service/like/v2"
)

// HandleGetIncomingLikes ...
// Example Request: curl -v localhost:8080/{userid}/v2/like/incoming
func HandleGetIncomingLikes(w http.ResponseWriter, r *http.Request) {

	// Get the userID from the request
	userID, err := auth.Authenticate(r)
	if err != nil {
		clog.Error(err.Error())
		http.Error(w, "Could not authenticate the user", http.StatusUnauthorized)
		return
	}

	// convert the id to the user.UserID type alias

	clog.Debugf("userID: %d", userID)

	// Get the incoming likes for the user
	likes, err := like.GetIncomingLikesByUserID(userID)
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

// HandlePostLike ...
// Example Request: curl -v -X "POST" localhost:8080/{userid}/v1/like -d '{"ReceiverID": 3}'
func HandlePostLike(w http.ResponseWriter, r *http.Request) {

	// Get the userID from the request
	userID, err := auth.Authenticate(r)
	if err != nil {
		clog.Error(err.Error())
		http.Error(w, "Could not authenticate the user", http.StatusUnauthorized)
		return
	}

	// Read the HTTP request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		clog.Error(err.Error())
		http.Error(w, "There was an error reading the request", http.StatusBadRequest)
		return
	}

	// Json unmarshal the request into the BasicLike struct
	var blike like.BasicLike
	err = json.Unmarshal(body, &blike)
	if err != nil {
		clog.Error(err.Error())
		http.Error(w, "There was an error json unmarshaling the request", http.StatusBadRequest)
		return
	}

	// Validate that the profile is has all required info
	if err := blike.Validate(); err != nil {
		clog.Error(err.Error())
		http.Error(w, fmt.Sprintf("There was an error validating the request: %v", err), http.StatusBadRequest)
		return
	}

	// Update the profile of the given user
	like, err := like.NewLike(userID, blike)
	if err != nil {
		clog.Error(err.Error())
		http.Error(w, rest.CleanAPIErrMessage, http.StatusInternalServerError)
		return
	}

	// Json marshal the updated profile so we can send it back
	resp, err := json.Marshal(like)
	if err != nil {
		clog.Error(err.Error())
		http.Error(w, rest.CleanAPIErrMessage, http.StatusInternalServerError)
		return
	}

	// Write the update profile to the http response
	_, err = w.Write(resp)
	if err != nil {
		clog.Error(err.Error())
		http.Error(w, rest.CleanAPIErrMessage, http.StatusInternalServerError)
	}

	clog.Info("Request succesfully processed")

}
