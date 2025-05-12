package utlities

import (
	"golang.org/x/sync/singleflight"
	"sync"
)

var PoolSingleFlight = sync.Pool{
	New: func() interface{} {
		return singleflight.Group{}
	},
}
