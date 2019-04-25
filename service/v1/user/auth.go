package user

import (
	"fmt"
	"net/http"

	"github.com/teejays/clog"

	"github.com/teejays/matchapi/lib/pk"
	"github.com/teejays/matchapi/lib/rest"
)

// Authenticate should implement the authentication logic. It should should at the auth token
// and figure out what user context. Currently, this is not implemented and it only relies on
// and explicitly passed userID in the route.
func Authenticate(r *http.Request) (pk.ID, error) {

	// Get the passed userID (token)
	id, err := rest.GetUserIDMuxVar(r)
	if err != nil {
		return id, err
	}

	// Validate that the user exists
	user, err := GetUserByID(id)
	if err != nil {
		clog.Errorf("%v", err)
		return id, fmt.Errorf("error while fetching the authorization user")
	}

	if user == nil {
		return id, fmt.Errorf("could not find a valid user with the provided ID")
	}

	return user.ID, nil
}
