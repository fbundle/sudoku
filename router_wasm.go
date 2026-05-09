//go:build js && wasm

package main

import (
	"sync"

	"github.com/fbundle/sudoku/sudoku"
	"github.com/fbundle/http_transport/http_transport"
)

func setup() (http_transport.Router, sudoku.Store, func()) {
	return http_transport.New(), &localStore{data: make(map[string]any)}, func() { select {} }
}

// localStore is a simple in-memory store for single-player WASM sessions.
type localStore struct {
	mu   sync.RWMutex
	data map[string]any
}

func (s *localStore) Set(key string, value any) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
}

func (s *localStore) Get(key string) any {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.data[key]
}

func (s *localStore) NumActiveKey() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.data)
}
