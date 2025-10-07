package repo

import (
	"languages-api/internal/models"

	"errors"
	"reflect"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Test_Ping_ShouldReturnRepoError(t *testing.T) {
	expected := errors.New("ping error")

	err := (&MockRepo{Err: expected}).Ping()
	if !errors.Is(err, expected) {
		t.Errorf("expected %v, got %v", expected, err)
	}
}

func Test_GetLanguages_ShouldReturnRepoLanguages(t *testing.T) {
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

	result, err := (&MockRepo{languages: expected}).GetLanguages(models.Language{})
	if err != nil {
		t.Error("Error getting languages:", err)
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func Test_GetLanguages_ShouldReturnRepoError(t *testing.T) {
	expected := errors.New("getLanguages error")

	_, err := (&MockRepo{Err: expected}).GetLanguages(models.Language{})
	if !errors.Is(err, expected) {
		t.Errorf("expected %v, got %v", expected, err)
	}
}

func Test_GetLanguage_ShouldReturnRepoLanguage(t *testing.T) {
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

	result, err := (&MockRepo{language: expected}).GetLanguage("")
	if err != nil {
		t.Error("Error getting language:", err)
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func Test_GetLanguage_ShouldReturnRepoError(t *testing.T) {
	expected := errors.New("getLanguage error")

	_, err := (&MockRepo{Err: expected}).GetLanguage("")
	if !errors.Is(err, expected) {
		t.Errorf("expected %v, got %v", expected, err)
	}
}

func Test_PostLanguage_ShouldReturnRepoId(t *testing.T) {
	expected := "id"

	result, err := (&MockRepo{id: expected}).PostLanguage(models.Language{})
	if err != nil {
		t.Error("Error posting language id:", err)
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func Test_PostLanguage_ShouldReturnRepoError(t *testing.T) {
	expected := errors.New("postLanguage error")

	_, err := (&MockRepo{Err: expected}).PostLanguage(models.Language{})
	if !errors.Is(err, expected) {
		t.Errorf("expected %v, got %v", expected, err)
	}
}

func Test_PutLanguage_ShouldReturnRepoIsUpserted(t *testing.T) {
	result, err := (&MockRepo{isUpserted: true}).PutLanguage("", models.Language{})
	if err != nil {
		t.Error("Error posting language id:", err)
	}

	if !result {
		t.Errorf("expected true, got %v", result)
	}
}

func Test_PutLanguage_ShouldReturnRepoError(t *testing.T) {
	expected := errors.New("putLanguage error")

	_, err := (&MockRepo{Err: expected}).PutLanguage("", models.Language{})
	if !errors.Is(err, expected) {
		t.Errorf("expected %v, got %v", expected, err)
	}
}

func Test_PatchLanguage_ShouldReturnRepoError(t *testing.T) {
	expected := errors.New("patchLanguage error")

	err := (&MockRepo{Err: expected}).PatchLanguage("", models.Language{})
	if !errors.Is(err, expected) {
		t.Errorf("expected %v, got %v", expected, err)
	}
}

func Test_DeleteLanguage_ShouldReturnRepoError(t *testing.T) {
	expected := errors.New("deleteLanguage error")

	err := (&MockRepo{Err: expected}).DeleteLanguage("")
	if !errors.Is(err, expected) {
		t.Errorf("expected %v, got %v", expected, err)
	}
}

func Test_Close_ShouldReturnRepoError(t *testing.T) {
	expected := errors.New("close error")

	err := (&MockRepo{Err: expected}).Close()
	if !errors.Is(err, expected) {
		t.Errorf("expected %v, got %v", expected, err)
	}
}
