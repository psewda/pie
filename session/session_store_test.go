package session

import (
	"testing"
	"time"

	"github.com/psewda/pie/utils/strutil"
	"github.com/stretchr/testify/assert"
)

func TestAddSessionInStore(t *testing.T) {
	// test happy path
	s, _ := NewSession(client, time.Second)
	store := NewSessionStore()
	defer store.Dispose()
	err := store.Add(s)
	assert.Nil(t, err)

	// test invalid inputs
	if err = store.Add(Session{}); true {
		assert.Error(t, err)
	}
}

func TestGetSessionFromStore(t *testing.T) {
	// test happy path
	s, _ := NewSession(client, time.Second)
	store := NewSessionStore()
	defer store.Dispose()
	store.Add(s)
	if s, ok := store.Get(s.GetInfo().Id); true {
		assert.True(t, ok)
		assert.NotZero(t, s)
	}

	// test invalid inputs
	data := []string{strutil.Empty, "non-found"}
	for _, item := range data {
		s, ok := store.Get(uid(item))
		assert.False(t, ok)
		assert.Zero(t, s)
	}
}

func TestRemoveSessionFromStore(t *testing.T) {
	s, _ := NewSession(client, time.Second)
	store := NewSessionStore()
	defer store.Dispose()
	store.Add(s)

	data := []struct {
		id       uid
		expected bool
	}{
		{
			id:       s.GetInfo().Id,
			expected: true,
		},
		{
			id:       uid(strutil.Empty),
			expected: false,
		},
		{
			id:       "not-found",
			expected: false,
		},
	}

	for _, item := range data {
		ok := store.Remove(item.id)
		assert.Equal(t, item.expected, ok)
	}
}

func TestRemoveSessionAfterTimeout(t *testing.T) {
	s, _ := NewSession(client, time.Second)
	store := NewSessionStore()
	defer store.Dispose()
	store.Add(s)

	time.Sleep(time.Millisecond * 1200) // 1.2 second delay
	if s, ok := store.Get(s.GetInfo().Id); true {
		assert.False(t, ok)
		assert.Zero(t, s)
	}
}
