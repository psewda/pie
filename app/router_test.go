package app

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/psewda/pie/app/models"
	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
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
