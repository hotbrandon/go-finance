package models

import "time"

type JobStatus struct {
	ID   int       `json:"id"`
	Name string    `json:"name,omitempty"`
	Spec string    `json:"spec,omitempty"`
	Next time.Time `json:"next,omitempty"`
	Prev time.Time `json:"prev,omitempty"`
}
