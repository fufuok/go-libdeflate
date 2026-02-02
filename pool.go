package libdeflate

import (
	"github.com/fufuok/bytespool"

	"github.com/fufuok/go-libdeflate/native"
)

// InitPools initialize the bytes pools, default: 32B ~ 1MiB
func InitPools(minSize, maxSize int) {
	native.InitPools(minSize, maxSize)
}

func SetWithStats(t bool) {
	native.SetWithStats(t)
}

func BytesPoolStats(topN int) bytespool.RuntimeSummary {
	return native.BytesPoolStats(topN)
}
