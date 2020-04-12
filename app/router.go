package app

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/psewda/pie/app/models"
)

type DefaultRouter struct {
	Handler http.Handler
}

func NewRouter() *DefaultRouter {
	r := mux.NewRouter()
	r = r.PathPrefix("/api/v1/pie-store").Subrouter()

	dr := &DefaultRouter{
		Handler: r,
	}

	r.HandleFunc("/version", dr.version).Methods("GET")
	return dr
}

func (dr *DefaultRouter) version(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	v := models.Version{
		Version:   Version,
		Golang:    Golang,
		GitCommit: GitCommit,
		Built:     Built,
		OsArch:    OsArch,
	}
	json.NewEncoder(w).Encode(v)
}
