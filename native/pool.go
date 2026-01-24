package native

import (
	"github.com/fufuok/bytespool"
)

var (
	// default: true
	reduceMemoryUsage = true

	// default byte slice pool: [32, 1MB]
	bspool = bytespool.NewCapacityPools(32, 1*1024*1024)
)

func InitPools(minSize, maxSize int) {
	if minSize > 0 && maxSize > 0 {
		bspool = bytespool.NewCapacityPools(minSize, maxSize)
	}
}

func SetWithStats(t bool) {
	bspool.SetWithStats(t)
}

func BytesPoolStats(topN int) bytespool.RuntimeSummary {
	return bytespool.RuntimeStatsSummary(topN, bspool)
}

func SetReduceMemoryUsage(b bool) {
	reduceMemoryUsage = b
}
