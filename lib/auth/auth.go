package auth

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/teejays/clog"

	"github.com/teejays/matchapi/lib/auth/jwt"
	"github.com/teejays/matchapi/lib/pk"
)

var JWTSecretKey = "JWT is awesome!"

type contextKey string

const ctxKeyForIsAuthenticated = contextKey("is_authenticated")
const ctxKeyForPayload = contextKey("jwt_payload")
const ctxKeyForToken = contextKey("jwt_token")

// TokenPayload is the payload type that goes in the JWT token
type TokenPayload struct {
	UserID pk.ID
	Email  string
}

// NewPayload creates a new payload for the JWT token
func NewPayload(userID pk.ID, email string) (TokenPayload, error) {
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
		Email:  email,
	}

	return payload, nil
}

// AuthenticateRequest should implement the authentication logic. It should should at the auth token
// and figure out what user context. Currently, this is not implemented and it only relies on
// and explicitly passed userID in the route.
func AuthenticateRequest(r *http.Request) (*http.Request, error) {
	clog.Debug("AuthenticateRequest() called...")
	// Get the authentication header
	val := r.Header.Get("Authorization")
	clog.Debugf("Authenticate Header: %v", val)
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

	// TODO: What if JWT client isn't initialized?
	if !jwt.IsClientInitialized() {
		err := InitJWTClient()
		if err != nil {
			return r, err
		}
	}
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

// GetHash returns the hash of the message
func GetHash(message, secret string) ([]byte, error) {
	h, err := hash([]byte(message), []byte(secret))
	clog.Debugf("Password Hash: %v", h)
	return h, err
}

func hash(message, secret []byte) ([]byte, error) {
	hash := hmac.New(sha256.New, secret)
	_, err := hash.Write(message)
	if err != nil {
		return nil, err
	}
	return hash.Sum(nil), nil
}

// IsEqualHash compares if two hashes are equal
func IsEqualHash(h1, h2 []byte) bool {
	clog.Debugf("H1: %v\nH2: %v", h1, h2)
	return hmac.Equal(h1, h2)
}

func InitJWTClient() error {
	return jwt.InitClient(JWTSecretKey, time.Hour*48)
}
