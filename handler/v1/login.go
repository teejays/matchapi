package handler

import (
	"strings"
	"io/ioutil"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/teejays/clog"
	
	"github.com/teejays/matchapi/service/auth/v1"
	"github.com/teejays/matchapi/lib/rest"
)

type LoginRequest struct {
	Email string
	Password string
}

// HandleLogin logins a user 
func HandleLogin(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// Read the HTTP request body
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		clog.Error(err.Error())
		http.Error(w, "There was an error reading the request", http.StatusBadRequest)
		return
	}

	// Json unmarshal into type LoginRequest
	var creds auth.LoginRequest
	err = json.Unmarshal(data, &creds)
	if err != nil {
		clog.Error(err.Error())
		http.Error(w, "There was an error json unmarshalling the request", http.StatusBadRequest)
		return
	}

	// Validate the request
	if strings.TrimSpace(creds.Email) == "" {
		clog.Error(err.Error())
		http.Error(w, fmt.Sprintf("There was an error validating the request: %s", "email cannot be empty"), http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(creds.Password) == "" {
		clog.Error(err.Error())
		http.Error(w, fmt.Sprintf("There was an error validating the request: %s", "password cannot be empty"), http.StatusBadRequest)
		return
	}

	// Find the user with these creds?
	token, err := auth.Login(creds)
	if err != nil {
		clog.Error(err.Error())
		http.Error(w, rest.CleanAPIErrMessage, http.StatusInternalServerError)
		return
	}

	type loginResponse struct {
		Token string
	}
	var respJSON loginResponse
	respJSON.Token = token

	// Json marshal the response
	resp, err := json.Marshal(respJSON)
	if err != nil {
		clog.Error(err.Error())
		http.Error(w, rest.CleanAPIErrMessage, http.StatusInternalServerError)
		return
	}

	// Write the response
	w.Header().Set("Content-Type", "application/json")

	_, err = w.Write(resp)
	if err != nil {
		clog.Error(err.Error())
		http.Error(w, rest.CleanAPIErrMessage, http.StatusInternalServerError)
	}

	_, err = w.Write([]byte(resp))
	if err != nil {
		clog.Error(err.Error())
		http.Error(w, rest.CleanAPIErrMessage, http.StatusInternalServerError)
		return
	}

}