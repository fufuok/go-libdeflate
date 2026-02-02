package native

/*
#include "libdeflate.h"
#include "helper.h"
#include <stddef.h>
#include <stdlib.h>
#include <stdint.h>

typedef struct libdeflate_decompressor decomp;
*/
import "C"
import (
	"errors"
	"unsafe"
)

// Decompressor decompresses any DEFLATE, zlib or gzip compressed data at any level
type Decompressor struct {
	dc                     *C.decomp
	isClosed               bool
	maxDecompressionFactor int
}

// NewDecompressor returns a new Decompressor with maxDecompressionFactor = 32 or and error if out of memory
func NewDecompressor() (*Decompressor, error) {
	return NewDecompressorWithExtendedDecompression(32)
}

// NewDecompressorWithExtendedDecompression returns a new Decompressor with maxDecompressionFactor or and error if out of memory
func NewDecompressorWithExtendedDecompression(maxDecompressionFactor int) (*Decompressor, error) {
	dc := C.libdeflate_alloc_decompressor()
	if dc == nil {
		return nil, ErrorOutOfMemory
	}

	return &Decompressor{dc, false, maxDecompressionFactor}, nil
}

// Decompress decompresses the given data from in to out and returns out and an error if something went wrong.
// If error != nil, then the data in out is undefined.
// If you pass a buffer to out, the size of this buffer must exactly match the length of the decompressed data.
// If you pass nil as out, this function will allocate a sufficient buffer and return it.
// Returns the number of consumed bytes from 'in'
func (dc *Decompressor) Decompress(in, out []byte, f decompress) (int, []byte, error) {
	if dc.isClosed {
		panic(ErrorAlreadyClosed)
	}
	if len(in) == 0 {
		return 0, out, ErrorNoInput
	}

	if len(out) > 0 {
		cons, n, err := dc.decompress(in, out, true, f)
		return cons, out[:n], err
	}

	cons := 0
	n := 0
	decompFactor := 4
	if dc.maxDecompressionFactor < decompFactor {
		decompFactor = dc.maxDecompressionFactor
	}

	size := 32 + len(in)*decompFactor
	firstAttempt := true
	tryMaxSize := true
	err := ErrorInsufficientSpace
	for errors.Is(err, ErrorInsufficientSpace) {
		if !firstAttempt && len(out) > 0 {
			bspool.Put(out)
		}
		firstAttempt = false

		out = bspool.New(size)
		cons, n, err = dc.decompress(in, out, false, f)
		if err == nil {
			break
		}

		if decompFactor > dc.maxDecompressionFactor {
			if tryMaxSize && size < bspool.MaxSize() {
				tryMaxSize = false
				size = bspool.MaxSize()
				continue
			}

			outSmallCap := bspool.NewBytes(out[:n])
			bspool.Put(out)
			return cons, outSmallCap, ErrorInsufficientDecompressionFactor
		}

		decompFactor *= 2
		size = 32 + len(in)*decompFactor
	}

	outSmallCap := bspool.NewBytes(out[:n])
	bspool.Put(out)
	return cons, outSmallCap, err
}

func (dc *Decompressor) decompress(in, out []byte, fit bool, f decompress) (int, int, error) {
	inAddr := startMemAddr(in)
	outAddr := startMemAddr(out)

	var (
		cons int
		n    int
	)

	consPtr := uintptr(unsafe.Pointer(&cons))
	sPtr := uintptr(unsafe.Pointer(&n))
	if fit {
		sPtr = 0
	}

	err := f(dc.dc, inAddr, outAddr, len(in), len(out), consPtr, sPtr)

	if fit {
		n = len(out)
	}

	return cons, n, err
}

// Close frees the memory allocated by C objects
func (dc *Decompressor) Close() {
	if dc.isClosed {
		panic(ErrorAlreadyClosed)
	}
	C.libdeflate_free_decompressor(dc.dc)
	dc.isClosed = true
}

// PanicFreeClose is like Close but doesn't panic if the decompressor is already closed. This is useful for the higher-level autoclose functionality.
func (c *Decompressor) PanicFreeClose() {
	if c.isClosed {
		return
	}
	c.Close()
}
