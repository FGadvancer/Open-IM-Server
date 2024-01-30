package localcache

import (
	"context"
	"time"
)

func defaultOption() *option {
	return &option{
		localSlotNum:    500,
		localSlotSize:   20000,
		linkSlotNum:     500,
		expirationEvict: false,
		localSuccessTTL: time.Minute,
		localFailedTTL:  time.Second * 5,
		delFn:           make([]func(ctx context.Context, key ...string), 0, 2),
		hook:            emptyHook{},
	}
}

type option struct {
	localSlotNum  int
	localSlotSize int
	linkSlotNum   int
	// expirationEvict: true means that the cache will be actively cleared when the timer expires,
	// false means that the cache will be lazily deleted.
	expirationEvict bool
	localSuccessTTL time.Duration
	localFailedTTL  time.Duration
	delFn           []func(ctx context.Context, key ...string)
	hook            Hook
}

type Option func(o *option)

func WithExpirationEvict() Option {
	return func(o *option) {
		o.expirationEvict = true
	}
}

func WithLazy() Option {
	return func(o *option) {
		o.expirationEvict = false
	}
}

func WithLocalDisable() Option {
	return WithLinkSlotNum(0)
}

func WithLinkDisable() Option {
	return WithLinkSlotNum(0)
}

func WithLinkSlotNum(linkSlotNum int) Option {
	return func(o *option) {
		o.linkSlotNum = linkSlotNum
	}
}

func WithLocalSlotNum(localSlotNum int) Option {
	return func(o *option) {
		o.localSlotNum = localSlotNum
	}
}

func WithLocalSlotSize(localSlotSize int) Option {
	return func(o *option) {
		o.localSlotSize = localSlotSize
	}
}

func WithLocalSuccessTTL(localSuccessTTL time.Duration) Option {
	if localSuccessTTL < 0 {
		panic("localSuccessTTL should be greater than 0")
	}
	return func(o *option) {
		o.localSuccessTTL = localSuccessTTL
	}
}

func WithLocalFailedTTL(localFailedTTL time.Duration) Option {
	if localFailedTTL < 0 {
		panic("localFailedTTL should be greater than 0")
	}
	return func(o *option) {
		o.localFailedTTL = localFailedTTL
	}
}

func WithHook(hook Hook) Option {
	if hook == nil {
		panic("hook should not be nil")
	}
	return func(o *option) {
		o.hook = hook
	}
}

func WithDeleteKeyBefore(fn func(ctx context.Context, key ...string)) Option {
	if fn == nil {
		panic("fn should not be nil")
	}
	return func(o *option) {
		o.delFn = append(o.delFn, fn)
	}
}

type emptyHook struct{}

func (e emptyHook) IncrementHit() {}

func (e emptyHook) IncrementMiss() {}

func (e emptyHook) IncrementDelHit() {}

func (e emptyHook) IncrementDelMiss() {}
