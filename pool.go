package libdeflate

import (
	"github.com/fufuok/go-libdeflate/native"
)

// InitPools initialize the bytes pools, default: 32B ~ 4MiB
func InitPools(minSize, maxSize int) {
	native.InitPools(minSize, maxSize)
}

// SetReduceMemoryUsage For: Return smaller slices when compressing #22
func SetReduceMemoryUsage(b bool) {
	native.SetReduceMemoryUsage(b)
}
