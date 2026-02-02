package libdeflate

import "errors"

var (
	ErrorInvalidModeCompressor   = errors.New("libdeflate: compressor: invalid mode")
	ErrorInvalidModeDecompressor = errors.New("libdeflate: decompressor: invalid mode")
)
