package like

import (
	"time"
)

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
