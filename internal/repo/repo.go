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
	GetLanguages(language models.Language) (languages models.Languages, err error)
	GetLanguage(id string) (language *models.Language, err error)
	PostLanguage(language *models.Language) (insertedId string, err error)
	PutLanguage(id string, language *models.Language) (isUpserted bool, err error)
	PatchLanguage(id string, update models.Language) (err error)
	DeleteLanguage(id string) (err error)
}

type Repo struct {
	client     mgo.Client
	collection string
	database   string
}

func New(cfg config.Config, c mgo.Connector) (r *Repo, err error) {
	r = &Repo{}
	r.client, err = c.Connect(cfg.DBURL)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create database client")
		return
	}

	err = r.client.Ping()
	if err != nil {
		log.Error().Err(err).Msg("Failed to ping database")
		return
	}

	r.database = cfg.Database
	r.collection = cfg.Collection

	return
}

func (r *Repo) Close() error {
	return r.client.Disconnect()
}

func (r *Repo) Ping() error {
	return r.client.Ping()
}

func (r *Repo) GetLanguages(language models.Language) (languages models.Languages, err error) {
	languages.Languages = []models.Language{}

	collection := r.client.Database(r.database).Collection(r.collection)

	cur, err := collection.Find(language)
	if err != nil {
		log.Error().Err(err).Msg("Failed to find languages")
		return
	}

	defer func() {
		err := cur.Close()
		if err != nil {
			log.Error().Err(err).Msg("Failed to close database cursor")
		}
	}()

	err = cur.All(&languages.Languages)

	return
}

func (r *Repo) GetLanguage(id string) (language *models.Language, err error) {
	collection := r.client.Database(r.database).Collection(r.collection)
	err = collection.FindOne(id).Decode(&language)

	return
}

func (r *Repo) PostLanguage(language *models.Language) (insertedId string, err error) {
	collection := r.client.Database(r.database).Collection(r.collection)
	insertedId, err = collection.InsertOne(language)

	return
}

func (r *Repo) PutLanguage(id string, language *models.Language) (isUpserted bool, err error) {
	collection := r.client.Database(r.database).Collection(r.collection)
	isUpserted, err = collection.ReplaceOne(id, language)

	return
}

func (r *Repo) PatchLanguage(id string, update models.Language) (err error) {
	collection := r.client.Database(r.database).Collection(r.collection)
	err = collection.UpdateOne(id, update)

	return
}

func (r *Repo) DeleteLanguage(id string) (err error) {
	collection := r.client.Database(r.database).Collection(r.collection)
	err = collection.DeleteOne(id)

	return
}
