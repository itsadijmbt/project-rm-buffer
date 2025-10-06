package main

import (
	"context"
	"path/filepath"
	"sync"
)

// limit total concurrent jobs
type WorkerPool struct{ tokens chan struct{} }

func NewWorkerPool(n int) *WorkerPool {
	if n < 1 {
		n = 1
	}
	return &WorkerPool{tokens: make(chan struct{}, n)}
}

func (p *WorkerPool) Do(ctx context.Context, fn func() error) error {
	select {
	case p.tokens <- struct{}{}:
		defer func() { <-p.tokens }()
	case <-ctx.Done():
		return ctx.Err()
	}
	return fn()
}

// per-path mutex map
type PathLocker struct {
	mu   sync.Mutex
	pool map[string]*sync.Mutex
}

func NewPathLocker() *PathLocker { return &PathLocker{pool: make(map[string]*sync.Mutex)} }

func (pl *PathLocker) lockFor(key string) *sync.Mutex {
	key = filepath.Clean(key)
	pl.mu.Lock()
	m, ok := pl.pool[key]
	if !ok {
		m = &sync.Mutex{}
		pl.pool[key] = m
	}
	pl.mu.Unlock()
	return m
}

func (pl *PathLocker) With(key string, fn func() error) error {
	m := pl.lockFor(key)
	m.Lock()
	defer m.Unlock()
	return fn()
}
