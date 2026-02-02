package helper

import (
	"bytes"
	"fmt"
	"math/rand"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/fufuok/go-libdeflate"
)

// Compressor is an alias for libdeflate.Compressor for test purposes
type Compressor = libdeflate.Compressor

// TestZipUnzip_Basic tests basic zip/unzip functionality
func TestZipUnzip_Basic(t *testing.T) {
	testCases := []struct {
		name        string
		data        []byte
		expectedErr bool
	}{
		{"empty data", []byte{}, false},
		{"small data", []byte("hello world"), false},
		{"medium data", bytes.Repeat([]byte("test"), 100), false},
		{"large data", bytes.Repeat([]byte("test data"), 1000), false},
		{"binary data", []byte{0x00, 0x01, 0x02, 0x03, 0xff}, false},
		{"unicode data", []byte("你好，世界！Hello World!"), false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test Zip
			compressed, err := Zip(tc.data)
			if (err != nil) != tc.expectedErr {
				t.Errorf("Zip() error = %v, expectedErr %v", err, tc.expectedErr)
			}

			if len(tc.data) == 0 {
				if compressed != nil {
					t.Errorf("Zip() for empty data = %v, want nil", compressed)
				}
				return
			}

			// Test Unzip
			decompressed, err := Unzip(compressed)
			if (err != nil) != tc.expectedErr {
				t.Errorf("Unzip() error = %v, expectedErr %v", err, tc.expectedErr)
			}

			// Verify data integrity
			if !bytes.Equal(tc.data, decompressed) {
				t.Errorf("Unzip() = %v, want %v", decompressed, tc.data)
			}

			// Test ZipLevel
			for _, level := range []int{1, 6, 9} {
				compressedLevel, err := ZipLevel(tc.data, level)
				if err != nil {
					t.Errorf("ZipLevel() level %d error = %v", level, err)
				}

				decompressedLevel, err := Unzip(compressedLevel)
				if err != nil {
					t.Errorf("Unzip() for level %d error = %v", level, err)
				}

				if !bytes.Equal(tc.data, decompressedLevel) {
					t.Errorf("Unzip() for level %d = %v, want %v", level, decompressedLevel, tc.data)
				}
			}
		})
	}
}

// TestZipUnzip_Parallel tests zip/unzip functionality under concurrent load
func TestZipUnzip_Parallel(t *testing.T) {
	var (
		count      atomic.Int64
		errorCount atomic.Int64
		wg         sync.WaitGroup
	)

	n := 1000
	wg.Add(n)

	for i := 0; i < n; i++ {
		go func(i int) {
			defer wg.Done()

			data := []byte(strconv.Itoa(i) + "-" + time.Now().String())

			// Zip
			compressed, err := Zip(data)
			if err != nil {
				errorCount.Add(1)
				return
			}

			// Unzip
			decompressed, err := Unzip(compressed)
			if err != nil {
				errorCount.Add(1)
				return
			}

			if bytes.Equal(data, decompressed) {
				count.Add(1)
			} else {
				errorCount.Add(1)
			}
		}(i)
	}

	wg.Wait()

	if errorCount.Load() > 0 {
		t.Errorf("Parallel test failed with %d errors", errorCount.Load())
	}

	if count.Load() != int64(n) {
		t.Errorf("Parallel test: count = %d, want %d", count.Load(), n)
	}
}

// BenchmarkZip benchmarks Zip function performance
func BenchmarkZip(b *testing.B) {
	testData := [][]byte{
		[]byte("small data"),
		bytes.Repeat([]byte("medium data"), 100),
		bytes.Repeat([]byte("large data"), 1000),
	}

	for _, data := range testData {
		b.Run(fmt.Sprintf("Size%d", len(data)), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := Zip(data)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// BenchmarkUnzip benchmarks Unzip function performance
func BenchmarkUnzip(b *testing.B) {
	testData := [][]byte{
		[]byte("small data"),
		bytes.Repeat([]byte("medium data"), 100),
		bytes.Repeat([]byte("large data"), 1000),
	}

	for _, data := range testData {
		compressed, err := Zip(data)
		if err != nil {
			b.Fatal(err)
		}

		b.Run(fmt.Sprintf("Size%d", len(data)), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := Unzip(compressed)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// BenchmarkZip_Parallel benchmarks Zip function under concurrent load
func BenchmarkZip_Parallel(b *testing.B) {
	testData := [][]byte{
		[]byte("small data"),
		bytes.Repeat([]byte("medium data"), 100),
		bytes.Repeat([]byte("large data"), 1000),
	}

	for _, data := range testData {
		b.Run(fmt.Sprintf("Size%d", len(data)), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					_, err := Zip(data)
					if err != nil {
						b.Fatal(err)
					}
				}
			})
		})
	}
}

// BenchmarkUnzip_Parallel benchmarks Unzip function under concurrent load
func BenchmarkUnzip_Parallel(b *testing.B) {
	testData := [][]byte{
		[]byte("small data"),
		bytes.Repeat([]byte("medium data"), 100),
		bytes.Repeat([]byte("large data"), 1000),
	}

	for _, data := range testData {
		compressed, err := Zip(data)
		if err != nil {
			b.Fatal(err)
		}

		b.Run(fmt.Sprintf("Size%d", len(data)), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					_, err := Unzip(compressed)
					if err != nil {
						b.Fatal(err)
					}
				}
			})
		})
	}
}

// printStack prints the current goroutine stack
func printStack() {
	buf := make([]byte, 4096)
	n := runtime.Stack(buf, true)
	fmt.Printf("\n=== STACK TRACE ===\n%s\n===================\n", buf[:n])
}

// FuzzZipUnzip fuzz tests zip/unzip functionality for robustness
func FuzzZipUnzip(f *testing.F) {
	// Seed corpus with common inputs
	f.Add([]byte(""))
	f.Add([]byte("hello"))
	f.Add([]byte{0x00, 0x01, 0x02})
	f.Add([]byte("你好世界"))

	// Run fuzz test
	f.Fuzz(func(t *testing.T, data []byte) {
		// 单协程测试
		t.Run("SingleGoroutine", func(t *testing.T) {
			// Zip
			compressed, err := Zip(data)
			if err != nil {
				fmt.Printf("SingleGoroutine Zip failed: %v\n", err)
				printStack()
				t.Fatalf("Zip failed: %v", err)
			}

			// Unzip
			decompressed, err := Unzip(compressed)
			if err != nil {
				fmt.Printf("SingleGoroutine Unzip failed: %v\n", err)
				printStack()
				t.Fatalf("Unzip failed: %v", err)
			}

			// Verify data integrity
			if !bytes.Equal(data, decompressed) {
				fmt.Printf("SingleGoroutine Data mismatch: input=%v, output=%v\n", data, decompressed)
				printStack()
				t.Fatalf("Data mismatch: input=%v, output=%v", data, decompressed)
			}
		})

		// 多协程并发测试
		t.Run("ConcurrentGoroutines", func(t *testing.T) {
			const goroutines = 10
			var wg sync.WaitGroup
			wg.Add(goroutines)

			for i := 0; i < goroutines; i++ {
				go func(id int) {
					defer wg.Done()

					// Zip
					compressed, err := Zip(data)
					if err != nil {
						fmt.Printf("Concurrent Goroutine %d Zip failed: %v\n", id, err)
						printStack()
						t.Fatalf("Concurrent Zip failed: %v", err)
					}

					// Unzip
					decompressed, err := Unzip(compressed)
					if err != nil {
						fmt.Printf("Concurrent Goroutine %d Unzip failed: %v\n", id, err)
						printStack()
						t.Fatalf("Concurrent Unzip failed: %v", err)
					}

					// Verify data integrity
					if !bytes.Equal(data, decompressed) {
						fmt.Printf("Concurrent Goroutine %d Data mismatch: input=%v, output=%v\n", id, data, decompressed)
						printStack()
						t.Fatalf("Concurrent data mismatch: input=%v, output=%v", data, decompressed)
					}
				}(i)
			}

			wg.Wait()
		})

		// 混合并发测试：同时进行压缩和解压
		t.Run("MixedConcurrentOperations", func(t *testing.T) {
			const operations = 20
			var wg sync.WaitGroup
			wg.Add(operations)

			for i := 0; i < operations; i++ {
				go func(id int) {
					defer wg.Done()

					// 随机选择操作类型
					if rand.Intn(2) == 0 {
						// 压缩操作
						_, err := Zip(data)
						if err != nil {
							fmt.Printf("Mixed Goroutine %d Zip failed: %v\n", id, err)
							printStack()
							t.Fatalf("Mixed Zip failed: %v", err)
						}
					} else {
						// 先压缩再解压
						compressed, err := Zip(data)
						if err != nil {
							fmt.Printf("Mixed Goroutine %d Zip failed: %v\n", id, err)
							printStack()
							t.Fatalf("Mixed Zip failed: %v", err)
						}

						_, err = Unzip(compressed)
						if err != nil {
							fmt.Printf("Mixed Goroutine %d Unzip failed: %v\n", id, err)
							printStack()
							t.Fatalf("Mixed Unzip failed: %v", err)
						}
					}
				}(i)
			}

			wg.Wait()
		})
	})
}

// TestZipLevel_AllLevels tests ZipLevel with all compression levels
func TestZipLevel_AllLevels(t *testing.T) {
	data := []byte("test data for compression levels")

	for level := 1; level <= 12; level++ {
		t.Run(fmt.Sprintf("Level%d", level), func(t *testing.T) {
			compressed, err := ZipLevel(data, level)
			if err != nil {
				t.Errorf("ZipLevel() level %d error = %v", level, err)
			}

			decompressed, err := Unzip(compressed)
			if err != nil {
				t.Errorf("Unzip() for level %d error = %v", level, err)
			}

			if !bytes.Equal(data, decompressed) {
				t.Errorf("Unzip() for level %d = %v, want %v", level, decompressed, data)
			}
		})
	}
}

// TestEdgeCases tests edge cases for zip/unzip
func TestEdgeCases(t *testing.T) {
	// Test very large data
	largeData := bytes.Repeat([]byte("test"), 10000) // 40KB
	compressed, err := Zip(largeData)
	if err != nil {
		t.Errorf("Zip() large data error = %v", err)
	}
	t.Logf("Compressed size: %d", len(compressed))

	decompressed, err := Unzip(compressed)
	if err != nil {
		t.Errorf("Unzip() large data error = %v", err)
	}

	if !bytes.Equal(largeData, decompressed) {
		t.Errorf("Unzip() large data mismatch")
	}

	// Test random data (hard to compress)
	randomData := make([]byte, 1000)
	rand.Read(randomData)
	compressed, err = Zip(randomData)
	if err != nil {
		t.Errorf("Zip() random data error = %v", err)
	}
	t.Logf("Compressed size: %d", len(compressed))

	decompressed, err = Unzip(compressed)
	if err != nil {
		t.Errorf("Unzip() random data error = %v", err)
	}

	if !bytes.Equal(randomData, decompressed) {
		t.Errorf("Unzip() random data mismatch")
	}
}

// TestPoolReuse tests that the sync.Pool is properly reusing objects
func TestPoolReuse(t *testing.T) {
	// Track object reuse by counting how many times New is called
	objectCount := 0
	originalNew := zipPool.New
	zipPool.New = func() interface{} {
		objectCount++
		return originalNew()
	}

	// Reuse objects
	for i := 0; i < 100; i++ {
		c := zipPool.Get().(*Compressor)
		zipPool.Put(c)
	}

	// Restore original New function
	zipPool.New = originalNew

	// Should have created far fewer objects than total operations
	if objectCount >= 100 {
		t.Errorf("Pool not reusing objects: created %d objects for 100 operations", objectCount)
	}

	t.Logf("Pool reuse test: created %d objects for 100 operations", objectCount)
}
