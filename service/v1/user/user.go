package user

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/teejays/matchapi/lib/pk"
)

// ShareableProfile is a part of the profile that can be shared with other users
type ShareableProfile struct {
	FirstName string
	Gender    int
	Images    []string
}

// Profile represents the editable part of the User
type Profile struct {
	ShareableProfile
	LastName string
	Email    string
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

// User represents the primary user object of the app
type User struct {
	ID              pk.ID
	DatetimeCreated time.Time
	DatetimeUpdated time.Time
	IsDeleted       bool
	Profile
}

// UpdateProfile ...
func (u *User) UpdateProfile(profile Profile) error {
	if err := profile.Validate(); err != nil {
		return err
	}
	u.Profile = profile
	return nil
}

const (
	GenderInvalid int = iota
	GenderMale
	GenderFemale
	GenderOther
)

// ErrEntityDoesNotExist is used when the requested entity does not exist in the system
var ErrEntityDoesNotExist = errors.New("the requested entity does not exist")

// GetUsers returns all the users in the system
func GetUsers() ([]User, error) {
	var users = getMockUsers()
	return users, nil
}

// GetUserByID returns the user object corresponding to the provided userID
func GetUserByID(id pk.ID) (User, error) {
	user, err := getMockUserByID(id)
	return user, err
}

// UpdateUserByID ...
func UpdateUserByID(id pk.ID, profile Profile) (User, error) {
	var user User

	err := updateMockUserByID(id, profile)
	if err != nil {
		return user, err
	}

	return GetUserByID(id)
}

var mockUsers = map[pk.ID]User{
	1: User{
		ID:              1,
		DatetimeCreated: time.Now(),
		DatetimeUpdated: time.Now(),
		Profile: Profile{
			ShareableProfile: ShareableProfile{
				FirstName: "John",
				Gender:    GenderMale,
			},
			LastName: "Doe",
			Email:    "john.doe@email.com",
		},
	},
	2: User{
		ID:              2,
		DatetimeCreated: time.Now(),
		DatetimeUpdated: time.Now(),
		Profile: Profile{
			ShareableProfile: ShareableProfile{
				FirstName: "Jane",
				Gender:    GenderFemale,
			},
			LastName: "Doe",
			Email:    "jane.doe@email.com",
		},
	},
	3: User{
		ID:              3,
		DatetimeCreated: time.Now(),
		DatetimeUpdated: time.Now(),
		Profile: Profile{
			ShareableProfile: ShareableProfile{
				FirstName: "Jack",
				Gender:    GenderOther,
			},
			LastName: "Does",
			Email:    "jack.does@email.com",
		},
	},
}

func getMockUserByID(id pk.ID) (User, error) {
	var user User
	var exists bool
	user, exists = mockUsers[id]
	if !exists {
		return user, ErrEntityDoesNotExist
	}
	return user, nil
}

func getMockUsers() []User {
	var users []User
	for _, usr := range mockUsers {
		users = append(users, usr)
	}
	return users

}

func updateMockUserByID(id pk.ID, profile Profile) error {
	user, exists := mockUsers[id]
	if !exists {
		return ErrEntityDoesNotExist
	}

	err := (&user).UpdateProfile(profile)
	if err != nil {
		return err
	}

	mockUsers[id] = user

	return nil
}
