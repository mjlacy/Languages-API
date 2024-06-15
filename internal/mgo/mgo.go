package mgo

import (
	"languages-api/internal/models"

	"context"
	"errors"
	"time"

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

// Client is for wrappers of mongo.Client
type Client interface {
	Ping() error
	Disconnect() error
	Database(name string) Database
}

type Database interface {
	Collection(name string) Collection
}

type Cursor interface {
	All(result interface{}) error
	Close() error
}

type Decoder interface {
	Decode(result interface{}) error
}

type Collection interface {
	Find(filter interface{}) (Cursor, error)
	FindOne(id string) Decoder
	InsertOne(document interface{}) (insertedId string, err error)
	ReplaceOne(id string, document interface{}) (isUpserted bool, err error)
	UpdateOne(id string, update interface{}) (err error)
	DeleteOne(id string) (err error)
}

// MongoClient implements the Client interface
type MongoClient struct {
	*mongo.Client
}

type MongoDatabase struct {
	*mongo.Database
}

type MongoCollection struct {
	*mongo.Collection
}

type MongoCursor struct {
	*mongo.Cursor
}

type MongoSingleResult struct {
	*mongo.SingleResult
	error
}

// Ping checks the connection to mongo
func (mc MongoClient) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), TenSeconds)
	defer cancel()

	return mc.Client.Ping(ctx, readpref.Primary())
}

// Disconnect terminates the connection to mongo
func (mc MongoClient) Disconnect() error {
	ctx, cancel := context.WithTimeout(context.Background(), TenSeconds)
	defer cancel()

	return mc.Client.Disconnect(ctx)
}

func (mc MongoClient) Database(name string) Database {
	opts := options.Database()
	return MongoDatabase{Database: mc.Client.Database(name, opts)}
}

func (md MongoDatabase) Collection(name string) Collection {
	opts := options.Collection()
	return MongoCollection{Collection: md.Database.Collection(name, opts)}
}

func (mcoll MongoCollection) Find(filter interface{}) (Cursor, error) {
	conditions := bson.M{}

	language := filter.(models.Language)

	if language.Name != "" {
		conditions["name"] = bson.M{"$eq": language.Name}
	}

	if len(language.Creators) > 0 {
		conditions["creators"] = bson.M{"$all": language.Creators}
	}

	if len(language.Extensions) > 0 {
		conditions["extensions"] = bson.M{"$all": language.Extensions}
	}

	if language.FirstAppeared != nil {
		conditions["firstAppeared"] = bson.M{"$eq": language.FirstAppeared}
	}

	if language.Year != 0 {
		conditions["year"] = bson.M{"$eq": language.Year}
	}

	if language.Wiki != "" {
		conditions["wiki"] = bson.M{"$eq": language.Wiki}
	}

	ctx, cancel := context.WithTimeout(context.Background(), FiveSeconds)
	defer cancel()

	opts := options.FindOptions{}
	cursor, err := mcoll.Collection.Find(ctx, conditions, &opts)

	return MongoCursor{Cursor: cursor}, err
}

func (mcoll MongoCollection) FindOne(id string) Decoder {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return MongoSingleResult{SingleResult: nil, error: models.ErrInvalidId}
	}

	ctx, cancel := context.WithTimeout(context.Background(), FiveSeconds)
	defer cancel()
	opts := options.FindOneOptions{}

	return MongoSingleResult{SingleResult: mcoll.Collection.FindOne(ctx, bson.M{"_id": objectId}, &opts), error: nil}
}

func (mcoll MongoCollection) InsertOne(document interface{}) (insertedId string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), FiveSeconds)
	defer cancel()

	opts := options.InsertOneOptions{}
	res, err := mcoll.Collection.InsertOne(ctx, document, &opts)

	return res.InsertedID.(primitive.ObjectID).Hex(), err
}

func (mcoll MongoCollection) ReplaceOne(id string, document interface{}) (isUpserted bool, err error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, models.ErrInvalidId
	}

	ctx, cancel := context.WithTimeout(context.Background(), FiveSeconds)
	defer cancel()

	upsert := options.ReplaceOptions{}
	res, err := mcoll.Collection.ReplaceOne(ctx, bson.M{"_id": objectId}, document, upsert.SetUpsert(true))

	return res.UpsertedCount > 0, err
}

func (mcoll MongoCollection) UpdateOne(id string, update interface{}) (err error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.ErrInvalidId
	}

	ctx, cancel := context.WithTimeout(context.Background(), FiveSeconds)
	defer cancel()

	lang := update.(models.Language)

	res, err := mcoll.Collection.UpdateOne(ctx, bson.M{"_id": objectId}, bson.M{"$set": buildMap(lang)})
	if err == nil && res.ModifiedCount == 0 && res.MatchedCount == 0 {
		err = models.ErrNotFound
	}

	return
}

func (mcoll MongoCollection) DeleteOne(id string) (err error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.ErrInvalidId
	}

	ctx, cancel := context.WithTimeout(context.Background(), FiveSeconds)
	defer cancel()

	res, err := mcoll.Collection.DeleteOne(ctx, bson.M{"_id": objectId})
	if err == nil && res.DeletedCount == 0 {
		err = models.ErrNotFound
	}

	return
}

func (mcur MongoCursor) All(result interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), FiveSeconds)
	defer cancel()

	return mcur.Cursor.All(ctx, result)
}

func (mcur MongoCursor) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), FiveSeconds)
	defer cancel()

	return mcur.Cursor.Close(ctx)
}

func (msr MongoSingleResult) Decode(result interface{}) error {
	if msr.error != nil {
		return msr.error
	} else {
		err := msr.SingleResult.Decode(result)
		if errors.Is(err, mongo.ErrNoDocuments) {
			err = models.ErrNotFound
		}
		return err
	}
}

// Connector specifies the methods needed to connect to mongo
type Connector interface {
	Connect(DBURL string) (Client, error)
}

// MongoConnector implements the Connector interface
type MongoConnector struct{}

// Connect establishes the connection to mongo
func (mc MongoConnector) Connect(DBURL string) (Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), TenSeconds)
	defer cancel()
	opts := options.Client().ApplyURI(DBURL)

	client, err := mongo.Connect(ctx, opts)
	return &MongoClient{Client: client}, err
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
