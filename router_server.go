//go:build !js

package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/fbundle/http_transport/http_transport"
	"github.com/fbundle/sudoku/sudoku"
)

func setup() (http_transport.Router, sudoku.Store, func()) {
	mux := http.NewServeMux()
	group := http_transport.NewGo(mux).Group("sudoku")
	mux.Handle("/sudoku/", http.StripPrefix("/sudoku/", http.FileServer(http.Dir("docs/"))))
	store := newSession(60 * time.Second)
	return group, store, func() {
		fmt.Println("Server is up at: http://0.0.0.0:3000/sudoku/")
		http.ListenAndServe("0.0.0.0:3000", mux)
	}
}

// session is a TTL key-value store for multiplayer game sessions.
type entry struct {
	lastAccess time.Time
	data       any
}

type session struct {
	timeout time.Duration
	pool    map[string]*entry
	mtx     sync.RWMutex
}

func newSession(timeout time.Duration) *session {
	s := &session{timeout: timeout, pool: make(map[string]*entry)}
	go s.cleanLoop()
	return s
}

func (s *session) cleanLoop() {
	for {
		time.Sleep(s.timeout)
		s.mtx.Lock()
		for key, e := range s.pool {
			if time.Since(e.lastAccess) > s.timeout {
				delete(s.pool, key)
			}
		}
		s.mtx.Unlock()
	}
}

func (s *session) Set(key string, data any) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.pool[key] = &entry{lastAccess: time.Now(), data: data}
}

func (s *session) Get(key string) any {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	e, ok := s.pool[key]
	if !ok {
		return nil
	}
	e.lastAccess = time.Now()
	return e.data
}

func (s *session) NumActiveKey() int {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	return len(s.pool)
}
