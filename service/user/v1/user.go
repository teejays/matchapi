package user

import (
	"errors"
	"fmt"
	"strings"
	"time"

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
	ID        pk.ID
	IsDeleted bool
	Profile
	meta
}

// meta represents the system generated and required info that we don't need to worry about
type meta struct {
	DatetimeCreated time.Time
	DatetimeUpdated time.Time
}

// Profile represents the editable part of the User
type Profile struct {
	ShareableProfile
	LastName string
	Email    string
}

// ShareableProfile is a part of the profile that can be shared with other users
type ShareableProfile struct {
	FirstName string
	Gender    int
	Images    []string
}

// NewUser creates a new instance of a user object and stores it in the database
func NewUser(profile Profile) (*User, error) {

	// Validate that user data is okay
	if err := profile.Validate(); err != nil {
		return nil, err
	}

	// Create a new user object and populate it with data
	var u User
	u.Profile = profile
	u.DatetimeCreated = time.Now()
	u.DatetimeUpdated = time.Now()

	// Save it to DB and get the new ID
	id, err := db.SaveNewEntity(db.UserCollection, &u)
	if err != nil {
		return nil, err
	}

	// Fetch the new entity from DB and return it
	var user User
	err = db.GetEntityByID(db.UserCollection, id, &user)
	if err != nil {
		return &user, err
	}

	return &user, nil
}

// GetUserByID returns the user object corresponding to the provided userID
func GetUserByID(id pk.ID) (*User, error) {
	var user User
	err := db.GetEntityByID(db.UserCollection, id, &user)
	if err != nil {
		return nil, err
	}

	return &user, err
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