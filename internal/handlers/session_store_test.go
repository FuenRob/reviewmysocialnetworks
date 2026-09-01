package handlers

import (
	"testing"
	"time"
)

func TestAuthSessionStoreExpiresAndDeletesToken(t *testing.T) {
	store := newAuthSessionStore(2)
	if !store.Store("session", authSession{accessToken: "sensitive-token", expiresAt: time.Now().Add(25 * time.Millisecond)}) {
		t.Fatal("expected session to be stored")
	}

	deadline := time.Now().Add(time.Second)
	for {
		store.mu.Lock()
		remaining := len(store.entries)
		store.mu.Unlock()
		if remaining == 0 {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("expired session was not removed automatically")
		}
		time.Sleep(5 * time.Millisecond)
	}

	if _, ok := store.LoadAndDelete("session"); ok {
		t.Fatal("expired session remained available")
	}
}

func TestAuthSessionStoreIsBoundedAndOneTime(t *testing.T) {
	store := newAuthSessionStore(1)
	expiresAt := time.Now().Add(time.Minute)
	if !store.Store("first", authSession{accessToken: "first-token", expiresAt: expiresAt}) {
		t.Fatal("expected first session to be stored")
	}
	if store.Store("second", authSession{accessToken: "second-token", expiresAt: expiresAt}) {
		t.Fatal("store accepted a session beyond its capacity")
	}

	session, ok := store.LoadAndDelete("first")
	if !ok || session.accessToken != "first-token" {
		t.Fatalf("unexpected stored session: %+v, ok=%v", session, ok)
	}
	if _, ok := store.LoadAndDelete("first"); ok {
		t.Fatal("session was available more than once")
	}
}
