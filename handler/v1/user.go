package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/teejays/clog"

	"github.com/teejays/matchapi/lib/auth"
	"github.com/teejays/matchapi/lib/rest"
	"github.com/teejays/matchapi/service/user/v1"
)

// HandleGetUser ...
// Example Request: curl -v localhost:8080/{userid}/v1/user
func HandleGetUser(w http.ResponseWriter, r *http.Request) {

	// Get the userID from the request
	userID, err := auth.Authenticate(r)
	if err != nil {
		clog.Error(err.Error())
		http.Error(w, "Could not authenticate the user", http.StatusUnauthorized)
		return
	}

	clog.Debugf("userID: %d", userID)

	// Get the incoming likes for the user
	user, err := user.GetUserByID(userID)
	if err != nil {
		clog.Error(err.Error())
		http.Error(w, rest.CleanAPIErrMessage, http.StatusInternalServerError)
		return
	}

	clog.Debugf("UserID fetched: %v", user)

	// Json marshal the response
	resp, err := json.Marshal(user)
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

// HandleCreateUser ...
// Example Request: curl -v -X "POST" localhost:8080//v1/user -d '{"FirstName":"Tom","LastName":"Harry", "Email": "tom.harry@email.com", "Gender": 3}'
func HandleCreateUser(w http.ResponseWriter, r *http.Request) {

	// Read the HTTP request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		clog.Error(err.Error())
		http.Error(w, "There was an error reading the request", http.StatusBadRequest)
		return
	}

	// Json unmarshal the request into the user.Profile
	var profile user.Profile
	err = json.Unmarshal(body, &profile)
	if err != nil {
		clog.Error(err.Error())
		http.Error(w, "There was an error json unmarshaling the request", http.StatusBadRequest)
		return
	}

	// Validate that the profile is has all required info
	err = profile.Validate()
	if err != nil {
		clog.Error(err.Error())
		http.Error(w, fmt.Sprintf("There was an error validating the request: %v", err), http.StatusBadRequest)
		return
	}

	// Update the profile of the given user
	user, err := user.NewUser(profile)
	if err != nil {
		clog.Error(err.Error())
		http.Error(w, rest.CleanAPIErrMessage, http.StatusInternalServerError)
		return
	}

	// Json marshal the updated profile so we can send it back
	resp, err := json.Marshal(user)
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

// HandleUpdateUserProfile ...
// Example Request: curl -v -X "PUT" localhost:8080/{userid}/v1/user -d '{"FirstName":"Jon","LastName":"Smith", "Email": "jon.smith@email.com", "Gender": 0}'
func HandleUpdateUserProfile(w http.ResponseWriter, r *http.Request) {

	// Get the userID from the request
	userID, err := auth.Authenticate(r)
	if err != nil {
		clog.Error(err.Error())
		http.Error(w, "Could not authenticate the user", http.StatusUnauthorized)
		return
	}

	clog.Debugf("userID: %d", userID)

	// Read the HTTP request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		clog.Error(err.Error())
		http.Error(w, "There was an error reading the request", http.StatusBadRequest)
		return
	}

	// Json unmarshal the request into the user.Profile
	var profile user.Profile
	err = json.Unmarshal(body, &profile)
	if err != nil {
		clog.Error(err.Error())
		http.Error(w, "There was an error json unmarshaling the request", http.StatusBadRequest)
		return
	}

	// Validate that the profile is has all required info
	err = profile.Validate()
	if err != nil {
		clog.Error(err.Error())
		http.Error(w, fmt.Sprintf("There was an error validating the request: %v", err), http.StatusBadRequest)
		return
	}

	// Get the user object
	user, err := user.GetUserByID(userID)
	if err != nil {
		clog.Error(err.Error())
		http.Error(w, rest.CleanAPIErrMessage, http.StatusInternalServerError)
		return
	}

	// Update the user object
	err = user.UpdateProfile(profile)
	if err != nil {
		clog.Error(err.Error())
		http.Error(w, rest.CleanAPIErrMessage, http.StatusInternalServerError)
		return
	}

	// Json marshal the updated profile so we can send it back
	resp, err := json.Marshal(user)
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
