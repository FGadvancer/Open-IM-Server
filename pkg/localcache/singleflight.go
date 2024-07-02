package localcache

import "sync"

type call[V any] struct {
	wg  sync.WaitGroup
	val V
	err error
}

type SingleFlight[K comparable, V any] struct {
	mu sync.Mutex
	m  map[K]*call[V]
}

func NewSingleFlight[K comparable, V any]() *SingleFlight[K, V] {
	return &SingleFlight[K, V]{m: make(map[K]*call[V])}
}

func (r *SingleFlight[K, V]) Do(key K, fn func() (V, error)) (V, error) {
	r.mu.Lock()
	if r.m == nil {
		r.m = make(map[K]*call[V])
	}
	if c, ok := r.m[key]; ok {
		r.mu.Unlock()
		c.wg.Wait()
		return c.val, c.err
	}
	c := new(call[V])
	c.wg.Add(1)
	r.m[key] = c
	r.mu.Unlock()

	c.val, c.err = fn()
	c.wg.Done()

	r.mu.Lock()
	delete(r.m, key)
	r.mu.Unlock()

	return c.val, c.err
}
