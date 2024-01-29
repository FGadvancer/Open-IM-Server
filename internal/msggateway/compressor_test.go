// Copyright © 2023 OpenIM. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package msggateway

import (
	"bytes"
	"crypto/rand"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testBufferPool = NewPool[*bytes.Buffer](func() *bytes.Buffer { return new(bytes.Buffer) },
	func(b *bytes.Buffer) {
		b.Reset()
	})

func mockRandom() []byte {
	bs := make([]byte, 50)
	rand.Read(bs)
	return bs
}

func TestCompressDecompress(t *testing.T) {

	compressor := NewGzipCompressor()

	for i := 0; i < 2000; i++ {
		src := mockRandom()

		// compress
		dest, err := compressor.CompressWithPool(src)
		assert.Equal(t, nil, err)

		// decompress
		res, err := compressor.DecompressWithPool(dest)
		assert.Equal(t, nil, err)

		// check
		assert.EqualValues(t, src, res)
	}
}

func TestCompressDecompressWithConcurrency(t *testing.T) {
	wg := sync.WaitGroup{}
	compressor := NewGzipCompressor()

	for i := 0; i < 200; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			src := mockRandom()

			// compress
			des, err := compressor.CompressWithPool(src)
			assert.Equal(t, nil, err)

			// decompress
			res, err := compressor.DecompressWithPool(des)
			assert.Equal(t, nil, err)

			// check
			assert.EqualValues(t, src, res)

		}()
	}
	wg.Wait()
}

func BenchmarkCompress(b *testing.B) {
	src := mockRandom()
	compressor := NewGzipCompressor()

	for i := 0; i < b.N; i++ {
		_, err := compressor.Compress(src)
		assert.Equal(b, nil, err)
	}
}

func BenchmarkCompressWithSyncPool(b *testing.B) {
	src := mockRandom()

	compressor := NewGzipCompressor()
	for i := 0; i < b.N; i++ {
		_, err := compressor.CompressWithPool(src)
		assert.Equal(b, nil, err)
	}
}

func BenchmarkCompressWithExternalPool(b *testing.B) {
	src := mockRandom()
	compressor := NewGzipCompressor()

	for i := 0; i < b.N; i++ {
		buf := testBufferPool.Get()
		err := compressor.CompressWithExternalPool(src, buf)
		assert.Equal(b, nil, err)
		testBufferPool.Put(buf)
	}
}

func BenchmarkDecompress(b *testing.B) {
	src := mockRandom()

	compressor := NewGzipCompressor()
	comdata, err := compressor.Compress(src)
	assert.Equal(b, nil, err)

	for i := 0; i < b.N; i++ {
		_, err := compressor.DeCompress(comdata)
		assert.Equal(b, nil, err)
	}
}

func BenchmarkDecompressWithSyncPool(b *testing.B) {
	src := mockRandom()

	compressor := NewGzipCompressor()
	comdata, err := compressor.Compress(src)
	assert.Equal(b, nil, err)

	for i := 0; i < b.N; i++ {
		_, err := compressor.DecompressWithPool(comdata)
		assert.Equal(b, nil, err)
	}
}

// goos: windows
// goarch: amd64
// pkg: github.com/openimsdk/open-im-server/v3/internal/msggateway
// cpu: Intel(R) Core(TM) i5-10400 CPU @ 2.90GHz
// BenchmarkCompress-4                                 6620            172501 ns/op          814112 B/op         20 allocs/op
// BenchmarkCompressWithSyncPool-4                    58910             20262 ns/op             253 B/op          3 allocs/op
// BenchmarkCompressWithExternalPool-4                59946             19697 ns/op              13 B/op          0 allocs/op
// BenchmarkDecompress-4                             116978              9703 ns/op           41751 B/op          7 allocs/op
// BenchmarkDecompressWithSyncPool-4                 859456              1426 ns/op             563 B/op          2 allocs/op
// BenchmarkDecompressWithExternalPool-4            1326120               845.5 ns/op            49 B/op          1 allocs/op
func BenchmarkDecompressWithExternalPool(b *testing.B) {
	src := mockRandom()

	compressor := NewGzipCompressor()
	comdata, err := compressor.Compress(src)
	assert.Equal(b, nil, err)

	for i := 0; i < b.N; i++ {
		buf := testBufferPool.Get()
		err := compressor.DecompressWithExternalPool(comdata, buf)
		assert.Equal(b, nil, err)
		testBufferPool.Put(buf)
	}
}
