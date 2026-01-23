package native

import (
	"github.com/fufuok/bytespool"
)

var (
	reduceMemoryUsage bool

	// default byte slice pool: [32, 4MB]
	bspool = bytespool.NewCapacityPools(32, 4*1024*1024)
)

func InitPools(minSize, maxSize int) {
	if minSize > 0 && maxSize > 0 {
		bspool = bytespool.NewCapacityPools(minSize, maxSize)
	}
}

func SetReduceMemoryUsage(b bool) {
	reduceMemoryUsage = b
}
