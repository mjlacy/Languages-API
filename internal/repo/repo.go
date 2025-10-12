package repo

import (
	"languages-api/internal/config"
	"languages-api/internal/mariadb"
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
	Client mariadb.Client
}

func New(cfg config.Config, c mariadb.Connector) (r *Repo, err error) {
	r = &Repo{}
	r.Client, err = c.Connect(cfg)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create database client")
		return
	}

	err = r.Client.Ping()
	if err != nil {
		log.Error().Err(err).Msgf("Error while pinging Maria: %v", err)
	}

	return
}

func (r *Repo) Close() error {
	return r.Client.Disconnect()
}

func (r *Repo) Ping() error {
	return r.Client.Ping()
}

func (r *Repo) GetLanguages(language models.Language) (languages models.Languages, errors []error) {
	return r.Client.Find(language)
}

func (r *Repo) GetLanguage(id string) (language models.Language, err error) {
	return r.Client.FindOne(id)
}

func (r *Repo) PostLanguage(language models.Language) (insertedId string, err error) {
	return r.Client.InsertOne(language)
}

func (r *Repo) PutLanguage(id string, language models.Language) (isUpserted bool, err error) {
	return r.Client.ReplaceOne(id, language)
}

func (r *Repo) PatchLanguage(id string, update models.Language) (err error) {
	return r.Client.UpdateOne(id, update)
}

func (r *Repo) DeleteLanguage(id string) (err error) {
	return r.Client.DeleteOne(id)
}
