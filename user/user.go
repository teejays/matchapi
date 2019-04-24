package user

import (
	"errors"
	"time"
)

// UserID is a type alias to allow more control over the underlying data type of the UserID field
type UserID int

// Profile represents the editable part of the User
type Profile struct {
	FirstName string
	LastName  string
	Email     string
	Gender    int
}

// User represents the primary user object of the app
type User struct {
	ID              UserID
	DatetimeCreated time.Time
	DatetimeUpdated time.Time
	IsDeleted       bool
	Profile
}

// UpdateProfile ...
func (u *User) UpdateProfile(profile Profile) error {
	u.Profile = profile
	return nil
}

const (
	GenderMale int = iota
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
func GetUserByID(id UserID) (User, error) {
	user, err := getMockUserByID(id)
	return user, err
}

// UpdateUserByID ...
func UpdateUserByID(id UserID, profile Profile) (User, error) {
	var user User

	err := updateMockUserByID(id, profile)
	if err != nil {
		return user, err
	}

	return GetUserByID(id)
}

var mockUsers = map[UserID]User{
	1: User{
		ID:              1,
		DatetimeCreated: time.Now(),
		DatetimeUpdated: time.Now(),
		Profile: Profile{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@email.com",
			Gender:    GenderMale,
		},
	},
	2: User{
		ID:              2,
		DatetimeCreated: time.Now(),
		DatetimeUpdated: time.Now(),
		Profile: Profile{
			FirstName: "Jane",
			LastName:  "Doe",
			Email:     "jane.doe@email.com",
			Gender:    GenderFemale,
		},
	},
	3: User{
		ID:              3,
		DatetimeCreated: time.Now(),
		DatetimeUpdated: time.Now(),
		Profile: Profile{
			FirstName: "Jack",
			LastName:  "Does",
			Email:     "jack.does@email.com",
			Gender:    GenderOther,
		},
	},
}

func getMockUserByID(id UserID) (User, error) {
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

func updateMockUserByID(id UserID, profile Profile) error {
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
