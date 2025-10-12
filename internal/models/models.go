package models

import (
	"errors"
	"time"
)

var (
	// ErrNotFound indicates that a language that matches the given criteria was not found
	ErrNotFound = errors.New("language not found")
	// ErrInvalidId indicates an invalid id was sent to the application
	ErrInvalidId = errors.New("invalid id provided")
	// ErrInvalidRequestBody indicates that the given request body is invalid
	ErrInvalidRequestBody = errors.New("invalid request body")
)

type Languages struct {
	Languages []Language `json:"languages" bson:"languages"`
}

type Language struct {
	Id            int        `json:"_id,omitempty"`
	Name          string     `json:"name"`
	Creators      []string   `json:"creators"`
	Extensions    []string   `json:"extensions"`
	FirstAppeared *time.Time `json:"firstAppeared"`
	Year          int32      `json:"year"`
	Wiki          string     `json:"wiki"`
}
