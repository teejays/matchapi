package like

import (
	"time"

	"github.com/teejays/clog"

	"github.com/teejays/matchapi/lib/pk"
	"github.com/teejays/matchapi/service/v1/user"
)

// Like represents the like action
type Like struct {
	GiverID    pk.ID
	ReceiverID pk.ID
	Datetime   time.Time
}

type IncomingLike struct {
	Giver user.ShareableProfile
	Like
}

// GetIncomingLikesByUserID returns all the users that have like the provided UserID
func GetIncomingLikesByUserID(userID pk.ID) ([]IncomingLike, error) {
	var incomingLikes = []IncomingLike{}

	for _, l := range mockLikes {
		if l.ReceiverID == userID {
			giver, err := user.GetUserByID(l.GiverID)
			if err != nil {
				clog.Warnf("There was an error trying GetUserByID(%d): %v", l.GiverID, err)
				continue
			}
			var incomingLike IncomingLike
			incomingLike.Like = l
			incomingLike.Giver = giver.ShareableProfile
			incomingLikes = append(incomingLikes, incomingLike)
		}
	}
	return incomingLikes, nil
}

var mockLikes = map[int]Like{
	1: Like{
		GiverID:    2,
		ReceiverID: 1,
		Datetime:   time.Now(),
	},
	2: Like{
		GiverID:    3,
		ReceiverID: 1,
		Datetime:   time.Now(),
	},
	3: Like{
		GiverID:    3,
		ReceiverID: 2,
		Datetime:   time.Now(),
	},
}
