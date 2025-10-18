package repo

import (
	"languages-api/internal/config"
	"languages-api/internal/mariadb"
	"strings"

	"testing"

	"database/sql"
)

func Test_New_ShouldReturnConnectError(t *testing.T) {
	_, err := New(config.Config{DBURL: "root:password@tcp(host.docker.internal:3306)/languages?parseTime=true"}, mariadb.MariaConnector{})
	if !strings.Contains(err.Error(), "connect: connection refused") {
		t.Errorf("New() returned an unexpected error")
	}
}

func Test_Ping_ShouldReturnDisconnectError(t *testing.T) {
	db, err := sql.Open("mysql", "root:password@tcp(host.docker.internal:3306)/languages?parseTime=true")
	if err != nil {
		t.Error("Error creating db:", err)
	}

	repo := &Repo{Client: mariadb.MariaClient{DB: db}}

	err = repo.Close()
	if err != nil {
		t.Error("Error closing repo:", err)
	}

	err = repo.Ping()
	if err.Error() != "sql: database is closed" {
		t.Errorf("Ping() returned an unexpected error")
	}
}

//func Test_GetLanguages_ShouldReturnFindError(t *testing.T) {
//	c, err := mongo.NewClient()
//	if err != nil {
//		t.Error("Error creating client:", err)
//	}
//
//	_, errs := (&Repo{client: mgo.MongoClient{Client: c, DatabaseName: "test", CollectionName: "test"}}).GetLanguages(models.Language{})
//	if !errors.Is(errs[0], mongo.ErrClientDisconnected) {
//		t.Errorf("GetLanguages() returned an unexpected error: %v", errs[0])
//	}
//}
//
//func Test_GetLanguage_ShouldReturnFindOneError(t *testing.T) {
//	c, err := mongo.NewClient()
//	if err != nil {
//		t.Error("Error creating client:", err)
//	}
//
//	_, err = (&Repo{client: mgo.MongoClient{Client: c, DatabaseName: "test", CollectionName: "test"}}).GetLanguage(primitive.NewObjectID().Hex())
//	if !errors.Is(err, mongo.ErrClientDisconnected) {
//		t.Errorf("GetLanguage() returned an unexpected error: %v", err)
//	}
//}
//
//func Test_PostLanguage_ShouldReturnInsertOneError(t *testing.T) {
//	c, err := mongo.NewClient()
//	if err != nil {
//		t.Error("Error creating client:", err)
//	}
//
//	_, err = (&Repo{client: mgo.MongoClient{Client: c, DatabaseName: "test", CollectionName: "test"}}).PostLanguage(models.Language{})
//	if !errors.Is(err, mongo.ErrClientDisconnected) {
//		t.Errorf("PostLanguage() returned an unexpected error: %v", err)
//	}
//}
//
//func Test_PutLanguage_ShouldReturnReplaceOneError(t *testing.T) {
//	c, err := mongo.NewClient()
//	if err != nil {
//		t.Error("Error creating client:", err)
//	}
//
//	_, err = (&Repo{client: mgo.MongoClient{Client: c, DatabaseName: "test", CollectionName: "test"}}).PutLanguage(primitive.NewObjectID().Hex(), models.Language{})
//	if !errors.Is(err, mongo.ErrClientDisconnected) {
//		t.Errorf("PutLanguage() returned an unexpected error: %v", err)
//	}
//}
//
//func Test_PatchLanguage_ShouldReturnUpdateOneError(t *testing.T) {
//	c, err := mongo.NewClient()
//	if err != nil {
//		t.Error("Error creating client:", err)
//	}
//
//	err = (&Repo{client: mgo.MongoClient{Client: c, DatabaseName: "test", CollectionName: "test"}}).PatchLanguage(primitive.NewObjectID().Hex(), models.Language{})
//	if !errors.Is(err, mongo.ErrClientDisconnected) {
//		t.Errorf("PatchLanguage() returned an unexpected error: %v", err)
//	}
//}
//
//func Test_DeleteLanguage_ShouldReturnUpdateOneError(t *testing.T) {
//	c, err := mongo.NewClient()
//	if err != nil {
//		t.Error("Error creating client:", err)
//	}
//
//	err = (&Repo{client: mgo.MongoClient{Client: c, DatabaseName: "test", CollectionName: "test"}}).DeleteLanguage(primitive.NewObjectID().Hex())
//	if !errors.Is(err, mongo.ErrClientDisconnected) {
//		t.Errorf("DeleteLanguage() returned an unexpected error: %v", err)
//	}
//}
