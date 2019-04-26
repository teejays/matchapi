package like

import (
	"time"

	"github.com/teejays/matchapi/lib/pk"
	likev2 "github.com/teejays/matchapi/service/like/v2"
)

// Like represents the like action
type Like struct {
	ID       pk.ID
	LikedBy  pk.ID `json:"GiverID"`
	Liked    pk.ID `json:"ReceiverID"`
	Datetime time.Time
}

// ConvertV2IncomimgLikeToV1 ...
func ConvertV2IncomimgLikeToV1(l likev2.IncomingLike) Like {
	return Like{
		ID:       l.ID,
		LikedBy:  l.GiverID,
		Liked:    l.ReceiverID,
		Datetime: l.Datetime,
	}
}
