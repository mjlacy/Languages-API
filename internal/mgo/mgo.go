package mgo

import (
	"languages-api/internal/config"
	"languages-api/internal/models"

	"context"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	FiveSeconds = 5 * time.Second
	TenSeconds  = 10 * time.Second
)

type Repository interface {
	Close() error
	GetLanguages(lanugage models.QueryString) (lanugages models.Languages, total int, filteredTotal int, err error)
	GetLanguage(id string) (lanugage *models.Language, err error)
	PostLanguage(lanugage *models.Language) (id string, err error)
	PutLanguage(id string, lanugage *models.Language) (isCreated bool, err error)
	PatchLanguage(id string, update models.Language) (err error)
	DeleteLanguage(id string) (err error)
	Ping() error
}

type Repo struct {
	client *mongo.Client
	config config.Config
}

func New(cfg config.Config) (r *Repo, err error) {
	r = &Repo{
		config: cfg,
	}

	err = r.InitClient()

	return
}

func (r *Repo) InitClient() (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), TenSeconds)
	defer cancel()

	opts := options.Client().ApplyURI(r.config.DBURL)

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create mongo client")
		return
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Error().Err(err).Msg("Failed to ping mongo")
		return
	}

	r.client = client
	return
}

func (r *Repo) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), FiveSeconds)
	defer cancel()
	return r.client.Ping(ctx, readpref.Primary())
}

func (r *Repo) GetLanguages(queryString models.QueryString) (languages models.Languages, total int, filteredTotal int, err error) {
	conditions := bson.M{}
	opts := options.Find()

	if queryString.Name != "" {
		conditions["name"] = bson.M{"$eq": queryString.Name}
	}

	if len(queryString.Creators) > 0 {
		conditions["creators"] = bson.M{"$all": queryString.Creators}
	}

	if len(queryString.Extensions) > 0 {
		conditions["extensions"] = bson.M{"$all": queryString.Extensions}
	}

	if queryString.FirstAppeared != nil {
		conditions["firstAppeared"] = bson.M{"$eq": queryString.FirstAppeared}
	}

	if queryString.Year != 0 {
		conditions["year"] = bson.M{"$eq": queryString.Year}
	}

	if queryString.Wiki != "" {
		conditions["wiki"] = bson.M{"$eq": queryString.Wiki}
	}

	ctx, cancel := context.WithTimeout(context.Background(), FiveSeconds)
	defer cancel()

	collection := r.client.Database(r.config.Database).Collection(r.config.Collection)

	languages.Languages = []models.Language{}

	count, err := collection.CountDocuments(ctx, bson.D{})
	if err != nil {
		log.Error().Err(err).Msg("Failed to count languages")
		return
	} else {
		total = int(count)
	}

	if len(conditions) > 0 {
		filteredCount := int64(0)
		filteredCount, err = collection.CountDocuments(ctx, conditions)
		if err != nil {
			log.Error().Err(err).Msg("Failed to count filtered languages")
			return
		} else {
			filteredTotal = int(filteredCount)
		}
	} else {
		filteredTotal = 0
	}

	if *queryString.Size != -1 {
		opts.SetLimit(int64(*queryString.Size))
	} else {
		opts.SetLimit(count)
	}

	opts.SetSkip(*opts.Limit * int64(*queryString.Page - 1))

	if queryString.SortBy != "" {
		sortBySlice := strings.Split(queryString.SortBy, ",")

		sortDoc := bson.D{}
		for i := 0; i < len(sortBySlice); i++ {
			elements := strings.Split(sortBySlice[i], " ")
			if elements[1] == "asc" {
				sortDoc = append(sortDoc, bson.E{Key: elements[0], Value: 1})
			} else if elements[1] == "desc" {
				sortDoc = append(sortDoc, bson.E{Key: elements[0], Value: -1})
			}
		}

		opts.SetSort(sortDoc)
	}

	cur, err := collection.Find(ctx, conditions, opts)
	if err != nil {
		log.Error().Err(err).Msg("Failed to find languages")
		return
	}

	defer cur.Close(ctx)

	err = cur.All(ctx, &languages.Languages)

	return
}

func (r *Repo) GetLanguage(id string) (Language *models.Language, err error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, models.ErrInvalidId
	}

	ctx, cancel := context.WithTimeout(context.Background(), FiveSeconds)
	defer cancel()

	collection := r.client.Database(r.config.Database).Collection(r.config.Collection)

	err = collection.FindOne(ctx, bson.M{"_id": objectId}).Decode(&Language)
	if err == mongo.ErrNoDocuments {
		err = models.ErrNotFound
	}

	return
}

func (r *Repo) PostLanguage(Language *models.Language) (id string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), FiveSeconds)
	defer cancel()

	collection := r.client.Database(r.config.Database).Collection(r.config.Collection)

	result, err := collection.InsertOne(ctx, Language)
	if err != nil {
		return
	}

	return result.InsertedID.(primitive.ObjectID).Hex(), err
}

func (r *Repo) PutLanguage(id string, Language *models.Language) (isCreated bool, err error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, models.ErrInvalidId
	}

	upsert := options.ReplaceOptions{}
	ctx, cancel := context.WithTimeout(context.Background(), FiveSeconds)
	defer cancel()

	collection := r.client.Database(r.config.Database).Collection(r.config.Collection)

	result, err := collection.ReplaceOne(ctx, bson.M{"_id": objectId}, Language, upsert.SetUpsert(true))
	if err != nil {
		return false, err
	}

	return result.UpsertedCount > 0, err
}

func (r *Repo) PatchLanguage(id string, update models.Language) (err error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.ErrInvalidId
	}

	ctx, cancel := context.WithTimeout(context.Background(), FiveSeconds)
	defer cancel()

	collection := r.client.Database(r.config.Database).Collection(r.config.Collection)

	result, err := collection.UpdateOne(ctx, bson.M{"_id": objectId}, bson.M{"$set": buildMap(update)})
	if err == nil && result.ModifiedCount == 0 && result.MatchedCount == 0 {
		err = models.ErrNotFound
	}
	return
}

func (r *Repo) DeleteLanguage(id string) (err error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.ErrInvalidId
	}

	ctx, cancel := context.WithTimeout(context.Background(), FiveSeconds)
	defer cancel()

	collection := r.client.Database(r.config.Database).Collection(r.config.Collection)

	result, err := collection.DeleteOne(ctx, bson.M{"_id": objectId})
	if err == nil && result.DeletedCount == 0 {
		err = models.ErrNotFound
	}

	return
}

func buildMap(language models.Language) bson.M {
	update := make(bson.M)

	if language.Name != "" {
		update["name"] = language.Name
	}

	if len(language.Creators) > 0 {
		update["creators"] = language.Creators
	}

	if len(language.Extensions) > 0 {
		update["extensions"] = language.Extensions
	}

	if language.FirstAppeared != nil {
		update["firstAppeared"] = language.FirstAppeared
	}

	if language.Year != 0 {
		update["year"] = language.Year
	}

	if language.Wiki != "" {
		update["wiki"] = language.Wiki
	}

	return update
}

func (r *Repo) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), FiveSeconds)
	defer cancel()

	return r.client.Disconnect(ctx)
}
