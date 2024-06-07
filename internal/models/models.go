package models

import (
	"errors"
	"time"

	// "go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	// ErrNotFound indicates that a language that matches the given criteria was not found
	ErrNotFound = errors.New("language not found")
	// ErrInvalidId indicates an invalid id was sent to the application
	ErrInvalidId = errors.New("invalid id provided")
)

type Languages struct {
	Languages []Language `json:"languages" bson:"languages"`
}

type Language struct {
	Id            string             `json:"_id" bson:"_id,omitempty"` // Changed to string to simplify responses
	Name          string             `json:"name" bson:"name"`
	Creators      []string           `json:"creators" bson:"creators"`
	Extensions    []string           `json:"extensions" bson:"extensions"`
	FirstAppeared *time.Time         `json:"firstAppeared" bson:"firstAppeared"`
	Year          int32              `json:"year" bson:"year"`
	Wiki          string             `json:"wiki" bson:"wiki"`
}

type LanguagesResponse struct {
	Languages []LanguageResponse `json:"languages"`
	Links []Links `json:"links"`
	Total int `json:"total"`
}

type LanguageResponse struct {
	Language
	Links []Links `json:"links"`
}

type Links struct {
	Rel  string `json:"rel"`
	Href string `json:"href"`
}

type QueryString struct {
	Language
	Page *int `json:"page"` // pointer to differentiate between nil and 0
	Size *int `json:"size"` // pointer to differentiate between nil and 0
}
