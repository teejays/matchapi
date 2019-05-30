package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/teejays/clog"

	authLib "github.com/teejays/matchapi/lib/auth"
	"github.com/teejays/matchapi/lib/rest"
	"github.com/teejays/matchapi/service/auth/v1"
	"github.com/teejays/matchapi/service/user/v1"
)

// HandleGetUser ...
// Example Request: curl -v localhost:8080/{userid}/v1/user
func HandleGetUser(w http.ResponseWriter, r *http.Request) {

	// Get the userID from the request
	userID, err := authLib.GetUserIdFromRequest(r)
	if err != nil {
		clog.Error(err.Error())
		http.Error(w, "Could not authenticate the user", http.StatusUnauthorized)
		return
	}

	clog.Debugf("userID: %d", userID)

	// Get the incoming likes for the user
	usr, err := user.GetUserByID(userID)
	if err != nil {
		clog.Error(err.Error())
		http.Error(w, rest.CleanAPIErrMessage, http.StatusInternalServerError)
		return
	}

	clog.Debugf("UserID fetched: %v", usr)

	// Json marshal the response
	resp, err := json.Marshal(usr.Profile)
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

// CreateUserRequest represents that request object that is used when
// a new user is created
type CreateUserRequest struct {
	user.Profile
	Password string
}

// HandleCreateUser ...
// Example Request: curl -X "POST" localhost:8080/v1/user -d '{"FirstName":"Tom","LastName":"Harry", "Email": "tom.harry@email.com", "Gender": 3}'
func HandleCreateUser(w http.ResponseWriter, r *http.Request) {

	// Read the HTTP request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		clog.Error(err.Error())
		http.Error(w, "There was an error reading the request", http.StatusBadRequest)
		return
	}

	// Json unmarshal the request into the user.Profile
	var req CreateUserRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		clog.Error(err.Error())
		http.Error(w, "There was an error json unmarshaling the request", http.StatusBadRequest)
		return
	}

	// Validate that the profile is has all required info
	err = req.Profile.Validate()
	if err != nil {
		clog.Error(err.Error())
		http.Error(w, fmt.Sprintf("There was an error validating the request: %v", err), http.StatusBadRequest)
		return
	}

	// Get the password hash after validating it
	err = IsValidPassword(req.Password)
	if err != nil {
		clog.Error(err.Error())
		http.Error(w, fmt.Sprintf("The password is invalid: %v", err), http.StatusBadRequest)
		return
	}

	passwordHash, err := authLib.GetHash(req.Password, auth.PasswordSecretKey)
	if err != nil {
		clog.Error(err.Error())
		http.Error(w, rest.CleanAPIErrMessage, http.StatusInternalServerError)
		return
	}

	newUserReq := user.NewUserRequest{
		Profile:      req.Profile,
		PasswordHash: passwordHash,
	}

	// Update the profile of the given user
	usr, err := user.NewUser(newUserReq)
	if err == user.ErrEmailAlreadyExist {
		clog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err != nil {
		clog.Error(err.Error())
		http.Error(w, rest.CleanAPIErrMessage, http.StatusInternalServerError)
		return
	}

	// Json marshal the updated profile so we can send it back
	resp, err := json.Marshal(usr)
	if err != nil {
		clog.Error(err.Error())
		http.Error(w, rest.CleanAPIErrMessage, http.StatusInternalServerError)
		return
	}

	// Write the update profile to the http response
	w.Header().Set("Content-Type", "application/json")
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
	userID, err := authLib.GetUserIdFromRequest(r)
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
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(resp)
	if err != nil {
		clog.Error(err.Error())
		http.Error(w, rest.CleanAPIErrMessage, http.StatusInternalServerError)
	}

	clog.Info("Request succesfully processed")

}

// IsValidPassword validates that the password is good enough to be used
func IsValidPassword(password string) error {

	// password is not empty?
	if strings.TrimSpace(password) == "" {
		return fmt.Errorf("no password provided")
	}

	// password is not too short
	minLength := 6
	if len(password) < minLength {
		return fmt.Errorf("password is too short, needs a minimum of %d characters", minLength)
	}

	return nil
}
