package main

import (
	"sync"

	"github.com/fufuok/go-libdeflate"
)

var (
	zipPool   sync.Pool
	unZipPool sync.Pool
)

func init() {
	// Initialize compressor pool with default compression level
	zipPool.New = func() interface{} {
		c, err := libdeflate.NewCompressorAutoClose()
		if err != nil {
			panic(err)
		}
		return &c
	}

	// Initialize decompressor pool
	unZipPool.New = func() interface{} {
		dc, err := libdeflate.NewDecompressorAutoClose()
		if err != nil {
			panic(err)
		}
		return &dc
	}
}

// Zip compresses data using DefaultCompressionLevel
func Zip(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, nil
	}

	c := zipPool.Get().(*libdeflate.Compressor)
	defer zipPool.Put(c)

	_, dst, err := c.Compress(data, nil, libdeflate.ModeZlib)
	return dst, err
}

// ZipLevel compresses data with specified compression level, without using object pool
func ZipLevel(data []byte, level int) ([]byte, error) {
	if len(data) == 0 {
		return nil, nil
	}

	c, err := libdeflate.NewCompressorLevelAutoClose(level)
	if err != nil {
		return nil, err
	}
	_, dst, err := c.Compress(data, nil, libdeflate.ModeZlib)
	return dst, err
}

// Unzip decompresses data
func Unzip(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, nil
	}

	dc := unZipPool.Get().(*libdeflate.Decompressor)
	defer unZipPool.Put(dc)

	_, src, err := dc.Decompress(data, nil, libdeflate.ModeZlib)
	return src, err
}
