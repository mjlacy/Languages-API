package models

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	// ErrNotFound indicates that a language that matches the given criteria was not found
	ErrNotFound = errors.New("language not found")
	// ErrInvalidId indicates an invalid id was sent to the application
	ErrInvalidId = errors.New("invalid id provided")
	// ErrCursorNil indicates that the given cursor is nil, so no functions can be called off it
	ErrCursorNil = errors.New("cursor is nil")
)

type Languages struct {
	Languages []Language `json:"languages" bson:"languages"`
}

type Language struct {
	Id            primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Name          string             `json:"name" bson:"name"`
	Creators      []string           `json:"creators" bson:"creators"`
	Extensions    []string           `json:"extensions" bson:"extensions"`
	FirstAppeared *time.Time         `json:"firstAppeared" bson:"firstAppeared"`
	Year          int32              `json:"year" bson:"year"`
	Wiki          string             `json:"wiki" bson:"wiki"`
}
