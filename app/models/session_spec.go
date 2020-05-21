package models

import "time"

type SessionSpec struct {
	Client  string        `json:"client"`
	Timeout time.Duration `json:"timeout"`
}
