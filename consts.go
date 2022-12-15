package libdeflate

import "github.com/fufuok/go-libdeflate/native"

// These constants specify several special compression levels
const (
	MinCompressionLevel        = native.MinCompressionLevel
	MaxStdZlibCompressionLevel = native.MaxStdZlibCompressionLevel
	MaxCompressionLevel        = native.MaxCompressionLevel
	DefaultCompressionLevel    = native.DefaultCompressionLevel
)
