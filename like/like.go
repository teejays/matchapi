package like

import (
	"time"

	"github.com/teejays/matchapi/user"
)

// Like represents the like action
type Like struct {
	LikedBy  user.UserID
	Liked    user.UserID
	Datetime time.Time
}

// GetIncomingLikesByUserID returns all the users that have like the provided UserID
func GetIncomingLikesByUserID(userID user.UserID) ([]Like, error) {
	var likes = []Like{}
	for _, l := range mockLikes {
		if l.Liked == userID {
			likes = append(likes, l)
		}
	}
	return likes, nil
}

var mockLikes = map[int]Like{
	1: Like{
		Liked:    1,
		LikedBy:  2,
		Datetime: time.Now(),
	},
	2: Like{
		Liked:    1,
		LikedBy:  3,
		Datetime: time.Now(),
	},
	3: Like{
		Liked:    2,
		LikedBy:  3,
		Datetime: time.Now(),
	},
}