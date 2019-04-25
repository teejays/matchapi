package user

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/teejays/matchapi/db"
)

func TestHandleGetUser(t *testing.T) {

	// Initialize the mock DB client
	err := db.InitMockClient()
	if err != nil {
		t.Fatal(err)
	}
	defer db.DestoryMockClient()

	tt := []struct {
		name             string
		target           string
		expectedCode     int
		expectedResponse *User
	}{
		{
			name:         "passing an invalid user auth id should say unauthorized",
			target:       "/42/v1/user",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:             "passing a valid user auth id should return it's own user",
			target:           "/1/v1/user",
			expectedCode:     http.StatusUnauthorized,
			expectedResponse: mockUsers[1],
		},
	}

	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {

			// Create the fake HTTP request
			r, err := http.NewRequest(http.MethodGet, test.target, nil)
			if err != nil {
				t.Fatal(err)
			}
			var w = httptest.NewRecorder()

			// Call the handler
			handler := http.HandlerFunc(HandleGetUser)
			handler.ServeHTTP(w, r)

			// Verify the status code
			assert.Equal(t, test.expectedCode, w.Code)

			// Verify the response
			if test.expectedResponse != nil {
				resp := w.Result()
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					t.Fatal(err)
				}
				var u User
				err = json.Unmarshal(body, &u)
				if err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, test.expectedResponse.Profile, u.Profile)

			}
		})
	}

}
