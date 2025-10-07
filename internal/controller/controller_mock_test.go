package controller

import (
	"languages-api/internal/models"

	"errors"
	"net/http"
	"reflect"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Test_Header_ShouldReturnStructHeader(t *testing.T) {
	header := http.Header{}

	mwr := mockResponseWriter{header: header}

	result := mwr.Header()

	if !reflect.DeepEqual(mwr.header, result) {
		t.Errorf("Header should return %v, but got %v", mwr.header, result)
	}
}

func Test_Write_ShouldSetStructMessage(t *testing.T) {
	message := "Golang"

	mwr := mockResponseWriter{}

	_, err := mwr.Write([]byte(message))
	if err != nil {
		t.Errorf("Write should not return error, but got %v", err)
	}

	if !reflect.DeepEqual(mwr.message, message) {
		t.Errorf("Header should return %v, but got %v", mwr.message, message)
	}
}

func Test_Write_ShouldReturnStructNum(t *testing.T) {
	mwr := mockResponseWriter{num: 42}

	num, err := mwr.Write([]byte(""))
	if err != nil {
		t.Errorf("Write should not return error, but got %v", err)
	}

	if !reflect.DeepEqual(mwr.num, num) {
		t.Errorf("Header should return %v, but got %v", mwr.num, num)
	}
}

func Test_Write_ShouldReturnStructError(t *testing.T) {
	e := errors.New("golang")
	mwr := mockResponseWriter{err: e}

	_, err := mwr.Write([]byte(""))
	if !errors.Is(err, e) {
		t.Errorf("Write should return struct error, but got %v", err)
	}
}

func Test_WriteHeader_ShouldSetStructStatusCode(t *testing.T) {
	mwr := mockResponseWriter{}

	mwr.WriteHeader(http.StatusTeapot)
	if !reflect.DeepEqual(mwr.statusCode, http.StatusTeapot) {
		t.Errorf("WriteHeader should return %v, but got %v", mwr.statusCode, http.StatusTeapot)
	}
}

func Test_Ping_ShouldReturnStructError(t *testing.T) {
	e := errors.New("golang")
	mr := mockRepository{err: e}

	err := mr.Ping()
	if !errors.Is(err, e) {
		t.Errorf("Ping should return struct error, but got %v", err)
	}
}

func Test_Close_ShouldReturnStructError(t *testing.T) {
	e := errors.New("golang")
	mr := mockRepository{err: e}

	err := mr.Close()
	if !errors.Is(err, e) {
		t.Errorf("Close should return struct error, but got %v", err)
	}
}

func Test_GetLanguages_ShouldReturnStructLanguages(t *testing.T) {
	firstAppeared, err := time.Parse(time.RFC3339, "2009-11-10T00:00:00Z")
	if err != nil {
		t.Error("Error parsing timestamp:", err)
	}

	expected := models.Languages{
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

	mr := mockRepository{ls: expected}

	langs, errs := mr.GetLanguages(models.Language{})
	if errs != nil {
		t.Errorf("GetLanguages should not return error, but got %v", errs)
	}

	if !reflect.DeepEqual(langs, expected) {
		t.Errorf("GetLanguages should return %v, but got %v", expected, langs)
	}
}

func Test_GetLanguages_ShouldReturnStructErrors(t *testing.T) {
	expected := []error{
		errors.New("golang"),
	}

	mr := mockRepository{errs: expected}

	_, errs := mr.GetLanguages(models.Language{})
	if !reflect.DeepEqual(errs, expected) {
		t.Errorf("GetLanguages should return %v, but got %v", expected, errs)
	}
}

func Test_GetLanguage_ShouldReturnStructLanguage(t *testing.T) {
	firstAppeared, err := time.Parse(time.RFC3339, "2009-11-10T00:00:00Z")
	if err != nil {
		t.Error("Error parsing timestamp:", err)
	}

	expected := models.Language{
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

	mr := mockRepository{l: expected}

	lang, err := mr.GetLanguage("")
	if err != nil {
		t.Errorf("GetLanguage should not return error, but got %v", err)
	}

	if !reflect.DeepEqual(lang, expected) {
		t.Errorf("GetLanguage should return %v, but got %v", expected, lang)
	}
}

func Test_GetLanguage_ShouldReturnStructError(t *testing.T) {
	expected := errors.New("golang")

	mr := mockRepository{err: expected}

	_, err := mr.GetLanguage("")
	if !reflect.DeepEqual(err, expected) {
		t.Errorf("GetLanguage should return %v, but got %v", expected, err)
	}
}

func Test_PostLanguage_ShouldReturnStructId(t *testing.T) {
	id := "id"
	mr := mockRepository{id: id}

	result, err := mr.PostLanguage(models.Language{})
	if err != nil {
		t.Errorf("PostLanguage should not return error, but got %v", err)
	}

	if !reflect.DeepEqual(result, id) {
		t.Errorf("PostLanguage should return %v, but got %v", id, result)
	}
}

func Test_PostLanguage_ShouldReturnStructError(t *testing.T) {
	expected := errors.New("golang")

	mr := mockRepository{err: expected}

	_, err := mr.PostLanguage(models.Language{})
	if !reflect.DeepEqual(err, expected) {
		t.Errorf("PostLanguage should return %v, but got %v", expected, err)
	}
}

func Test_PutLanguage_ShouldReturnStructId(t *testing.T) {
	mr := mockRepository{isUpserted: true}

	isUpserted, err := mr.PutLanguage("", models.Language{})
	if err != nil {
		t.Errorf("PutLanguage should not return error, but got %v", err)
	}

	if !isUpserted {
		t.Errorf("PutLanguage should return %v, but got %v", true, isUpserted)
	}
}

func Test_PutLanguage_ShouldReturnStructError(t *testing.T) {
	expected := errors.New("golang")

	mr := mockRepository{err: expected}

	_, err := mr.PutLanguage("", models.Language{})
	if !reflect.DeepEqual(err, expected) {
		t.Errorf("PutLanguage should return %v, but got %v", expected, err)
	}
}

func Test_PatchLanguage_ShouldReturnStructError(t *testing.T) {
	expected := errors.New("golang")

	mr := mockRepository{err: expected}

	err := mr.PatchLanguage("", models.Language{})
	if !reflect.DeepEqual(err, expected) {
		t.Errorf("PatchLanguage should return %v, but got %v", expected, err)
	}
}

func Test_DeleteLanguage_ShouldReturnStructError(t *testing.T) {
	expected := errors.New("golang")

	mr := mockRepository{err: expected}

	err := mr.DeleteLanguage("")
	if !reflect.DeepEqual(err, expected) {
		t.Errorf("DeleteLanguage should return %v, but got %v", expected, err)
	}
}
