package auth

import (
	"fmt"

	"github.com/teejays/matchapi/lib/auth"
	"github.com/teejays/matchapi/lib/auth/jwt"
	"github.com/teejays/matchapi/service/user/v1"
)

// LoginRequest holds the login information
type LoginRequest struct {
	Email    string
	Password string
}

var ErrInvalidEmail = fmt.Errorf("no accounts found with the given email")
var ErrInvalidPassword = fmt.Errorf("accounts found but could not match the password")

// Login takes a LoginRequest and verifies that login credentials
func Login(req LoginRequest) (string, error) {

	// find the user by email ID
	creds, err := user.GetUserCredsByEmail(req.Email)
	if err != nil {
		return "", err
	}

	if len(creds) < 1 {
		return "", ErrInvalidEmail
	}

	// there should be only one user with this email
	if len(creds) > 1 {
		return "", fmt.Errorf("email %s has mutiple accounts", req.Email)
	}

	// hash the password
	h, err := auth.GetHash(req.Password)
	if err != nil {
		return "", err
	}

	// Find the users (we're treating it as multiple users here)
	var u *user.User
	var found bool
	for _, c := range creds {
		if h == c.PasswordHash && req.Email == c.Email {
			// We have a user!
			u, err = user.GetUserByID(c.ID)
			if err != nil {
				return "", err
			}
			found = true
			break
		}
	}

	if !found {
		return "", ErrInvalidPassword
	}

	// Get the JWT client and create a token 
	cl, err := jwt.GetClient()
	if err != nil {
		return "", err
	}

	payload, err := auth.NewPayload(u.ID, u.Email)
	if err != nil {
		return "", fmt.Errorf("error creating payload for JWT token: %v", err)
	}

	token, err := cl.CreateToken(payload)
	if err != nil {
		return "", fmt.Errorf("error creating JWT token: %v", err)
	}

	return token, nil
}