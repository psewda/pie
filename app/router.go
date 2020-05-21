package app

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/psewda/pie/app/models"
	"github.com/psewda/pie/session"
)

const (
	ContentType string = "Content-Type"
	Json        string = "application/json"
)

type DefaultRouter struct {
	Handler http.Handler
	store   session.SessionStore
}

func NewRouter(store session.SessionStore) *DefaultRouter {
	r := mux.NewRouter()
	r = r.PathPrefix("/api/v1/pie-store").Subrouter()

	dr := &DefaultRouter{
		Handler: r,
		store:   store,
	}

	r.HandleFunc("/version", dr.version).Methods("GET")
	r.HandleFunc("/sessions", dr.create).Methods("POST")
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

func (router *DefaultRouter) create(w http.ResponseWriter, r *http.Request) {
	var spec models.SessionSpec
	if err := json.NewDecoder(r.Body).Decode(&spec); err != nil {
		writeBadRequest(w, "Invalid session spec passed. Supply valid session spec.")
		return
	}

	s, err := session.NewSession(spec.Client, spec.Timeout)
	if err != nil {
		writeBadRequest(w, err.Error())
		return
	}
	router.store.Add(s)

	sid := models.SessionId{}
	sid.Id = string(s.GetInfo().Id)
	writeCreated(w, sid)
}

func writeBadRequest(w http.ResponseWriter, value string) {
	err := models.Error{}
	err.Message = value
	write(w, http.StatusBadRequest, err)
}

func writeCreated(w http.ResponseWriter, value interface{}) {
	write(w, http.StatusCreated, value)
}

func writeOk(w http.ResponseWriter, value interface{}) {
	write(w, http.StatusOK, value)
}

func write(w http.ResponseWriter, statusCode int, value interface{}) {
	w.Header().Add(ContentType, Json)
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(value)
}
