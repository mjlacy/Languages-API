package repo

import (
	"languages-api/internal/config"
	"languages-api/internal/mgo"
	"languages-api/internal/models"

	"github.com/rs/zerolog/log"
)

type Repository interface {
	Close() error
	Ping() error
	GetLanguages(language models.Language) (languages models.Languages, errors []error)
	GetLanguage(id string) (language models.Language, err error)
	PostLanguage(language models.Language) (insertedId string, err error)
	PutLanguage(id string, language models.Language) (isUpserted bool, err error)
	PatchLanguage(id string, update models.Language) (err error)
	DeleteLanguage(id string) (err error)
}

type Repo struct {
	client mgo.Client
}

func New(cfg config.Config, c mgo.Connector) (r *Repo, err error) {
	r = &Repo{}
	r.client, err = c.Connect(cfg)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create database client")
		return
	}

	err = r.client.Ping()
	if err != nil {
		log.Error().Err(err).Msg("Failed to ping database")
		return
	}

	return
}

func (r *Repo) Close() error {
	return r.client.Disconnect()
}

func (r *Repo) Ping() error {
	return r.client.Ping()
}

func (r *Repo) GetLanguages(language models.Language) (languages models.Languages, errors []error) {
	return r.client.Find(language)
}

func (r *Repo) GetLanguage(id string) (language models.Language, err error) {
	return r.client.FindOne(id)
}

func (r *Repo) PostLanguage(language models.Language) (insertedId string, err error) {
	return r.client.InsertOne(language)
}

func (r *Repo) PutLanguage(id string, language models.Language) (isUpserted bool, err error) {
	return r.client.ReplaceOne(id, language)
}

func (r *Repo) PatchLanguage(id string, update models.Language) (err error) {
	return r.client.UpdateOne(id, update)
}

func (r *Repo) DeleteLanguage(id string) (err error) {
	return r.client.DeleteOne(id)
}
