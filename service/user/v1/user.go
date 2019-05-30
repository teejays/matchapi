package user

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/teejays/clog"

	"github.com/teejays/matchapi/db"
	"github.com/teejays/matchapi/lib/pk"
)

const (
	GenderInvalid int = iota
	GenderMale
	GenderFemale
	GenderOther
)

// ErrEntityDoesNotExist is used when the requested entity does not exist in the system
var ErrEntityDoesNotExist = errors.New("the requested entity does not exist")

// User represents the primary user object of the app
type User struct {
	ID pk.ID
	Profile
	PasswordHash []byte
	IsDeleted    bool
	meta
}

// meta represents the system generated and required info that we don't need to worry about
type meta struct {
	DatetimeCreated time.Time
	DatetimeUpdated time.Time
}

// ProfileUser represents the editable part of the User with ID
type ProfileUser struct {
	ID pk.ID
	Profile
}

// Profile represents the editable part of the User
type Profile struct {
	ShareableProfile
	LastName string
	Email    string
}

// ShareableProfileUser is a part of the profile that can be shared with other users
type ShareableProfileUser struct {
	ID pk.ID
	ShareableProfile
}

// ShareableProfile is a part of the profile that can be shared with other users
type ShareableProfile struct {
	FirstName string
	Gender    int
	Images    []string
}

// NewUserRequest represents that request object that is used when
// a new user is created
type NewUserRequest struct {
	Profile
	PasswordHash []byte
}

var ErrEmailAlreadyExist = fmt.Errorf("Email is already taken")

// NewUser creates a new instance of a user object and stores it in the database
func NewUser(req NewUserRequest) (*ProfileUser, error) {

	// Validate that user data is okay
	if err := req.Profile.Validate(); err != nil {
		return nil, err
	}

	// Make sure we don't have a user already with the email
	users, err := GetUserCredsByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if len(users) > 0 {
		return nil, ErrEmailAlreadyExist
	}
	// Create a new user object and populate it with data
	var u User
	u.Profile = req.Profile
	u.PasswordHash = req.PasswordHash
	u.DatetimeCreated = time.Now()
	u.DatetimeUpdated = time.Now()

	clog.Debugf("NewUser Hash: %v", u.PasswordHash)
	// Save it to DB and get the new ID
	id, err := db.SaveNewEntity(db.UserCollection, &u)
	if err != nil {
		return nil, err
	}
	clog.Debugf("User | NewUser(): new ID generates: %d", id)

	// Fetch the new entity from DB and return it
	var user User
	err = db.GetEntityByID(db.UserCollection, id, &user)
	if err != nil {
		return nil, err
	}
	// We can't/shouldn't return the entire user object for security purposes

	safeUser := ProfileUser{
		ID:      user.ID,
		Profile: user.Profile,
	}
	return &safeUser, nil
}

// GetUserByID returns the user object corresponding to the provided userID
func GetUserByID(id pk.ID) (*User, error) {
	var user User
	err := db.GetEntityByID(db.UserCollection, id, &user)
	clog.Debugf("user.GetUserById(%v): \nuser: %v\nerr:%v", id, user, err)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

type UserCred struct {
	ID           pk.ID
	Email        string
	PasswordHash []byte
}

func GetUserCredsByEmail(email string) ([]UserCred, error) {
	// Run the query to all the users
	results, err := db.Query(db.UserCollection, fmt.Sprintf("Email:%s", email))
	if err != nil {
		return nil, err
	}

	// Get the creds and return
	var creds []UserCred
	for _, _v := range results {
		v, ok := _v.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("error fetching users by email: one of the response is not a map[string]interface{}")
		}
		id, ok := v["ID"].(float64)
		if !ok {
			return nil, fmt.Errorf("error fetching users by email: one of the response ID is not a number")
		}

		var c UserCred
		usr, err := GetUserByID(pk.ID(id))
		if err != nil {
			return creds, err
		}

		c.ID = usr.ID
		c.Email = usr.Email
		c.PasswordHash = usr.PasswordHash

		creds = append(creds, c)
	}

	return creds, nil
}

// UpdateProfile ...
func (u *User) UpdateProfile(profile Profile) error {
	if err := profile.Validate(); err != nil {
		return err
	}
	u.Profile = profile
	u.DatetimeUpdated = time.Now()

	err := db.SaveEntityByID(db.UserCollection, u.ID, u)
	return err
}

// Validate validates a profile before saving
func (p Profile) Validate() error {
	var errs []error

	if strings.TrimSpace(p.FirstName) == "" {
		errs = append(errs, fmt.Errorf("first name cannot be empty"))
	}
	if strings.TrimSpace(p.LastName) == "" {
		errs = append(errs, fmt.Errorf("last name cannot be empty"))
	}
	if strings.TrimSpace(p.Email) == "" {
		errs = append(errs, fmt.Errorf("email cannot be empty"))
	}
	if p.Gender < 1 || p.Gender > 3 {
		errs = append(errs, fmt.Errorf("gender is invalid; possible values are 1 (male), 2 (female) and 3 (other)"))
	}

	if len(errs) > 0 {
		var errMsg string
		for i, err := range errs {
			errMsg = fmt.Sprintf("%s\n%d) %v", errMsg, i+1, err)
		}
		return fmt.Errorf(errMsg)
	}

	return nil
}
