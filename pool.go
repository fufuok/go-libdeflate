package libdeflate

import (
	"github.com/fufuok/bytespool"
)

func InitDefaultPools(minSize, maxSize int) {
	bytespool.InitDefaultPools(minSize, maxSize)
}
