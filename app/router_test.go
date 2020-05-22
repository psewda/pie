package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/psewda/pie/app/models"
	"github.com/psewda/pie/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockSessionStore struct {
	mock.Mock
}

func (m *mockSessionStore) Add(s session.Session) error {
	args := m.Called(s)
	return args.Error(0)
}

func TestEndpointVersion(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/v1/pie-store/version", nil)

	dr := DefaultRouter{}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(dr.version)
	handler.ServeHTTP(rr, req)

	v := func() models.Version {
		assert.Equal(t, http.StatusOK, rr.Code)
		var ver models.Version
		json.Unmarshal(rr.Body.Bytes(), &ver)
		return ver
	}()

	assert.Equal(t, Version, v.Version)
	assert.Equal(t, Golang, v.Golang)
	assert.Equal(t, GitCommit, v.GitCommit)
	assert.Equal(t, Built, v.Built)
	assert.Equal(t, OsArch, v.OsArch)
}

func TestEndpointCreateSession(t *testing.T) {
	url := "/api/v1/pie-store/sessions"
	store := session.NewSessionStore()
	r := NewRouter(store)

	// test happy path
	body := toBody("client", time.Second*2)
	req, _ := http.NewRequest("POST", url, body)
	rr := httptest.NewRecorder()
	r.Handler.ServeHTTP(rr, req)

	var sid models.SessionId
	assert.Equal(t, http.StatusCreated, rr.Code)
	json.Unmarshal(rr.Body.Bytes(), &sid)
	assert.Equal(t, 20, len(sid.Id))

	// test invalid inputs
	data := []string{"invalid", `{"client":"abc"}`}
	for _, item := range data {
		req, _ := http.NewRequest("POST", url, strings.NewReader(item))
		rr := httptest.NewRecorder()
		r.Handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	}
}

func TestEndpointSetValue(t *testing.T) {
	url := "/api/v1/pie-store/sessions/%s/items/i"
	s, _ := session.NewSession("client", time.Second*30)
	sid := s.GetInfo().Id
	store := session.NewSessionStore()
	store.Add(s)
	r := NewRouter(store)

	// test happy path
	body := bytes.NewReader([]byte("value"))
	req, _ := http.NewRequest("PUT", fmt.Sprintf(url, sid), body)
	rr := httptest.NewRecorder()
	r.Handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, string(rr.Body.Bytes()), `"Size":5`)

	// test invalid inputs
	data := []struct {
		sid    string
		value  []byte
		status int
	}{
		{
			sid:    s.GetInfo().Id,
			value:  []byte{},
			status: http.StatusBadRequest,
		},
		{
			sid:    "not-found",
			value:  []byte("value"),
			status: http.StatusNotFound,
		},
	}
	for _, item := range data {

		body := bytes.NewReader(item.value)
		req, _ := http.NewRequest("PUT", fmt.Sprintf(url, item.sid), body)
		rr := httptest.NewRecorder()
		r.Handler.ServeHTTP(rr, req)
		assert.Equal(t, item.status, rr.Code)
	}
}

func TestEndpointGetValue(t *testing.T) {
	url := "/api/v1/pie-store/sessions/%s/items/%s"
	s, _ := session.NewSession("client", time.Second*30)
	s.SetItem("k", []byte("value"))
	sid := s.GetInfo().Id
	store := session.NewSessionStore()
	store.Add(s)
	r := NewRouter(store)

	// test happy path
	req, _ := http.NewRequest("GET", fmt.Sprintf(url, sid, "k"), nil)
	rr := httptest.NewRecorder()
	r.Handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "value", string(rr.Body.Bytes()))

	// test invalid inputs
	data := []struct {
		sid string
		key string
	}{
		{
			sid: "invalid",
			key: "k",
		},
		{
			sid: s.GetInfo().Id,
			key: "invalid",
		},
	}
	for _, i := range data {
		req, _ := http.NewRequest("GET", fmt.Sprintf(url, i.sid, i.key), nil)
		rr := httptest.NewRecorder()
		r.Handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusNotFound, rr.Code)
	}
}

func toBody(client string, timeout time.Duration) io.Reader {
	spec := models.SessionSpec{
		Client:  client,
		Timeout: timeout,
	}
	json, _ := json.Marshal(spec)
	return bytes.NewReader(json)
}
