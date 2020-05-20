package session

import (
	"errors"
	"log"
	"time"
)

type SessionStore interface {
	Add(s Session) error
	Get(id uid) (Session, bool)
	Remove(id uid) bool
	Dispose()
}

type internalSession struct {
	session Session
	dieAt   time.Time
}

type inMemoryStore struct {
	sessions map[uid]internalSession
	quit     chan struct{}
}

func (store *inMemoryStore) Add(s Session) error {
	if !s.ready {
		err := errors.New("Session is not in ready state, please pass valid session.")
		return err
	}

	info := s.GetInfo()
	is := internalSession{
		session: s,
		dieAt:   time.Now().Add(info.Timeout),
	}
	store.sessions[info.Id] = is
	return nil
}

func (store *inMemoryStore) Get(id uid) (Session, bool) {
	if len(id) == 0 {
		return Session{}, false
	}

	is, ok := store.sessions[id]
	if ok {
		info := is.session.GetInfo()
		is.dieAt = time.Now().Add(info.Timeout)
	}
	return is.session, ok
}

func (store *inMemoryStore) Remove(id uid) bool {
	if len(id) > 0 {
		if _, ok := store.sessions[id]; ok {
			delete(store.sessions, id)
			log.Printf("session removed from memory store => [sid: %s]", id)
			return true
		}
	}
	return false
}

func (store *inMemoryStore) Dispose() {
	close(store.quit)
}

func NewSessionStore() SessionStore {
	store := &inMemoryStore{
		sessions: make(map[uid]internalSession),
		quit:     make(chan struct{}),
	}
	go runKiller(store)
	return store
}

func runKiller(store *inMemoryStore) {
forloop:
	for {
		select {
		case <-time.After(time.Second):
			for k, s := range store.sessions {
				if time.Now().After(s.dieAt) {
					store.Remove(k)
				}
			}
		case <-store.quit:
			break forloop
		}
	}
}
