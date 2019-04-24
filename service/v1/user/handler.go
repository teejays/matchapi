package user

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/teejays/clog"

	"github.com/teejays/matchapi/rest"
)

// HandleUpdateUserProfile ...
// Example Request: curl -v -X "POST" localhost:8080/{userid}/v1/user -d '{"FirstName":"Jon","LastName":"Smith", "Email": "jon.smith@email.com", "Gender": 0}'
func HandleUpdateUserProfile(w http.ResponseWriter, r *http.Request) {

	clog.Infof("Request received for %s", "HandleUpdateUserProfile")

	// Get the userID from the request
	id, err := rest.Authenticate(r)
	if err != nil {
		clog.Error(err.Error())
		http.Error(w, "Could not authenticate the user", http.StatusUnauthorized)
		return
	}

	// convert the id to the user.UserID type alias
	userID := UserID(id)

	clog.Debugf("userID: %d", userID)

	// Read the HTTP request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		clog.Error(err.Error())
		http.Error(w, "There was an error reading the request", http.StatusBadRequest)
		return
	}

	// Json unmarshal the request into the user.Profile
	var profile Profile
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
	user, err := UpdateUserByID(userID, profile)
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

// HandleGetUser ...
// Example Request: curl -v localhost:8080/{userid}/v1/user
func HandleGetUser(w http.ResponseWriter, r *http.Request) {

	clog.Infof("Request received for %s", "HandleGetUser")

	// Get the userID from the request
	id, err := rest.Authenticate(r)
	if err != nil {
		clog.Error(err.Error())
		http.Error(w, "Could not authenticate the user", http.StatusUnauthorized)
		return
	}

	// convert the id to the user.UserID type alias
	userID := UserID(id)

	clog.Debugf("userID: %d", userID)

	// Get the incoming likes for the user
	user, err := GetUserByID(userID)
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
