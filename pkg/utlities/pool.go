package utlities

import (
	"golang.org/x/sync/singleflight"
	"sync"
)

var SingleFlightPool = sync.Pool{
	New: func() interface{} {
		return new(singleflight.Group)
	},
}
