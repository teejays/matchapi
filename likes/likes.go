package likes

import (
	"time"

	"github.com/teejays/matchapi/users"
)

// LikeAction represents the like action
type LikeAction struct {
	LikedBy  users.UserID
	Liked    users.UserID
	Datetime time.Time
}

// GetIncomingLikesByUserID returns all the users that have like the provided UserID
func GetIncomingLikesByUserID(userID users.UserID) ([]LikeAction, error) {
	var likes = []LikeAction{}
	for _, l := range mockLikes {
		if l.Liked == userID {
			likes = append(likes, l)
		}
	}
	return likes, nil
}

var mockLikes = map[int]LikeAction{
	1: LikeAction{
		Liked:    1,
		LikedBy:  2,
		Datetime: time.Now(),
	},
	2: LikeAction{
		Liked:    1,
		LikedBy:  3,
		Datetime: time.Now(),
	},
	3: LikeAction{
		Liked:    2,
		LikedBy:  3,
		Datetime: time.Now(),
	},
}
