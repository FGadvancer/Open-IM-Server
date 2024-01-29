package msggateway

import (
	"bytes"
	"compress/gzip"
	"sync"
)

var (
	bufferPool = NewPool[*bytes.Buffer](func() *bytes.Buffer { return new(bytes.Buffer) },
		func(b *bytes.Buffer) {
			b.Reset()
		})
	reqPool = NewPool[*Req](func() *Req { return new(Req) },
		func(r *Req) {
			r.Data = nil
			r.MsgIncr = ""
			r.OperationID = ""
			r.ReqIdentifier = 0
			r.SendID = ""
			r.Token = ""
		})
	gzipWriterPool = NewPool[*gzip.Writer](func() *gzip.Writer { return gzip.NewWriter(nil) }, nil)
	gzipReaderPool = NewPool[*gzip.Reader](func() *gzip.Reader { return new(gzip.Reader) }, nil)
)

// Pool is a generic sync.Pool
type Pool[T any] struct {
	pool      sync.Pool
	resetFunc func(T)
}

// NewPool creates a new pool with a constructor and reset function.
func NewPool[T any](constructor func() T, resetFunc func(T)) *Pool[T] {
	return &Pool[T]{
		pool: sync.Pool{
			New: func() any {
				return constructor()
			},
		},
		resetFunc: resetFunc,
	}
}

// Get retrieves an item from the pool, clearing it if it's not new.
func (p *Pool[T]) Get() T {
	item := p.pool.Get().(T)

	// Assuming T has a Reset method and p.resetFunc is not nil exec reset to clear the data before returning it
	if p.resetFunc != nil {
		p.resetFunc(item)
	}

	return item
}

// Put returns an item to the pool.
func (p *Pool[T]) Put(item T) {
	p.pool.Put(item)
}
