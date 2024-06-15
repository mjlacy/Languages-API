package mgo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mockCollection struct {
	cursor          *mongo.Cursor
	singleResult    *mongo.SingleResult
	insertOneResult *mongo.InsertOneResult
	updateResult    *mongo.UpdateResult
	deleteResult    *mongo.DeleteResult
	err             error
}

func (mc mockCollection) Find(_ context.Context, _ interface{}, _ ...*options.FindOptions) (cur *mongo.Cursor, err error) {
	return mc.cursor, mc.err
}

func (mc mockCollection) FindOne(_ context.Context, _ interface{}, _ ...*options.FindOneOptions) *mongo.SingleResult {
	return mc.singleResult
}

func (mc mockCollection) InsertOne(_ context.Context, document interface{}, _ ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	return mc.insertOneResult, mc.err
}

func (mc mockCollection) ReplaceOne(_ context.Context, _ interface{}, replacement interface{}, _ ...*options.ReplaceOptions) (*mongo.UpdateResult, error) {
	return mc.updateResult, mc.err
}

func (mc mockCollection) UpdateOne(_ context.Context, _ interface{}, update interface{}, _ ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return mc.updateResult, mc.err
}

func (mc mockCollection) DeleteOne(_ context.Context, _ interface{}, _ ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	return mc.deleteResult, mc.err
}
