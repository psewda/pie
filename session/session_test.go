package session

import (
	"testing"
	"time"

	"github.com/psewda/pie/utils/strutil"
	"github.com/stretchr/testify/assert"
)

func TestNewSession(t *testing.T) {
	// test happy path
	s, err := NewSession(client, time.Second)
	assert.Nil(t, err)
	assert.NotZero(t, s)
	if info := s.GetInfo(); true {
		assert.Equal(t, 20, len(info.Id))
		assert.Equal(t, client, info.Client)
		assert.Equal(t, time.Second, info.Timeout)
	}

	// test invalid inputs
	data := []struct {
		client string
		timout time.Duration
	}{
		{
			client: strutil.Empty,
			timout: time.Second,
		},
		{
			client: client,
			timout: time.Microsecond * 10,
		},
	}
	for _, item := range data {
		s, err := NewSession(item.client, item.timout)
		assert.Error(t, err)
		assert.Zero(t, s)
	}
}

func TestSetItemInSession(t *testing.T) {
	// test happy path
	s, _ := NewSession(client, time.Second)
	err := s.SetItem("k1", []byte("v1"))
	assert.Nil(t, err)

	// test invalid inputs
	if err := s.SetItem(strutil.Empty, []byte("v1")); true {
		assert.Error(t, err)
	}
}

func TestGetItemFromSession(t *testing.T) {
	// test happy path
	s, _ := NewSession(client, time.Second)
	s.SetItem("k1", []byte("v1"))
	v, ok := s.GetItem("k1")
	assert.True(t, ok)
	assert.Equal(t, []byte("v1"), v)

	// test invalid inputs
	data := []string{strutil.Empty, "not-found"}
	for _, item := range data {
		v, ok := s.GetItem(item)
		assert.False(t, ok)
		assert.Nil(t, v)
	}
}

func TestKeySanitize(t *testing.T) {
	s, _ := NewSession(client, time.Second)
	s.SetItem("  KEY  ", []byte("v1"))
	v, ok := s.GetItem("key")
	assert.True(t, ok)
	assert.Equal(t, []byte("v1"), v)
}
