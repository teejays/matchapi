package like

import (
	"time"
)

var mockLikes = map[int]Like{
	1: Like{
		GiverID:  2,
		Datetime: time.Now(),
		BasicLike: BasicLike{
			ReceiverID: 1,
		},
	},
	2: Like{
		GiverID:  3,
		Datetime: time.Now(),
		BasicLike: BasicLike{
			ReceiverID: 1,
		},
	},
	3: Like{
		GiverID:  3,
		Datetime: time.Now(),
		BasicLike: BasicLike{
			ReceiverID: 2,
		},
	},
}
