package user

import (
	"time"

	"github.com/teejays/clog"

	"github.com/teejays/matchapi/db"
	"github.com/teejays/matchapi/lib/pk"
)

// HelperPopulateMockData populates the mock DB client with test users
func HelperPopulateMockData() error {
	clog.Debugf("Populating the mock users data")
	for id, u := range MockUsers {
		err := db.SaveEntityByID(db.UserCollection, id, u)
		if err != nil {
			return err
		}
	}
	return nil
}

// MockUsers is the test user data
var MockUsers = map[pk.ID]*User{
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
