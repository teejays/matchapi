package user

import (
	"errors"
	"time"
)

// UserID is a type alias to allow more control over the underlying data type of the UserID field
type UserID int

// User represents the primary user object of the app
type User struct {
	ID              UserID
	FirstName       string
	LastName        string
	Email           string
	Gender          int
	DatetimeCreated time.Time
	DatetimeUpdated time.Time
	IsDeleted       bool
}

const (
	GenderMale int = iota
	GenderFemale
	GenderOther
)

// ErrEntityDoesNotExist is used when the requested entity does not exist in the system
var ErrEntityDoesNotExist error = errors.New("the requested entity does not exist")

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

var mockUsers = map[UserID]User{
	1: User{
		ID:              1,
		FirstName:       "John",
		LastName:        "Doe",
		Email:           "john.doe@email.com",
		Gender:          GenderMale,
		DatetimeCreated: time.Now(),
		DatetimeUpdated: time.Now(),
	},
	2: User{
		ID:              2,
		FirstName:       "Jane",
		LastName:        "Doe",
		Email:           "jane.doe@email.com",
		Gender:          GenderFemale,
		DatetimeCreated: time.Now(),
		DatetimeUpdated: time.Now(),
	},
	3: User{
		ID:              3,
		FirstName:       "Jack",
		LastName:        "Does",
		Email:           "jack.does@email.com",
		Gender:          GenderOther,
		DatetimeCreated: time.Now(),
		DatetimeUpdated: time.Now(),
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
