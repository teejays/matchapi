package auth

import (
	"fmt"
	"context"
	"net/http"
	"crypto/hmac"
	"crypto/sha256"
	"strings"

	"github.com/teejays/matchapi/lib/pk"
	"github.com/teejays/matchapi/lib/auth/jwt"
)

type contextKey string
const ctxKeyForIsAuthenticated = contextKey("is_authenticated")
const ctxKeyForPayload = contextKey("jwt_payload")
const ctxKeyForToken = contextKey("jwt_token")

// TokenPayload is the payload type that goes in the JWT token
type TokenPayload struct {
	UserID pk.ID
	Email string
}

// NewPayload creates a new payload for the JWT token
func NewPayload(userID pk.ID, email string) (TokenPayload, error){
	var payload TokenPayload
	
	// Validate the params before making the payload
	if err := userID.Validate(); err != nil {
		return payload, err 
	}
	if strings.TrimSpace(email) == "" {
		return payload, fmt.Errorf("cannot create a paylaod with empty email")
	}

	payload = TokenPayload{
		UserID: userID,
		Email: email,
	}

	return payload, nil
}

// AuthenticateRequest should implement the authentication logic. It should should at the auth token
// and figure out what user context. Currently, this is not implemented and it only relies on
// and explicitly passed userID in the route.
func AuthenticateRequest(r *http.Request) (*http.Request, error) {

	// Get the authentication header
	val := r.Header.Get("Authorization")
	// In JWT, we're looking for the Bearer type token
	// This means that the val should be like: Bearer <token>
	// - split by the space
	valParts := strings.Split(val, " ")
	if len(valParts) != 2 {
		return r, fmt.Errorf("Authorization header has an invalid form: it's not `Authorization:Bearer <token>")
	}
	if valParts[0] != "Bearer" {
		return r, fmt.Errorf("Authorization header has an invalid form: it's not `Authorization:Bearer <token>")
	}

	token := valParts[1]

	cl, err := jwt.GetClient()
	if err != nil {
		return r, err
	}

	var payload TokenPayload
	err = cl.VerifyAndDecode(token, &payload)
	if err != nil {
		return r, err
	}

	// Authentication succesful
	// Add the authentication payload to the context
	ctx := r.Context()
	ctx = context.WithValue(ctx, ctxKeyForIsAuthenticated, true)
	ctx = context.WithValue(ctx, ctxKeyForPayload, payload)
	ctx = context.WithValue(ctx, ctxKeyForToken, token)
	
	return r.WithContext(ctx), nil

}

func IsRequestAuthenticated(r *http.Request) bool {
	defer func() {
		if r := recover(); r != nil {
			panic(fmt.Sprintf("Could not verify whether http.Request is authenticated: %v", r))
		}
	}()
	
	// Not authenticated if the request itself is nil
	if r == nil {
		return false
	}

	ctx := r.Context()
	isAuthenticated := (ctx.Value(ctxKeyForIsAuthenticated)).(bool)
	
	return isAuthenticated

}

var ErrNotAuthenticated = fmt.Errorf("not authenticated")

func GetPayloadFromRequest(r *http.Request) (TokenPayload, error) {
	defer func() {
		if r := recover(); r != nil {
			panic(fmt.Sprintf("Could not extract payload from http.Request: %v", r))
		}
	}()

	var payload TokenPayload

	if !IsRequestAuthenticated(r) {
		return payload, ErrNotAuthenticated
	}

	ctx := r.Context()
	payload = ctx.Value(ctxKeyForPayload).(TokenPayload)

	return payload, nil
}

func GetUserIdFromRequest(r *http.Request) (pk.ID, error) {
	var userID pk.ID
	payload, err := GetPayloadFromRequest(r)
	if err != nil {
		return userID, err
	}
	
	return payload.UserID, nil

}

// IsStrongPassword validates that the password is good enough to be used
func IsStrongPassword(password string) error {

	// password is not empty?
	if password != "" {
		return fmt.Errorf("password is empty")
	}

	// password is not too short
	minLength := 10
	if len(password) < minLength {
		return fmt.Errorf("password is too short, needs a minimum of %d characters", minLength)
	}

	return nil
}

const passwordSecretKey = "I am disco dancer"
// GetHash returns the hash of the message
func GetHash(message string) (string, error) {
	h, err := hash([]byte(message))
	return string(h), err
}

func hash(message []byte) ([]byte, error) {
	hash := hmac.New(sha256.New, []byte(passwordSecretKey))
	_, err := hash.Write(message)
	if err != nil {
		return nil, err
	}
	return hash.Sum(message), nil
}
