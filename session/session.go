package session

import (
	"errors"
	"strings"
	"time"

	"github.com/rs/xid"
)

type SessionInfo struct {
	Id      string
	Client  string
	Timeout time.Duration
}

type Session struct {
	info  SessionInfo
	items map[string][]byte
	ready bool
}

func (s *Session) GetInfo() SessionInfo {
	return SessionInfo{
		Id:      s.info.Id,
		Client:  s.info.Client,
		Timeout: s.info.Timeout,
	}
}

func (s *Session) SetItem(key string, value []byte) error {
	if len(key) == 0 {
		err := errors.New("Key is empty, pass valid key.")
		return err
	}
	s.items[sanitize(key)] = value
	return nil
}

func (s *Session) GetItem(key string) ([]byte, bool) {
	if len(key) == 0 {
		return nil, false
	}
	value, ok := s.items[sanitize(key)]
	return value, ok
}

func NewSession(client string, timeout time.Duration) (Session, error) {
	if len(client) == 0 {
		err := errors.New("Invalid client value passed, pass valid value.")
		return Session{}, err
	} else if timeout < time.Second {
		err := errors.New("Timeout must be equal to or greater than 1 second.")
		return Session{}, err
	}

	s := Session{
		info: SessionInfo{
			Id:      xid.New().String(),
			Client:  client,
			Timeout: timeout,
		},
		items: make(map[string][]byte),
		ready: true,
	}
	return s, nil
}

func sanitize(value string) string {
	if len(value) > 0 {
		return strings.TrimSpace(strings.ToLower(value))
	}
	return value
}
