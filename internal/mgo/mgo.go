package mgo

//
//import (
//	"languages-api/internal/config"
//	"languages-api/internal/models"
//
//	"context"
//	"errors"
//	"time"
//
//	"github.com/rs/zerolog/log"
//	"go.mongodb.org/mongo-driver/bson"
//	"go.mongodb.org/mongo-driver/bson/primitive"
//	"go.mongodb.org/mongo-driver/mongo"
//	"go.mongodb.org/mongo-driver/mongo/options"
//	"go.mongodb.org/mongo-driver/mongo/readpref"
//)
//
//const (
//	FiveSeconds = 5 * time.Second
//	TenSeconds  = 10 * time.Second
//)
//
//// Client is for wrappers of mongo.Client
//type Client interface {
//	Ping() error
//	Disconnect() error
//	Find(filter interface{}) (languages models.Languages, errors []error)
//	FindOne(id string) (language models.Language, err error)
//	InsertOne(document interface{}) (insertedId string, err error)
//	ReplaceOne(id string, document interface{}) (isUpserted bool, err error)
//	UpdateOne(id string, update interface{}) (err error)
//	DeleteOne(id string) (err error)
//}
//
//// MongoClient implements the Client interface
//type MongoClient struct {
//	*mongo.Client
//	DatabaseName   string
//	CollectionName string
//}
//
//type MongoDatabase struct {
//	*mongo.Database
//}
//
//func (mc MongoCursor) DecodeAll() (languages models.Languages, err error) {
//	if mc.Cursor == nil {
//		return languages, models.ErrCursorNil
//	}
//
//	ctx, cancel := context.WithTimeout(context.Background(), FiveSeconds)
//	defer cancel()
//
//	defer func() {
//		err := mc.Cursor.Close(ctx)
//		if err != nil {
//			log.Error().Err(err).Msg("Failed to close database cursor")
//		}
//	}()
//
//	err = mc.Cursor.All(ctx, &languages.Languages)
//
//	return
//}
//
//type MongoCollection struct {
//	*mongo.Collection
//}
//
//type MongoCursor struct {
//	*mongo.Cursor
//}
//
//type MongoSingleResult struct {
//	*mongo.SingleResult
//}
//
//type MongoInsertOneResult struct {
//	*mongo.InsertOneResult
//}
//
//type MongoUpdateResult struct {
//	*mongo.UpdateResult
//}
//
//type MongoDeleteResult struct {
//	*mongo.DeleteResult
//}
//
//// Ping checks the connection to mongo
//func (mc MongoClient) Ping() error {
//	ctx, cancel := context.WithTimeout(context.Background(), TenSeconds)
//	defer cancel()
//
//	return mc.Client.Ping(ctx, readpref.Primary())
//}
//
//// Disconnect terminates the connection to mongo
//func (mc MongoClient) Disconnect() error {
//	ctx, cancel := context.WithTimeout(context.Background(), TenSeconds)
//	defer cancel()
//
//	return mc.Client.Disconnect(ctx)
//}
//
//func (mc MongoClient) Find(filter interface{}) (languages models.Languages, errs []error) {
//	conditions := bson.M{}
//
//	language := filter.(models.Language)
//
//	if language.Name != "" {
//		conditions["name"] = bson.M{"$eq": language.Name}
//	}
//
//	if len(language.Creators) > 0 {
//		conditions["creators"] = bson.M{"$all": language.Creators}
//	}
//
//	if len(language.Extensions) > 0 {
//		conditions["extensions"] = bson.M{"$all": language.Extensions}
//	}
//
//	if language.FirstAppeared != nil {
//		conditions["firstAppeared"] = bson.M{"$eq": language.FirstAppeared}
//	}
//
//	if language.Year != 0 {
//		conditions["year"] = bson.M{"$eq": language.Year}
//	}
//
//	if language.Wiki != "" {
//		conditions["wiki"] = bson.M{"$eq": language.Wiki}
//	}
//
//	ctx, cancel := context.WithTimeout(context.Background(), FiveSeconds)
//	defer cancel()
//
//	cursor, err := mc.Client.Database(mc.DatabaseName).Collection(mc.CollectionName).Find(ctx, conditions)
//	if err != nil {
//		errs = append(errs, err)
//	}
//
//	languages, err = MongoCursor{Cursor: cursor}.DecodeAll()
//	if err != nil {
//		errs = append(errs, err)
//	}
//
//	if len(languages.Languages) == 0 {
//		languages.Languages = []models.Language{}
//	}
//
//	return
//}
//
//func (mc MongoClient) FindOne(id string) (language models.Language, err error) {
//	objectId, err := primitive.ObjectIDFromHex(id)
//	if err != nil {
//		return models.Language{}, models.ErrInvalidId
//	}
//
//	ctx, cancel := context.WithTimeout(context.Background(), FiveSeconds)
//	defer cancel()
//
//	err = MongoSingleResult{SingleResult: mc.Client.Database(mc.DatabaseName).Collection(mc.CollectionName).FindOne(ctx, bson.M{"_id": objectId})}.Decode(&language)
//
//	return
//}
//
//func (mc MongoClient) InsertOne(document interface{}) (insertedId string, err error) {
//	ctx, cancel := context.WithTimeout(context.Background(), FiveSeconds)
//	defer cancel()
//
//	ior, err := mc.Client.Database(mc.DatabaseName).Collection(mc.CollectionName).InsertOne(ctx, document)
//
//	insertedId = MongoInsertOneResult{InsertOneResult: ior}.GetId()
//
//	return
//}
//
//func (mc MongoClient) ReplaceOne(id string, document interface{}) (isUpserted bool, err error) {
//	objectId, err := primitive.ObjectIDFromHex(id)
//	if err != nil {
//		return false, models.ErrInvalidId
//	}
//
//	ctx, cancel := context.WithTimeout(context.Background(), FiveSeconds)
//	defer cancel()
//
//	upsert := options.ReplaceOptions{}
//	ur, err := mc.Client.Database(mc.DatabaseName).Collection(mc.CollectionName).ReplaceOne(ctx, bson.M{"_id": objectId}, document, upsert.SetUpsert(true))
//
//	isUpserted = MongoUpdateResult{UpdateResult: ur}.GetIsUpserted()
//
//	return
//}
//
//func (mc MongoClient) UpdateOne(id string, update interface{}) (err error) {
//	objectId, err := primitive.ObjectIDFromHex(id)
//	if err != nil {
//		return models.ErrInvalidId
//	}
//
//	ctx, cancel := context.WithTimeout(context.Background(), FiveSeconds)
//	defer cancel()
//
//	lang := update.(models.Language)
//
//	ur, err := mc.Client.Database(mc.DatabaseName).Collection(mc.CollectionName).UpdateOne(ctx, bson.M{"_id": objectId}, bson.M{"$set": buildMap(lang)})
//
//	modifiedCount, matchedCount := MongoUpdateResult{UpdateResult: ur}.GetUpdateCounts()
//
//	if err == nil && modifiedCount == 0 && matchedCount == 0 {
//		err = models.ErrNotFound
//	}
//
//	return
//}
//
//func (mc MongoClient) DeleteOne(id string) (err error) {
//	objectId, err := primitive.ObjectIDFromHex(id)
//	if err != nil {
//		return models.ErrInvalidId
//	}
//
//	ctx, cancel := context.WithTimeout(context.Background(), FiveSeconds)
//	defer cancel()
//
//	dr, err := mc.Client.Database(mc.DatabaseName).Collection(mc.CollectionName).DeleteOne(ctx, bson.M{"_id": objectId})
//
//	deletedCount := MongoDeleteResult{DeleteResult: dr}.GetDeletedCount()
//	if err == nil && deletedCount == 0 {
//		err = models.ErrNotFound
//	}
//
//	return
//}
//
//func (mc MongoCursor) All(ctx context.Context, results interface{}) error {
//	return mc.Cursor.All(ctx, results)
//}
//
//func (mc MongoCursor) Close(ctx context.Context) error {
//	return mc.Cursor.Close(ctx)
//}
//
//func (msr MongoSingleResult) Decode(result interface{}) (err error) {
//	err = msr.SingleResult.Decode(result)
//	if errors.Is(err, mongo.ErrNoDocuments) {
//		err = models.ErrNotFound
//	}
//	return err
//}
//
//func (mior MongoInsertOneResult) GetId() string {
//	if mior.InsertOneResult != nil {
//		return mior.InsertedID.(primitive.ObjectID).Hex()
//	} else {
//		return ""
//	}
//}
//
//func (mur MongoUpdateResult) GetIsUpserted() bool {
//	if mur.UpdateResult != nil {
//		return mur.UpsertedCount > 0
//	} else {
//		return false
//	}
//}
//
//func (mur MongoUpdateResult) GetUpdateCounts() (int64, int64) {
//	if mur.UpdateResult != nil {
//		return mur.ModifiedCount, mur.MatchedCount
//	} else {
//		return 0, 0
//	}
//}
//
//func (mdr MongoDeleteResult) GetDeletedCount() int64 {
//	if mdr.DeleteResult != nil {
//		return mdr.DeletedCount
//	} else {
//		return 0
//	}
//}
//
//// Connector specifies the methods needed to connect to mongo
//type Connector interface {
//	Connect(cfg config.Config) (Client, error)
//}
//
//// MongoConnector implements the Connector interface
//type MongoConnector struct{}
//
//// Connect establishes the connection to mongo
//func (mc MongoConnector) Connect(cfg config.Config) (Client, error) {
//	ctx, cancel := context.WithTimeout(context.Background(), TenSeconds)
//	defer cancel()
//	opts := options.Client().ApplyURI(cfg.DBURL)
//
//	client, err := mongo.Connect(ctx, opts)
//	return &MongoClient{Client: client, DatabaseName: cfg.Database, CollectionName: cfg.Collection}, err
//}
//
//func buildMap(language models.Language) bson.M {
//	update := make(bson.M)
//
//	if language.Name != "" {
//		update["name"] = language.Name
//	}
//
//	if len(language.Creators) > 0 {
//		update["creators"] = language.Creators
//	}
//
//	if len(language.Extensions) > 0 {
//		update["extensions"] = language.Extensions
//	}
//
//	if language.FirstAppeared != nil {
//		update["firstAppeared"] = language.FirstAppeared
//	}
//
//	if language.Year != 0 {
//		update["year"] = language.Year
//	}
//
//	if language.Wiki != "" {
//		update["wiki"] = language.Wiki
//	}
//
//	return update
//}
