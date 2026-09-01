package handlers

import (
	"sync"
	"time"
)

const maxAuthSessions = 10_000

type storedAuthSession struct {
	value authSession
	timer *time.Timer
}

type authSessionStore struct {
	mu       sync.Mutex
	entries  map[string]*storedAuthSession
	capacity int
}

func newAuthSessionStore(capacity int) *authSessionStore {
	return &authSessionStore{
		entries:  make(map[string]*storedAuthSession),
		capacity: capacity,
	}
}

func (s *authSessionStore) Store(id string, session authSession) bool {
	now := time.Now()
	if id == "" || !session.expiresAt.After(now) {
		return false
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.deleteExpiredLocked(now)
	if len(s.entries) >= s.capacity {
		return false
	}

	entry := &storedAuthSession{value: session}
	s.entries[id] = entry
	entry.timer = time.AfterFunc(time.Until(session.expiresAt), func() {
		s.deleteIfCurrent(id, entry)
	})
	return true
}

func (s *authSessionStore) LoadAndDelete(id string) (authSession, bool) {
	s.mu.Lock()
	entry, ok := s.entries[id]
	if ok {
		delete(s.entries, id)
		if entry.timer != nil {
			entry.timer.Stop()
		}
	}
	s.mu.Unlock()

	if !ok || !entry.value.expiresAt.After(time.Now()) {
		return authSession{}, false
	}
	return entry.value, true
}

func (s *authSessionStore) deleteIfCurrent(id string, expected *storedAuthSession) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if current, ok := s.entries[id]; ok && current == expected {
		delete(s.entries, id)
	}
}

func (s *authSessionStore) deleteExpiredLocked(now time.Time) {
	for id, entry := range s.entries {
		if !entry.value.expiresAt.After(now) {
			delete(s.entries, id)
			if entry.timer != nil {
				entry.timer.Stop()
			}
		}
	}
}
