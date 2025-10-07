package mgo

import (
	"languages-api/internal/config"
	"languages-api/internal/models"

	"errors"
	"reflect"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func Test_DecodeAll_ShouldReturnErrCursorNilIfCursorIsNil(t *testing.T) {
	_, err := MongoCursor{Cursor: nil}.DecodeAll()
	if !errors.Is(err, models.ErrCursorNil) {
		t.Errorf("Cursor expected to be ErrCursorNil got %v", err)
	}
}

func Test_DecodeAll_ShouldReturnLanguages(t *testing.T) {
	firstAppeared, err := time.Parse(time.RFC3339, "2009-11-10T00:00:00Z")
	if err != nil {
		t.Error("Error parsing timestamp:", err)
	}

	langs := models.Languages{
		Languages: []models.Language{
			models.Language{
				Id:   primitive.NewObjectID(),
				Name: "Golang",
				Creators: []string{
					"Robert Griesemer",
					"Rob Pike",
					"Ken Thompson",
				},
				Extensions: []string{
					".go",
				},
				FirstAppeared: &firstAppeared,
				Year:          2009,
				Wiki:          "https://en.wikipedia.org/wiki/Go_(programming_language)",
			},
		},
	}

	interfaces := make([]interface{}, len(langs.Languages))

	for i, v := range langs.Languages {
		interfaces[i] = v
	}

	mcur, err := mongo.NewCursorFromDocuments(interfaces, nil, nil)
	if err != nil {
		t.Error("Error creating cursor:", err)
	}

	results, err := MongoCursor{Cursor: mcur}.DecodeAll()
	if err != nil {
		t.Error("Error decoding results:", err)

	}

	if !reflect.DeepEqual(results, langs) {
		t.Errorf("DecodeAll returned wrong results: got %v want %v", results, langs)
	}
}

func Test_Ping_ShouldReturnClientPing(t *testing.T) {
	c, err := mongo.NewClient()
	if err != nil {
		t.Error("Error creating client:", err)
	}

	mc := MongoClient{Client: c}
	err = mc.Ping()
	if !errors.Is(err, mongo.ErrClientDisconnected) {
		t.Errorf("Unexpected error pinging client: %v", err)
	}
}

func Test_Disconnect_ShouldReturnClientDisconnect(t *testing.T) {
	c, err := mongo.NewClient()
	if err != nil {
		t.Error("Error creating client:", err)
	}

	mc := MongoClient{Client: c}
	err = mc.Disconnect()
	if !errors.Is(err, mongo.ErrClientDisconnected) {
		t.Errorf("Unexpected error pinging client: %v", err)
	}
}

func Test_Find_ShouldReturnClientFindErrorWithNoFilter(t *testing.T) {
	c, err := mongo.NewClient()
	if err != nil {
		t.Error("Error creating client:", err)
	}

	mc := MongoClient{Client: c, DatabaseName: "test", CollectionName: "test"}
	_, errs := mc.Find(models.Language{})
	if !errors.Is(errs[0], mongo.ErrClientDisconnected) {
		t.Errorf("Unexpected error in Find: %v", errs[0])
	}
}

func Test_Find_ShouldReturnClientFindErrorWithFilter(t *testing.T) {
	c, err := mongo.NewClient()
	if err != nil {
		t.Error("Error creating client:", err)
	}

	mc := MongoClient{Client: c, DatabaseName: "test", CollectionName: "test"}
	firstAppeared, err := time.Parse(time.RFC3339, "2009-11-10T00:00:00Z")
	if err != nil {
		t.Error("Error parsing timestamp:", err)
	}

	_, errs := mc.Find(models.Language{
		Id:   primitive.NewObjectID(),
		Name: "Golang",
		Creators: []string{
			"Robert Griesemer",
			"Rob Pike",
			"Ken Thompson",
		},
		Extensions: []string{
			".go",
		},
		FirstAppeared: &firstAppeared,
		Year:          2009,
		Wiki:          "https://en.wikipedia.org/wiki/Go_(programming_language)",
	})
	if !errors.Is(errs[0], mongo.ErrClientDisconnected) {
		t.Errorf("Unexpected error in Find: %v", errs[0])
	}
}

func Test_FindOne_ShouldReturnErrInvalidIdIfGivenInvalidId(t *testing.T) {
	c, err := mongo.NewClient()
	if err != nil {
		t.Error("Error creating client:", err)
	}

	mc := MongoClient{Client: c, DatabaseName: "test", CollectionName: "test"}

	_, err = mc.FindOne("1")
	if !errors.Is(err, models.ErrInvalidId) {
		t.Errorf("Unexpected error in FindOne: %v", err)
	}
}

func Test_FindOne_ShouldReturnBlankLanguageIfGivenInvalidId(t *testing.T) {
	c, err := mongo.NewClient()
	if err != nil {
		t.Error("Error creating client:", err)
	}

	mc := MongoClient{Client: c, DatabaseName: "test", CollectionName: "test"}

	lang, err := mc.FindOne("1")
	if !errors.Is(err, models.ErrInvalidId) {
		t.Errorf("Unexpected error in FindOne: %v", err)
	}

	if !reflect.DeepEqual(lang, models.Language{}) {
		t.Errorf("FindOne returned wrong results: got %v want %v", lang, models.Language{})
	}
}

func Test_FindOne_ShouldReturnDecodeError(t *testing.T) {
	c, err := mongo.NewClient()
	if err != nil {
		t.Error("Error creating client:", err)
	}

	mc := MongoClient{Client: c, DatabaseName: "test", CollectionName: "test"}

	_, err = mc.FindOne(primitive.NewObjectID().Hex())
	if !errors.Is(err, mongo.ErrClientDisconnected) {
		t.Errorf("Unexpected error in FindOne: %v", err)
	}
}

func Test_InsertOne_ShouldReturnInsertOneError(t *testing.T) {
	c, err := mongo.NewClient()
	if err != nil {
		t.Error("Error creating client:", err)
	}

	mc := MongoClient{Client: c, DatabaseName: "test", CollectionName: "test"}

	_, err = mc.InsertOne(models.Language{})
	if !errors.Is(err, mongo.ErrClientDisconnected) {
		t.Errorf("Unexpected error in InsertOne: %v", err)
	}
}

func Test_ReplaceOne_ShouldReturnErrInvalidIdIfGivenInvalidId(t *testing.T) {
	c, err := mongo.NewClient()
	if err != nil {
		t.Error("Error creating client:", err)
	}

	mc := MongoClient{Client: c, DatabaseName: "test", CollectionName: "test"}

	_, err = mc.ReplaceOne("1", models.Language{})
	if !errors.Is(err, models.ErrInvalidId) {
		t.Errorf("Unexpected error in ReplaceOne: %v", err)
	}
}

func Test_ReplaceOne_ShouldReturnFalseIfGivenInvalidId(t *testing.T) {
	c, err := mongo.NewClient()
	if err != nil {
		t.Error("Error creating client:", err)
	}

	mc := MongoClient{Client: c, DatabaseName: "test", CollectionName: "test"}

	isUpserted, err := mc.ReplaceOne("1", models.Language{})
	if !errors.Is(err, models.ErrInvalidId) {
		t.Errorf("Unexpected error in ReplaceOne: %v", err)
	}

	if isUpserted {
		t.Errorf("FindOne returned wrong results: got %v want %v", isUpserted, false)
	}
}

func Test_ReplaceOne_ShouldReturnReplaceOneError(t *testing.T) {
	c, err := mongo.NewClient()
	if err != nil {
		t.Error("Error creating client:", err)
	}

	mc := MongoClient{Client: c, DatabaseName: "test", CollectionName: "test"}

	_, err = mc.ReplaceOne(primitive.NewObjectID().Hex(), models.Language{})
	if !errors.Is(err, mongo.ErrClientDisconnected) {
		t.Errorf("Unexpected error in ReplaceOne: %v", err)
	}
}

func Test_UpdateOne_ShouldReturnErrInvalidIdIfGivenInvalidId(t *testing.T) {
	c, err := mongo.NewClient()
	if err != nil {
		t.Error("Error creating client:", err)
	}

	mc := MongoClient{Client: c, DatabaseName: "test", CollectionName: "test"}

	err = mc.UpdateOne("1", models.Language{})
	if !errors.Is(err, models.ErrInvalidId) {
		t.Errorf("Unexpected error in UpdateOne: %v", err)
	}
}

func Test_UpdateOne_ShouldReturnUpdateOneError(t *testing.T) {
	c, err := mongo.NewClient()
	if err != nil {
		t.Error("Error creating client:", err)
	}

	mc := MongoClient{Client: c, DatabaseName: "test", CollectionName: "test"}

	err = mc.UpdateOne(primitive.NewObjectID().Hex(), models.Language{})
	if !errors.Is(err, mongo.ErrClientDisconnected) {
		t.Errorf("Unexpected error in UpdateOne: %v", err)
	}
}

func Test_DeleteOne_ShouldReturnErrInvalidIdIfGivenInvalidId(t *testing.T) {
	c, err := mongo.NewClient()
	if err != nil {
		t.Error("Error creating client:", err)
	}

	mc := MongoClient{Client: c, DatabaseName: "test", CollectionName: "test"}

	err = mc.DeleteOne("1")
	if !errors.Is(err, models.ErrInvalidId) {
		t.Errorf("Unexpected error in DeleteOne: %v", err)
	}
}

func Test_DeleteOne_ShouldReturnUpdateOneError(t *testing.T) {
	c, err := mongo.NewClient()
	if err != nil {
		t.Error("Error creating client:", err)
	}

	mc := MongoClient{Client: c, DatabaseName: "test", CollectionName: "test"}

	err = mc.DeleteOne(primitive.NewObjectID().Hex())
	if !errors.Is(err, mongo.ErrClientDisconnected) {
		t.Errorf("Unexpected error in DeleteOne: %v", err)
	}
}

func Test_All_ShouldReturnAllError(t *testing.T) {
	firstAppeared, err := time.Parse(time.RFC3339, "2009-11-10T00:00:00Z")
	if err != nil {
		t.Error("Error parsing timestamp:", err)
	}

	langs := models.Languages{
		Languages: []models.Language{
			models.Language{
				Id:   primitive.NewObjectID(),
				Name: "Golang",
				Creators: []string{
					"Robert Griesemer",
					"Rob Pike",
					"Ken Thompson",
				},
				Extensions: []string{
					".go",
				},
				FirstAppeared: &firstAppeared,
				Year:          2009,
				Wiki:          "https://en.wikipedia.org/wiki/Go_(programming_language)",
			},
		},
	}

	interfaces := make([]interface{}, len(langs.Languages))

	for i, v := range langs.Languages {
		interfaces[i] = v
	}

	mcur, err := mongo.NewCursorFromDocuments(interfaces, nil, nil)
	if err != nil {
		t.Error("Error creating cursor:", err)
	}

	err = MongoCursor{Cursor: mcur}.All(nil, &[]models.Language{})
	if err != nil {
		t.Error("Error calling All:", err)
	}
}

func Test_Close_ShouldReturnCloseError(t *testing.T) {
	firstAppeared, err := time.Parse(time.RFC3339, "2009-11-10T00:00:00Z")
	if err != nil {
		t.Error("Error parsing timestamp:", err)
	}

	langs := models.Languages{
		Languages: []models.Language{
			models.Language{
				Id:   primitive.NewObjectID(),
				Name: "Golang",
				Creators: []string{
					"Robert Griesemer",
					"Rob Pike",
					"Ken Thompson",
				},
				Extensions: []string{
					".go",
				},
				FirstAppeared: &firstAppeared,
				Year:          2009,
				Wiki:          "https://en.wikipedia.org/wiki/Go_(programming_language)",
			},
		},
	}

	interfaces := make([]interface{}, len(langs.Languages))

	for i, v := range langs.Languages {
		interfaces[i] = v
	}

	mcur, err := mongo.NewCursorFromDocuments(interfaces, nil, nil)
	if err != nil {
		t.Error("Error creating cursor:", err)
	}

	err = MongoCursor{Cursor: mcur}.Close(nil)
	if err != nil {
		t.Error("Error calling All:", err)
	}
}

func Test_Decode_ShouldReturnDecodeError(t *testing.T) {
	firstAppeared, err := time.Parse(time.RFC3339, "2009-11-10T00:00:00Z")
	if err != nil {
		t.Error("Error parsing timestamp:", err)
	}

	lang := models.Language{
		Id:   primitive.NewObjectID(),
		Name: "Golang",
		Creators: []string{
			"Robert Griesemer",
			"Rob Pike",
			"Ken Thompson",
		},
		Extensions: []string{
			".go",
		},
		FirstAppeared: &firstAppeared,
		Year:          2009,
		Wiki:          "https://en.wikipedia.org/wiki/Go_(programming_language)",
	}

	sr := mongo.NewSingleResultFromDocument(lang, nil, nil)

	msr := MongoSingleResult{SingleResult: sr}

	err = msr.Decode(&lang)
	if err != nil {
		t.Error("Error calling Decode where none was expected:", err)
	}
}

func Test_Decode_ShouldReturnErrNotFoundOnErrNoDocuments(t *testing.T) {
	firstAppeared, err := time.Parse(time.RFC3339, "2009-11-10T00:00:00Z")
	if err != nil {
		t.Error("Error parsing timestamp:", err)
	}

	lang := models.Language{
		Id:   primitive.NewObjectID(),
		Name: "Golang",
		Creators: []string{
			"Robert Griesemer",
			"Rob Pike",
			"Ken Thompson",
		},
		Extensions: []string{
			".go",
		},
		FirstAppeared: &firstAppeared,
		Year:          2009,
		Wiki:          "https://en.wikipedia.org/wiki/Go_(programming_language)",
	}

	sr := mongo.NewSingleResultFromDocument(lang, mongo.ErrNoDocuments, nil)

	msr := MongoSingleResult{SingleResult: sr}

	err = msr.Decode(&lang)
	if !errors.Is(err, models.ErrNotFound) {
		t.Errorf("Unexpected error calling Decode where ErrNotFound was expected: %v", err)
	}
}

func Test_GetId_ShouldReturnInsertedIdIfResultIsNotNil(t *testing.T) {
	expected := primitive.NewObjectID()

	result := MongoInsertOneResult{InsertOneResult: &mongo.InsertOneResult{InsertedID: expected}}.GetId()
	if result != expected.Hex() {
		t.Errorf("GetId() should return %s, but got %s", expected.Hex(), result)
	}
}

func Test_GetId_ShouldReturnEmptyStringIfResultIsNil(t *testing.T) {
	result := MongoInsertOneResult{InsertOneResult: nil}.GetId()
	if result != "" {
		t.Errorf("GetId() should return an empty string, but got %s", result)
	}
}

func Test_GetIsUpserted_ShouldReturnIfUpsertedCountIsGreaterThanZero(t *testing.T) {
	result := MongoUpdateResult{UpdateResult: &mongo.UpdateResult{UpsertedCount: 1}}.GetIsUpserted()
	if !result {
		t.Errorf("GetIsUpserted() should return true, but got %v", result)
	}
}

func Test_GetIsUpserted_ShouldReturnFalseIfResultIsNil(t *testing.T) {
	result := MongoUpdateResult{UpdateResult: nil}.GetIsUpserted()
	if result {
		t.Errorf("GetIsUpserted() should return false, but got %v", result)
	}
}

func Test_GetUpdateCounts_ShouldReturnCountsIfResultIsNotNil(t *testing.T) {
	modifiedCt, matchedCt := MongoUpdateResult{UpdateResult: &mongo.UpdateResult{ModifiedCount: 1, MatchedCount: 1}}.GetUpdateCounts()
	if modifiedCt != 1 {
		t.Errorf("ModifiedCount should return 1, but got %v", modifiedCt)
	}
	if matchedCt != 1 {
		t.Errorf("MatchedCount should return 1, but got %v", modifiedCt)
	}
}

func Test_GetUpdateCounts_ShouldReturnZeroesIfResultIsNil(t *testing.T) {
	modifiedCt, matchedCt := MongoUpdateResult{UpdateResult: nil}.GetUpdateCounts()
	if modifiedCt != 0 {
		t.Errorf("ModifiedCount should return 0, but got %v", modifiedCt)
	}
	if matchedCt != 0 {
		t.Errorf("MatchedCount should return 0, but got %v", modifiedCt)
	}
}

func Test_GetDeletedCount_ShouldReturnIfUpsertedCountIsGreaterThanZero(t *testing.T) {
	result := MongoDeleteResult{DeleteResult: &mongo.DeleteResult{DeletedCount: 1}}.GetDeletedCount()
	if result != 1 {
		t.Errorf("GetDeletedCount() should return 1, but got %d", result)
	}
}

func Test_GetDeletedCount_ShouldReturnZeroIfResultIsNil(t *testing.T) {
	result := MongoDeleteResult{DeleteResult: nil}.GetDeletedCount()
	if result != 0 {
		t.Errorf("GetDeletedCount() should return 0, but got %d", result)
	}
}

func Test_Connect_ShouldReturnMongoConnectError(t *testing.T) {
	_, err := MongoConnector{}.Connect(config.Config{DBURL: "mongodb://fake"})
	if err != nil {
		t.Error("Unexpected error returned from Connect():", err)
	}
}

func Test_Connect_ShouldReturnMongoClient(t *testing.T) {
	c, err := MongoConnector{}.Connect(config.Config{DBURL: "mongodb://fake"})
	if err != nil {
		t.Error("Unexpected error returned from Connect():", err)
	}

	if reflect.TypeOf(c).String() != "*mgo.MongoClient" {
		t.Errorf("Connect() should return a MongoClient pointer, but got %s", reflect.TypeOf(c).String())
	}
}

func Test_buildMap_ShouldReturnBlankMapIfGivenBlankLanguage(t *testing.T) {
	expected := make(bson.M)

	result := buildMap(models.Language{})

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("buildMap() should return %v, but got %v", expected, result)
	}
}

func Test_buildMap_ShouldReturnFullMapIfGivenFullLanguage(t *testing.T) {
	firstAppeared, err := time.Parse(time.RFC3339, "2009-11-10T00:00:00Z")
	if err != nil {
		t.Error("Error parsing timestamp:", err)
	}

	lang := models.Language{
		Id:   primitive.NewObjectID(),
		Name: "Golang",
		Creators: []string{
			"Robert Griesemer",
			"Rob Pike",
			"Ken Thompson",
		},
		Extensions: []string{
			".go",
		},
		FirstAppeared: &firstAppeared,
		Year:          2009,
		Wiki:          "https://en.wikipedia.org/wiki/Go_(programming_language)",
	}

	expected := make(bson.M)
	expected["name"] = lang.Name
	expected["creators"] = lang.Creators
	expected["extensions"] = lang.Extensions
	expected["firstAppeared"] = lang.FirstAppeared
	expected["year"] = lang.Year
	expected["wiki"] = lang.Wiki

	result := buildMap(lang)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("buildMap() should return %v, but got %v", expected, result)
	}
}
