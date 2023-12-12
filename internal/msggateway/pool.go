package msggateway

import (
	"sync"
)

// Pool is a generic sync.Pool
type Pool[T any] struct {
	pool sync.Pool
}

// NewPool creates a new pool with a constructor function.
func NewPool[T any](constructor func() T) *Pool[T] {
	return &Pool[T]{
		pool: sync.Pool{
			New: func() any {
				return constructor()
			},
		},
	}
}

// Get retrieves an item from the pool, clearing it if it's not new.
func (p *Pool[T]) Get() T {
	item := p.pool.Get().(T)

	// Assuming T is a pointer to a struct, and the struct has a Reset method.
	// Otherwise, you'll need to clear the data in a way appropriate for T.
	if reseter, ok := any(item).(interface{ Reset() }); ok {
		reseter.Reset()
	}

	return item
}

// Put returns an item to the pool.
func (p *Pool[T]) Put(item T) {
	p.pool.Put(item)
}

// Example usage with a hypothetical type and Reset method.
type MyStruct struct {
	// fields
}

func (m *MyStruct) Reset() {
	// clear all fields
}

func main() {
	pool := NewPool(func() *MyStruct {
		return &MyStruct{}
	})

	item := pool.Get()
	// Use the item
	item.Reset() // Manually reset if needed before putting it back

	pool.Put(item)
}
