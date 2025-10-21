package libdeflate

import (
	"github.com/fufuok/bytespool"

	"github.com/fufuok/go-libdeflate/native"
)

// InitDefaultPools Initialize the default pools, default: 2B ~ 8MiB
func InitDefaultPools(minSize, maxSize int) {
	bytespool.InitDefaultPools(minSize, maxSize)
}

// SetReduceMemoryUsage For: Return smaller slices when compressing #22
func SetReduceMemoryUsage(b bool) {
	native.SetReduceMemoryUsage(b)
}
