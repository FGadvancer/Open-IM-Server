package msggateway

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

var encodeBufferPool = NewPool[*bytes.Buffer](func() *bytes.Buffer { return new(bytes.Buffer) },
	func(b *bytes.Buffer) {
		b.Reset()
	})

func TestGobEncoder_EncodeDecode(t *testing.T) {

	encoder := NewGobEncoder()
	for i := 0; i < 2000; i++ {
		src := map[string]string{"a": "b"}

		// compress
		des, err := encoder.Encode(src)
		assert.Equal(t, nil, err)

		// decompress
		decode := make(map[string]string)
		err = encoder.Decode(des, &decode)
		assert.Equal(t, nil, err)

		// check
		assert.EqualValues(t, src, decode)
	}
}

func TestGobEncoder_EncodeDecodeWithPool(t *testing.T) {

	encoder := NewGobEncoder()
	for i := 0; i < 2000; i++ {
		src := map[string]string{"a": "b"}

		buf := encodeBufferPool.Get()
		// compress
		err := encoder.EncodeWithExternalPool(src, buf)
		assert.Equal(t, nil, err)

		// decompress
		dest := make(map[string]string)
		err = encoder.DecodeWithExternalPool(buf, &dest)
		encodeBufferPool.Put(buf)
		assert.Equal(t, nil, err)

		// check
		assert.EqualValues(t, src, dest)
	}
}

func BenchmarkEncode(b *testing.B) {
	src := map[string]string{"a": "b"}
	encoder := NewGobEncoder()

	for i := 0; i < b.N; i++ {
		_, err := encoder.Encode(src)
		assert.Equal(b, nil, err)
	}
}

// pkg: github.com/openimsdk/open-im-server/v3/internal/msggateway
// cpu: Intel(R) Core(TM) i5-10400 CPU @ 2.90GHz
// BenchmarkEncode-4                         546922              2004 ns/op            1008 B/op         16 allocs/op
// BenchmarkEncodeWithExternalPool-4         655075              1949 ns/op             896 B/op         14 allocs/op
func BenchmarkEncodeWithExternalPool(b *testing.B) {
	src := map[string]string{"a": "b"}
	encoder := NewGobEncoder()

	for i := 0; i < b.N; i++ {
		buf := encodeBufferPool.Get()
		err := encoder.EncodeWithExternalPool(src, buf)
		assert.Equal(b, nil, err)
		encodeBufferPool.Put(buf)
	}
}
