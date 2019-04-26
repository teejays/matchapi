package handler

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/teejays/clog"

	"github.com/teejays/matchapi/db"
	"github.com/teejays/matchapi/service/user/v1"
)

func init() {
	// Let's turn logging off
	clog.LogLevel = 7
}
func TestHandleGetUser(t *testing.T) {

	// Initialize the mock DB client
	err := db.InitMockClient()
	if err != nil {
		t.Fatal(err)
	}
	defer db.DestoryMockClient()

	// Populate the User mock data in DB
	user.HelperPopulateMockData()

	// Setup the handler
	r := mux.NewRouter()

	r.PathPrefix("/{userid}").
		Subrouter().
		HandleFunc("/v1/user", HandleGetUser).
		Methods(http.MethodGet)

	http.Handle("/", r)
	defer func() { http.DefaultServeMux = new(http.ServeMux) }()
	tt := []struct {
		name             string
		target           string
		expectedCode     int
		expectedResponse *user.User
	}{
		{
			name:         "passing an invalid user auth id should say unauthorized",
			target:       "/42/v1/user",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:             "passing a valid user auth id should return it's own user",
			target:           "/1/v1/user",
			expectedCode:     http.StatusOK,
			expectedResponse: user.MockUsers[1],
		},
	}

	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {

			// Create the fake HTTP request
			req, err := http.NewRequest(http.MethodGet, test.target, nil)
			if err != nil {
				t.Fatal(err)
			}
			var w = httptest.NewRecorder()

			// Call the handler
			r.ServeHTTP(w, req)

			// Verify the status code
			assert.Equal(t, test.expectedCode, w.Code)

			// Verify the response
			if test.expectedResponse != nil {
				resp := w.Result()
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					t.Fatal(err)
				}
				var u user.User
				err = json.Unmarshal(body, &u)
				if err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, test.expectedResponse.Profile, u.Profile)

			}
		})
	}

}

func TestHandleCreatUser(t *testing.T) {

	// Initialize the mock DB client
	err := db.InitMockClient()
	if err != nil {
		t.Fatal(err)
	}
	defer db.DestoryMockClient()

	// Setup the handler
	r := mux.NewRouter()
	r.HandleFunc("/v1/user", HandleCreateUser).
		Methods(http.MethodPost)

	http.Handle("/", r)
	defer func() { http.DefaultServeMux = new(http.ServeMux) }()

	tt := []struct {
		name         string
		body         io.Reader
		expectedCode int
	}{
		{
			name:         "passing no data shoud give an error",
			body:         strings.NewReader(``),
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "passing an invalid user data should return an error",
			body:         strings.NewReader(`'{"FirstName":"","LastName":"", "Email": "", "Gender": 0}'`),
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "passing a valid user data should create a user and return it",
			body:         strings.NewReader(`{"FirstName":"Jon","LastName":"Harry", "Email": "jon.harry@email.com", "Gender": 3}`),
			expectedCode: http.StatusOK,
		},
	}

	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {

			// Create the fake HTTP request
			req, err := http.NewRequest(http.MethodPost, "/v1/user", test.body)
			if err != nil {
				t.Fatal(err)
			}
			var w = httptest.NewRecorder()

			// Call the handler
			r.ServeHTTP(w, req)

			// Verify the status code
			assert.Equal(t, test.expectedCode, w.Code)

			// Verify the response
			if test.expectedCode == http.StatusOK {
				resp := w.Result()
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					t.Fatal(err)
				}
				var u user.User
				err = json.Unmarshal(body, &u)
				if err != nil {
					t.Fatal(err)
				}
				assert.NotEqual(t, 0, u.ID)
				assert.Equal(t, "Jon", u.FirstName)

			}
		})
	}

}
