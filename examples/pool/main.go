package main

import (
	"bytes"
	"fmt"

	"github.com/fufuok/bytespool"

	"github.com/fufuok/go-libdeflate"
)

func init() {
	// Optional: Set to true to minimize the capacity of returned slices, but adds one copy operation
	// libdeflate.SetReduceMemoryUsage(true)

	// Optional: Set the initial and maximum size of the byte pool (default 2B ~ 4MiB)
	// libdeflate.InitPools(32, 128<<20)
}

func main() {
	smallDataDemo()
	poolDemo()
}

func smallDataDemo() {
	c, _ := libdeflate.NewCompressorAutoClose()
	n, _, err := c.Compress(nil, nil, libdeflate.ModeZlib)
	fmt.Println(n, err) // 0 libdeflate: native: empty input

	dst, err := Zip(nil)
	fmt.Println(len(dst), err) // 0 <nil>

	data := []byte("f")
	n, _, err = c.Compress(data, nil, libdeflate.ModeZlib)
	fmt.Println(n, err) // 12 <nil>

	out, err := Zip(data)
	fmt.Println(len(out), err) // 12 <nil>
	bytespool.Put(out)
}

func poolDemo() {
	data := bytes.Repeat([]byte("abc.123"), 1000)

	dst, err := Zip(data)
	fmt.Println("zip err:", err)      // zip err: <nil>
	fmt.Println("zip len:", len(dst)) // zip len: 43

	src, err := Unzip(dst)
	fmt.Println("unzip err:", err)                     // unzip err: <nil>
	fmt.Println("src == data", bytes.Equal(src, data)) // src == data true

	// Suggestion (non-mandatory): After use, bytes can be put back to bytespool
	bytespool.Put(dst)
	bytespool.Put(src)

	dst, err = ZipLevel(data, libdeflate.MaxCompressionLevel)
	fmt.Println("zip err:", err)      // zip err: <nil>
	fmt.Println("zip len:", len(dst)) // zip len: 43

	src, err = Unzip(dst)
	fmt.Println("unzip err:", err)                     // unzip err: <nil>
	fmt.Println("src == data", bytes.Equal(src, data)) // src == data true

	// Suggestion (non-mandatory): After use, bytes can be put back to bytespool
	bytespool.Put(src)
	bytespool.Put(dst)
}
