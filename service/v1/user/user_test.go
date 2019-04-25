package user

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/teejays/clog"

	"github.com/teejays/matchapi/db"
	"github.com/teejays/matchapi/lib/pk"
)

func init() {
	clog.LogLevel = 0
}

func helperPopulateMockData() error {
	clog.Debugf("Populating the mock users data")
	for id, u := range mockUsers {
		err := db.SaveEntityByID(db.UserCollection, id, u)
		if err != nil {
			return err
		}
	}
	return nil
}

func TestValidate(t *testing.T) {

	// Define the table tests
	tt := []struct {
		name      string
		profile   Profile
		shouldErr bool
	}{
		{
			name: "empty first name should give an error",
			profile: Profile{
				ShareableProfile: ShareableProfile{"", GenderMale, []string{}},
				LastName:         "Doe",
				Email:            "jon.doe@email.com",
			},
			shouldErr: true,
		},
		{
			name: "empty last name should give an error",
			profile: Profile{
				ShareableProfile: ShareableProfile{"Jon", GenderMale, []string{}},
				LastName:         "",
				Email:            "jon.doe@email.com",
			},
			shouldErr: true,
		},
		{
			name: "empty email should give an error",
			profile: Profile{
				ShareableProfile: ShareableProfile{"Jon", GenderMale, []string{}},
				LastName:         "Doe",
				Email:            "",
			},
			shouldErr: true,
		},
		{
			name: "invalid gender should give an error",
			profile: Profile{
				ShareableProfile: ShareableProfile{"Jon", 7, []string{}},
				LastName:         "Doe",
				Email:            "jon.doe@email.com",
			},
			shouldErr: true,
		},
		{
			name: "valid profile should not give an error",
			profile: Profile{
				ShareableProfile: ShareableProfile{"Jon", GenderMale, []string{}},
				LastName:         "Doe",
				Email:            "jon.doe@email.com",
			},
			shouldErr: false,
		},
	}

	// Run the table tests
	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			err := test.profile.Validate()
			assert.Equal(t, err != nil, test.shouldErr)
		})
	}
}

func TestNewUser(t *testing.T) {

	// Define the table tests
	tt := []struct {
		name      string
		profile   Profile
		shouldErr bool
	}{
		{
			name: "invalid profile should give an error",
			profile: Profile{
				ShareableProfile: ShareableProfile{"", GenderInvalid, []string{}},
				LastName:         "",
				Email:            "",
			},
			shouldErr: true,
		},
		{
			name: "valid profile should not give an error",
			profile: Profile{
				ShareableProfile: ShareableProfile{"Jon", GenderMale, []string{}},
				LastName:         "Doe",
				Email:            "jon.doe@email.com",
			},
			shouldErr: false,
		},
	}

	// Initialize the mock DB client
	err := db.InitMockClient()
	if err != nil {
		t.Error(err)
	}
	defer db.DestoryMockClient()

	// Run the table tests
	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			user, err := NewUser(test.profile)
			assert.Equal(t, err != nil, test.shouldErr)
			if !test.shouldErr {
				assert.NotEqual(t, 0, user.ID)
				assert.Equal(t, test.profile, user.Profile)
			}
		})
	}
}

func TestGetUserByID(t *testing.T) {

	// Initialize the mock DB client
	err := db.InitMockClient()
	if err != nil {
		t.Fatal(err)
	}
	defer db.DestoryMockClient()

	// Populate some data
	err = helperPopulateMockData()
	if err != nil {
		t.Fatal(err)
	}

	// Define table tests
	tt := []struct {
		name      string
		id        pk.ID
		user      *User
		shouldErr bool
	}{
		{
			name:      "fecthing a non-existent User should give an error",
			id:        42,
			user:      nil,
			shouldErr: true,
		},
		{
			name:      "fecthing an existing profile should return the right profile",
			id:        1,
			user:      mockUsers[1],
			shouldErr: false,
		},
	}

	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			u, err := GetUserByID(test.id)
			assert.Equal(t, test.shouldErr, err != nil)
			if test.user != nil {
				assert.NotNil(t, u)
				assert.Equal(t, test.user.ID, u.ID)
				assert.Equal(t, test.user.Profile, u.Profile)
			}
		})
	}
}

func TestUpdateProfile(t *testing.T) {

	// Initialize the mock DB client
	err := db.InitMockClient()
	if err != nil {
		t.Fatal(err)
	}
	defer db.DestoryMockClient()

	// Populate some data
	err = helperPopulateMockData()
	if err != nil {
		t.Fatal(err)
	}

	// Define table tests
	tt := []struct {
		name       string
		id         pk.ID
		newProfile Profile
		shouldErr  bool
	}{
		{
			name: "updating to an invalid profile should give an error",
			id:   1,
			newProfile: Profile{
				ShareableProfile: ShareableProfile{
					FirstName: "",
					Gender:    GenderInvalid,
				},
				LastName: "",
				Email:    "",
			},
			shouldErr: true,
		},
		{
			name:       "updating to a valid profile should be okay",
			id:         1,
			newProfile: mockUsers[2].Profile,
			shouldErr:  false,
		},
	}

	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			// Get the user
			u, err := GetUserByID(test.id)
			if err != nil {
				t.Fatal(err)
			}
			// Update the profile
			err = u.UpdateProfile(test.newProfile)
			assert.Equal(t, test.shouldErr, err != nil)

			if !test.shouldErr {
				assert.NotNil(t, u)
				assert.Equal(t, test.newProfile, u.Profile)
				// Get the profile from DB and ensure it's updated
				u2, err := GetUserByID(test.id)
				if err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, u.Profile, u2.Profile)
			}
		})
	}
}

var mockUsers = map[pk.ID]*User{
	1: &User{
		ID: 1,
		meta: meta{
			DatetimeCreated: time.Now(),
			DatetimeUpdated: time.Now(),
		},
		Profile: Profile{
			ShareableProfile: ShareableProfile{
				FirstName: "John",
				Gender:    GenderMale,
			},
			LastName: "Doe",
			Email:    "john.doe@email.com",
		},
	},
	2: &User{
		ID: 2,
		meta: meta{
			DatetimeCreated: time.Now(),
			DatetimeUpdated: time.Now(),
		},
		Profile: Profile{
			ShareableProfile: ShareableProfile{
				FirstName: "Jane",
				Gender:    GenderFemale,
			},
			LastName: "Doe",
			Email:    "jane.doe@email.com",
		},
	},
	3: &User{
		ID: 3,
		meta: meta{
			DatetimeCreated: time.Now(),
			DatetimeUpdated: time.Now(),
		},
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
