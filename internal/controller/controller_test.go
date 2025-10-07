package controller

import (
	"languages-api/internal/config"
	"languages-api/internal/models"

	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	cfg = config.Config{
		AppName:    config.AppName,
		ConfigPath: "",
		Collection: "",
		Database:   "",
		DBURL:      "",
		Port:       "",
		Version:    config.Version,
	}
	ctrl = Controller{Config: cfg}
)

func Test_New_ShouldReturnController(t *testing.T) {
	expected := &Controller{
		Config: cfg,
	}

	ctrl := New(cfg)
	if !reflect.DeepEqual(ctrl, expected) {
		t.Errorf("New should return %v, but returned %v", expected, ctrl)
	}
}

func Test_HealthCheckHandler_ShouldHaveContentTypeHeader(t *testing.T) {
	expected := "application/json"

	req, err := http.NewRequest(http.MethodGet, "/health", nil)
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	handler := ctrl.HealthCheckHandler(mockRepository{err: errors.New("ping")})

	handler.ServeHTTP(rr, req)

	ct := rr.Header().Get("Content-Type")

	if ct != expected {
		t.Errorf("Expected Content-Type of %s, but got %v", expected, ct)
	}
}

func Test_HealthCheckHandler_ShouldReturnStatus500OnPingError(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/health", nil)
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	handler := ctrl.HealthCheckHandler(mockRepository{err: errors.New("ping")})

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Expected 500 but got %v", rr.Code)
	}
}

func Test_HealthCheckHandler_ShouldReturnErrorMessageOnPingError(t *testing.T) {
	expected := HealthCheck{
		Info: Info{
			ApplicationName: ctrl.Config.AppName,
			Version:         ctrl.Config.Version,
		},
		HealthCodes: HealthCodes{
			Application:     http.StatusText(http.StatusOK),
			MongoConnection: http.StatusText(http.StatusInternalServerError),
		},
	}

	req, err := http.NewRequest(http.MethodGet, "/health", nil)
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	handler := ctrl.HealthCheckHandler(mockRepository{err: errors.New("ping")})

	handler.ServeHTTP(rr, req)

	var respBody HealthCheck

	err = json.Unmarshal(rr.Body.Bytes(), &respBody)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(respBody, expected) {
		t.Errorf("Expected %+v but got %+v", expected, respBody)
	}
}

func Test_HealthCheckHandler_ShouldReturnStatus200OnSuccess(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/health", nil)
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	handler := ctrl.HealthCheckHandler(mockRepository{err: nil})

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected 200 but got %v", rr.Code)
	}
}

func Test_HealthCheckHandler_ShouldReturnSuccessMessageOnSuccess(t *testing.T) {
	expected := HealthCheck{
		Info: Info{
			ApplicationName: ctrl.Config.AppName,
			Version:         ctrl.Config.Version,
		},
		HealthCodes: HealthCodes{
			Application:     http.StatusText(http.StatusOK),
			MongoConnection: http.StatusText(http.StatusOK),
		},
	}

	req, err := http.NewRequest(http.MethodGet, "/health", nil)
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	handler := ctrl.HealthCheckHandler(mockRepository{err: nil})

	handler.ServeHTTP(rr, req)

	var respBody HealthCheck

	err = json.Unmarshal(rr.Body.Bytes(), &respBody)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(respBody, expected) {
		t.Errorf("Expected %+v but got %+v", expected, respBody)
	}
}

func Test_GetLanguagesHandler_ShouldHaveContentTypeHeaderOnQueryDecodeError(t *testing.T) {
	expected := "text/plain; charset=utf-8"

	req, err := http.NewRequest(http.MethodGet, "/?fake=true", nil)
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	handler := ctrl.GetLanguagesHandler(mockRepository{})

	handler.ServeHTTP(rr, req)

	ct := rr.Header().Get("Content-Type")

	if ct != expected {
		t.Errorf("Expected Content-Type of %s, but got %v", expected, ct)
	}
}

func Test_GetLanguagesHandler_ShouldReturnStatus400OnQueryDecodeError(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/?fake=true", nil)
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	handler := ctrl.GetLanguagesHandler(mockRepository{})

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 but got %v", rr.Code)
	}
}

func Test_GetLanguagesHandler_ShouldReturnErrorMessageOnQueryDecodeError(t *testing.T) {
	expected := "Invalid query string"

	req, err := http.NewRequest(http.MethodGet, "/?fake=true", nil)
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	handler := ctrl.GetLanguagesHandler(mockRepository{})

	handler.ServeHTTP(rr, req)

	respBody := rr.Body.String()

	if !reflect.DeepEqual(respBody, expected) {
		t.Errorf("Expected %+v but got %+v", expected, respBody)
	}
}

func Test_GetLanguagesHandler_ShouldHaveContentTypeHeaderOnGetLanguagesError(t *testing.T) {
	expected := "text/plain; charset=utf-8"

	req, err := http.NewRequest(http.MethodGet, "/?fake=true", nil)
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	handler := ctrl.GetLanguagesHandler(mockRepository{errs: []error{errors.New("GetLanguages")}})

	handler.ServeHTTP(rr, req)

	ct := rr.Header().Get("Content-Type")

	if ct != expected {
		t.Errorf("Expected Content-Type of %s, but got %v", expected, ct)
	}
}

func Test_GetLanguagesHandler_ShouldReturnStatus500OnGetLanguagesError(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/?Creators=Robert Griesemer&Extensions=.go", nil)
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	handler := ctrl.GetLanguagesHandler(mockRepository{errs: []error{errors.New("GetLanguages")}})

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Expected 500 but got %v", rr.Code)
	}
}

func Test_GetLanguagesHandler_ShouldReturnErrorMessageOnGetLanguagesError(t *testing.T) {
	expected := "An error occurred processing this request"

	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	handler := ctrl.GetLanguagesHandler(mockRepository{errs: []error{errors.New("GetLanguages")}})

	handler.ServeHTTP(rr, req)

	respBody := rr.Body.String()

	if !reflect.DeepEqual(respBody, expected) {
		t.Errorf("Expected: %+v, but got: %+v", expected, respBody)
	}
}

func Test_GetLanguagesHandler_ShouldHaveContentTypeHeaderOnSuccess(t *testing.T) {
	expected := "application/json"

	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	handler := ctrl.GetLanguagesHandler(mockRepository{ls: models.Languages{}})

	handler.ServeHTTP(rr, req)

	ct := rr.Header().Get("Content-Type")

	if ct != expected {
		t.Errorf("Expected Content-Type of %s, but got %v", expected, ct)
	}
}

func Test_GetLanguagesHandler_ShouldReturnStatus200OnSuccess(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	handler := ctrl.GetLanguagesHandler(mockRepository{})

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected 200 but got %v", rr.Code)
	}
}

func Test_GetLanguagesHandler_ShouldReturnLanguagesOnSuccess(t *testing.T) {
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

	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	handler := ctrl.GetLanguagesHandler(mockRepository{ls: expected})

	handler.ServeHTTP(rr, req)

	var respBody models.Languages

	err = json.Unmarshal(rr.Body.Bytes(), &respBody)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(respBody, expected) {
		t.Errorf("Expected %+v but got %+v", expected, respBody)
	}
}

func Test_GetLanguageHandler_ShouldHaveContentTypeHeaderOnInvalidIdError(t *testing.T) {
	expected := "text/plain; charset=utf-8"

	req, err := http.NewRequest(http.MethodGet, "/1", nil)
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	handler := ctrl.GetLanguageHandler(mockRepository{err: models.ErrInvalidId})

	handler.ServeHTTP(rr, req)

	ct := rr.Header().Get("Content-Type")

	if ct != expected {
		t.Errorf("Expected Content-Type of %s, but got %v", expected, ct)
	}
}

func Test_GetLanguageHandler_ShouldReturnStatus400OnInvalidIdError(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/1", nil)
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	handler := ctrl.GetLanguageHandler(mockRepository{err: models.ErrInvalidId})

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 but got %v", rr.Code)
	}
}

func Test_GetLanguageHandler_ShouldReturnErrorMessageOnInvalidIdError(t *testing.T) {
	expected := "The given id is not a valid id"

	req, err := http.NewRequest(http.MethodGet, "/1", nil)
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	handler := ctrl.GetLanguageHandler(mockRepository{err: models.ErrInvalidId})

	handler.ServeHTTP(rr, req)

	respBody := rr.Body.String()

	if !reflect.DeepEqual(respBody, expected) {
		t.Errorf("Expected %+v but got %+v", expected, respBody)
	}
}

func Test_GetLanguageHandler_ShouldHaveContentTypeHeaderOnNotFoundError(t *testing.T) {
	expected := "text/plain; charset=utf-8"

	req, err := http.NewRequest(http.MethodGet, "/1", nil)
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	handler := ctrl.GetLanguageHandler(mockRepository{err: models.ErrNotFound})

	handler.ServeHTTP(rr, req)

	ct := rr.Header().Get("Content-Type")

	if ct != expected {
		t.Errorf("Expected Content-Type of %s, but got %v", expected, ct)
	}
}

func Test_GetLanguageHandler_ShouldReturnStatus404OnNotFoundError(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/1", nil)
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	handler := ctrl.GetLanguageHandler(mockRepository{err: models.ErrNotFound})

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("Expected 404 but got %v", rr.Code)
	}
}

func Test_GetLanguageHandler_ShouldReturnErrorMessageOnNotFoundError(t *testing.T) {
	expected := "No language found with that id"

	req, err := http.NewRequest(http.MethodGet, "/1", nil)
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	handler := ctrl.GetLanguageHandler(mockRepository{err: models.ErrNotFound})

	handler.ServeHTTP(rr, req)

	respBody := rr.Body.String()

	if !reflect.DeepEqual(respBody, expected) {
		t.Errorf("Expected %+v but got %+v", expected, respBody)
	}
}

func Test_GetLanguageHandler_ShouldHaveContentTypeHeaderOnInternalError(t *testing.T) {
	expected := "text/plain; charset=utf-8"

	req, err := http.NewRequest(http.MethodGet, "/1", nil)
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	handler := ctrl.GetLanguageHandler(mockRepository{err: errors.New("internal server error")})

	handler.ServeHTTP(rr, req)

	ct := rr.Header().Get("Content-Type")

	if ct != expected {
		t.Errorf("Expected Content-Type of %s, but got %v", expected, ct)
	}
}

func Test_GetLanguageHandler_ShouldReturnStatus500OnInternalError(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/1", nil)
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	handler := ctrl.GetLanguageHandler(mockRepository{err: errors.New("internal server error")})

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Expected 500 but got %v", rr.Code)
	}
}

func Test_GetLanguageHandler_ShouldReturnErrorMessageOnInternalError(t *testing.T) {
	expected := "An error occurred processing this request"

	req, err := http.NewRequest(http.MethodGet, "/1", nil)
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	handler := ctrl.GetLanguageHandler(mockRepository{err: errors.New("internal server error")})

	handler.ServeHTTP(rr, req)

	respBody := rr.Body.String()

	if !reflect.DeepEqual(respBody, expected) {
		t.Errorf("Expected %+v but got %+v", expected, respBody)
	}
}

func Test_GetLanguageHandler_ShouldHaveContentTypeHeaderOnSuccess(t *testing.T) {
	expected := "application/json"

	req, err := http.NewRequest(http.MethodGet, "/1", nil)
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	handler := ctrl.GetLanguageHandler(mockRepository{})

	handler.ServeHTTP(rr, req)

	ct := rr.Header().Get("Content-Type")

	if ct != expected {
		t.Errorf("Expected Content-Type of %s, but got %v", expected, ct)
	}
}

func Test_GetLanguageHandler_ShouldReturnStatus200OnSuccess(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/1", nil)
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	handler := ctrl.GetLanguageHandler(mockRepository{})

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected 200 but got %v", rr.Code)
	}
}

func Test_GetLanguageHandler_ShouldReturnLanguageOnSuccess(t *testing.T) {
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

	req, err := http.NewRequest(http.MethodGet, "/1", nil)
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	handler := ctrl.GetLanguageHandler(mockRepository{l: expected})

	handler.ServeHTTP(rr, req)

	var respBody models.Language

	err = json.Unmarshal(rr.Body.Bytes(), &respBody)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(respBody, expected) {
		t.Errorf("Expected %+v but got %+v", expected, respBody)
	}
}

func Test_CreateLanguageHandler_ShouldHaveContentTypeHeaderOnDecodeError(t *testing.T) {
	expected := "text/plain; charset=utf-8"

	req, err := http.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte("Invalid request body")))
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	handler := ctrl.CreateLanguageHandler(mockRepository{})

	handler.ServeHTTP(rr, req)

	ct := rr.Header().Get("Content-Type")

	if ct != expected {
		t.Errorf("Expected Content-Type of %s, but got %v", expected, ct)
	}
}

func Test_CreateLanguageHandler_ShouldReturnStatus400OnDecodeError(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte("Invalid request body")))
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	handler := ctrl.CreateLanguageHandler(mockRepository{})

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 but got %v", rr.Code)
	}
}

func Test_CreateLanguageHandler_ShouldReturnErrorMessageOnDecodeError(t *testing.T) {
	expected := "Invalid request body"

	req, err := http.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(expected)))
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	handler := ctrl.CreateLanguageHandler(mockRepository{})

	handler.ServeHTTP(rr, req)

	respBody := rr.Body.String()

	if !reflect.DeepEqual(respBody, expected) {
		t.Errorf("Expected %+v but got %+v", expected, respBody)
	}
}

func Test_CreateLanguageHandler_ShouldHaveContentTypeHeaderOnInternalError(t *testing.T) {
	expected := "text/plain; charset=utf-8"

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

	reqBody, err := json.Marshal(lang)
	if err != nil {
		t.Error(err)
	}

	req, err := http.NewRequest(http.MethodPost, "/", bytes.NewReader(reqBody))
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	handler := ctrl.CreateLanguageHandler(mockRepository{err: errors.New("internal server error")})

	handler.ServeHTTP(rr, req)

	ct := rr.Header().Get("Content-Type")

	if ct != expected {
		t.Errorf("Expected Content-Type of %s, but got %v", expected, ct)
	}
}

func Test_CreateLanguageHandler_ShouldReturnStatus500OnInternalError(t *testing.T) {
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

	reqBody, err := json.Marshal(lang)
	if err != nil {
		t.Error(err)
	}

	req, err := http.NewRequest(http.MethodPost, "/", bytes.NewReader(reqBody))
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	handler := ctrl.CreateLanguageHandler(mockRepository{err: errors.New("internal server error")})

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Expected 500 but got %v", rr.Code)
	}
}

func Test_CreateLanguageHandler_ShouldReturnErrorMessageOnInternalError(t *testing.T) {
	expected := "An error occurred processing this request"

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

	reqBody, err := json.Marshal(lang)
	if err != nil {
		t.Error(err)
	}

	req, err := http.NewRequest(http.MethodPost, "/", bytes.NewReader(reqBody))
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	handler := ctrl.CreateLanguageHandler(mockRepository{err: errors.New("internal server error")})

	handler.ServeHTTP(rr, req)

	respBody := rr.Body.String()

	if !reflect.DeepEqual(respBody, expected) {
		t.Errorf("Expected %+v but got %+v", expected, respBody)
	}
}

func Test_CreateLanguageHandler_ShouldHaveLocationHeaderOnSuccess(t *testing.T) {
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

	expected := fmt.Sprintf("/%v", url.PathEscape(lang.Id.String()))

	reqBody, err := json.Marshal(lang)
	if err != nil {
		t.Error(err)
	}

	req, err := http.NewRequest(http.MethodPost, "/", bytes.NewReader(reqBody))
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	handler := ctrl.CreateLanguageHandler(mockRepository{id: lang.Id.String()})

	handler.ServeHTTP(rr, req)

	location := rr.Header().Get("Location")

	if location != expected {
		t.Errorf("Expected Location of %s, but got %v", expected, location)
	}
}

func Test_CreateLanguageHandler_ShouldReturnStatus201OnSuccess(t *testing.T) {
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

	reqBody, err := json.Marshal(lang)
	if err != nil {
		t.Error(err)
	}

	req, err := http.NewRequest(http.MethodPost, "/", bytes.NewReader(reqBody))
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	handler := ctrl.CreateLanguageHandler(mockRepository{})

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("Expected 201 but got %v", rr.Code)
	}
}

func Test_CreateLanguageHandler_ShouldReturnNoMessageOnSuccess(t *testing.T) {
	expected := ""

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

	reqBody, err := json.Marshal(lang)
	if err != nil {
		t.Error(err)
	}

	req, err := http.NewRequest(http.MethodPost, "/", bytes.NewReader(reqBody))
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	handler := ctrl.CreateLanguageHandler(mockRepository{})

	handler.ServeHTTP(rr, req)

	respBody := rr.Body.String()

	if !reflect.DeepEqual(respBody, expected) {
		t.Errorf("Expected %+v but got %+v", expected, respBody)
	}
}

func Test_UpsertLanguageHandler_ShouldHaveContentTypeHeaderOnDecodeError(t *testing.T) {
	expected := "text/plain; charset=utf-8"

	req, err := http.NewRequest(http.MethodPut, "/1", bytes.NewReader([]byte("Invalid request body")))
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	handler := ctrl.UpsertLanguageHandler(mockRepository{})

	handler.ServeHTTP(rr, req)

	ct := rr.Header().Get("Content-Type")

	if ct != expected {
		t.Errorf("Expected Content-Type of %s, but got %v", expected, ct)
	}
}

func Test_UpsertLanguageHandler_ShouldReturnStatus400OnDecodeError(t *testing.T) {
	req, err := http.NewRequest(http.MethodPut, "/1", bytes.NewReader([]byte("Invalid request body")))
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	handler := ctrl.UpsertLanguageHandler(mockRepository{})

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 but got %v", rr.Code)
	}
}

func Test_UpsertLanguageHandler_ShouldReturnErrorMessageOnDecodeError(t *testing.T) {
	expected := "Invalid request body"

	req, err := http.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(expected)))
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	handler := ctrl.UpsertLanguageHandler(mockRepository{})

	handler.ServeHTTP(rr, req)

	respBody := rr.Body.String()

	if !reflect.DeepEqual(respBody, expected) {
		t.Errorf("Expected %+v but got %+v", expected, respBody)
	}
}

func Test_UpsertLanguageHandler_ShouldHaveContentTypeHeaderOnInvalidIdError(t *testing.T) {
	expected := "text/plain; charset=utf-8"

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

	reqBody, err := json.Marshal(lang)
	if err != nil {
		t.Error(err)
	}

	req, err := http.NewRequest(http.MethodPut, "/1", bytes.NewReader(reqBody))
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	handler := ctrl.UpsertLanguageHandler(mockRepository{err: models.ErrInvalidId})

	handler.ServeHTTP(rr, req)

	ct := rr.Header().Get("Content-Type")

	if ct != expected {
		t.Errorf("Expected Content-Type of %s, but got %v", expected, ct)
	}
}

func Test_UpsertLanguageHandler_ShouldReturnStatus400OnInvalidIdError(t *testing.T) {
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

	reqBody, err := json.Marshal(lang)
	if err != nil {
		t.Error(err)
	}

	req, err := http.NewRequest(http.MethodPut, "/1", bytes.NewReader(reqBody))
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	handler := ctrl.UpsertLanguageHandler(mockRepository{err: models.ErrInvalidId})

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 but got %v", rr.Code)
	}
}

func Test_UpsertLanguageHandler_ShouldReturnErrorMessageOnInvalidIdError(t *testing.T) {
	expected := "The given id is not a valid id"

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

	reqBody, err := json.Marshal(lang)
	if err != nil {
		t.Error(err)
	}

	req, err := http.NewRequest(http.MethodPut, "/1", bytes.NewReader(reqBody))
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	handler := ctrl.UpsertLanguageHandler(mockRepository{err: models.ErrInvalidId})

	handler.ServeHTTP(rr, req)

	respBody := rr.Body.String()

	if !reflect.DeepEqual(respBody, expected) {
		t.Errorf("Expected %+v but got %+v", expected, respBody)
	}
}

func Test_UpsertLanguageHandler_ShouldHaveContentTypeHeaderOnInternalError(t *testing.T) {
	expected := "text/plain; charset=utf-8"

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

	reqBody, err := json.Marshal(lang)
	if err != nil {
		t.Error(err)
	}

	req, err := http.NewRequest(http.MethodPut, "/1", bytes.NewReader(reqBody))
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	handler := ctrl.UpsertLanguageHandler(mockRepository{err: errors.New("internal server error")})

	handler.ServeHTTP(rr, req)

	ct := rr.Header().Get("Content-Type")

	if ct != expected {
		t.Errorf("Expected Content-Type of %s, but got %v", expected, ct)
	}
}

func Test_UpsertLanguageHandler_ShouldReturnStatus500OnInternalError(t *testing.T) {
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

	reqBody, err := json.Marshal(lang)
	if err != nil {
		t.Error(err)
	}

	req, err := http.NewRequest(http.MethodPut, "/1", bytes.NewReader(reqBody))
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	handler := ctrl.UpsertLanguageHandler(mockRepository{err: errors.New("internal server error")})

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Expected 500 but got %v", rr.Code)
	}
}

func Test_UpsertLanguageHandler_ShouldReturnErrorMessageOnInternalError(t *testing.T) {
	expected := "An error occurred processing this request"

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

	reqBody, err := json.Marshal(lang)
	if err != nil {
		t.Error(err)
	}

	req, err := http.NewRequest(http.MethodPut, "/1", bytes.NewReader(reqBody))
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	handler := ctrl.UpsertLanguageHandler(mockRepository{err: errors.New("internal server error")})

	handler.ServeHTTP(rr, req)

	respBody := rr.Body.String()

	if !reflect.DeepEqual(respBody, expected) {
		t.Errorf("Expected %+v but got %+v", expected, respBody)
	}
}

func Test_UpsertLanguageHandler_ShouldHaveLocationHeaderOnIsUpsertedSuccess(t *testing.T) {
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

	expected := "/1"

	reqBody, err := json.Marshal(lang)
	if err != nil {
		t.Error(err)
	}

	req, err := http.NewRequest(http.MethodPut, "/1", bytes.NewReader(reqBody))
	if err != nil {
		t.Error(err)
	}

	req = mux.SetURLVars(req, map[string]string{"id": "1"})

	rr := httptest.NewRecorder()
	handler := ctrl.UpsertLanguageHandler(mockRepository{isUpserted: true})

	handler.ServeHTTP(rr, req)

	location := rr.Header().Get("Location")

	if location != expected {
		t.Errorf("Expected Location of %s, but got %v", expected, location)
	}
}

func Test_UpsertLanguageHandler_ShouldReturnStatus201OnIsUpsertedSuccess(t *testing.T) {
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

	reqBody, err := json.Marshal(lang)
	if err != nil {
		t.Error(err)
	}

	req, err := http.NewRequest(http.MethodPut, "/1", bytes.NewReader(reqBody))
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	handler := ctrl.CreateLanguageHandler(mockRepository{isUpserted: true})

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("Expected 201 but got %v", rr.Code)
	}
}

func Test_UpsertLanguageHandler_ShouldReturnStatus200OnNonIsUpsertedSuccess(t *testing.T) {
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

	reqBody, err := json.Marshal(lang)
	if err != nil {
		t.Error(err)
	}

	req, err := http.NewRequest(http.MethodPut, "/1", bytes.NewReader(reqBody))
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	handler := ctrl.UpsertLanguageHandler(mockRepository{isUpserted: false})

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected 200 but got %v", rr.Code)
	}
}

func Test_UpsertLanguageHandler_ShouldReturnNoMessageOnSuccess(t *testing.T) {
	expected := ""

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

	reqBody, err := json.Marshal(lang)
	if err != nil {
		t.Error(err)
	}

	req, err := http.NewRequest(http.MethodPut, "/1", bytes.NewReader(reqBody))
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	handler := ctrl.UpsertLanguageHandler(mockRepository{})

	handler.ServeHTTP(rr, req)

	respBody := rr.Body.String()

	if !reflect.DeepEqual(respBody, expected) {
		t.Errorf("Expected %+v but got %+v", expected, respBody)
	}
}

func Test_UpdateLanguageHandler_ShouldHaveContentTypeHeaderOnDecodeError(t *testing.T) {
	expected := "text/plain; charset=utf-8"

	req, err := http.NewRequest(http.MethodPatch, "/1", bytes.NewReader([]byte("Invalid request body")))
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	handler := ctrl.UpdateLanguageHandler(mockRepository{})

	handler.ServeHTTP(rr, req)

	ct := rr.Header().Get("Content-Type")

	if ct != expected {
		t.Errorf("Expected Content-Type of %s, but got %v", expected, ct)
	}
}

func Test_UpdateLanguageHandler_ShouldReturnStatus400OnDecodeError(t *testing.T) {
	req, err := http.NewRequest(http.MethodPatch, "/1", bytes.NewReader([]byte("Invalid request body")))
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	handler := ctrl.UpdateLanguageHandler(mockRepository{})

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 but got %v", rr.Code)
	}
}

func Test_UpdateLanguageHandler_ShouldReturnErrorMessageOnDecodeError(t *testing.T) {
	expected := "Invalid request body"

	req, err := http.NewRequest(http.MethodPatch, "/", bytes.NewReader([]byte(expected)))
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	handler := ctrl.UpdateLanguageHandler(mockRepository{})

	handler.ServeHTTP(rr, req)

	respBody := rr.Body.String()

	if !reflect.DeepEqual(respBody, expected) {
		t.Errorf("Expected %+v but got %+v", expected, respBody)
	}
}

func Test_UpdateLanguageHandler_ShouldHaveContentTypeHeaderOnInvalidIdError(t *testing.T) {
	expected := "text/plain; charset=utf-8"

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

	reqBody, err := json.Marshal(lang)
	if err != nil {
		t.Error(err)
	}

	req, err := http.NewRequest(http.MethodPatch, "/1?Creators=Robert Griesemer&Extensions=.go", bytes.NewReader(reqBody))
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	handler := ctrl.UpdateLanguageHandler(mockRepository{err: models.ErrInvalidId})

	handler.ServeHTTP(rr, req)

	ct := rr.Header().Get("Content-Type")

	if ct != expected {
		t.Errorf("Expected Content-Type of %s, but got %v", expected, ct)
	}
}

func Test_UpdateLanguageHandler_ShouldReturnStatus400OnInvalidIdError(t *testing.T) {
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

	reqBody, err := json.Marshal(lang)
	if err != nil {
		t.Error(err)
	}

	req, err := http.NewRequest(http.MethodPatch, "/1", bytes.NewReader(reqBody))
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	handler := ctrl.UpdateLanguageHandler(mockRepository{err: models.ErrInvalidId})

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 but got %v", rr.Code)
	}
}

func Test_UpdateLanguageHandler_ShouldReturnErrorMessageOnInvalidIdError(t *testing.T) {
	expected := "The given id is not a valid id"

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

	reqBody, err := json.Marshal(lang)
	if err != nil {
		t.Error(err)
	}

	req, err := http.NewRequest(http.MethodPatch, "/1", bytes.NewReader(reqBody))
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	handler := ctrl.UpdateLanguageHandler(mockRepository{err: models.ErrInvalidId})

	handler.ServeHTTP(rr, req)

	respBody := rr.Body.String()

	if !reflect.DeepEqual(respBody, expected) {
		t.Errorf("Expected %+v but got %+v", expected, respBody)
	}
}

func Test_UpdateLanguageHandler_ShouldHaveContentTypeHeaderOnNotFoundError(t *testing.T) {
	expected := "text/plain; charset=utf-8"

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

	reqBody, err := json.Marshal(lang)
	if err != nil {
		t.Error(err)
	}

	req, err := http.NewRequest(http.MethodPatch, "/1", bytes.NewReader(reqBody))
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	handler := ctrl.UpdateLanguageHandler(mockRepository{err: models.ErrNotFound})

	handler.ServeHTTP(rr, req)

	ct := rr.Header().Get("Content-Type")

	if ct != expected {
		t.Errorf("Expected Content-Type of %s, but got %v", expected, ct)
	}
}

func Test_UpdateLanguageHandler_ShouldReturnStatus404OnNotFoundError(t *testing.T) {
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

	reqBody, err := json.Marshal(lang)
	if err != nil {
		t.Error(err)
	}

	req, err := http.NewRequest(http.MethodPatch, "/1", bytes.NewReader(reqBody))
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	handler := ctrl.UpdateLanguageHandler(mockRepository{err: models.ErrNotFound})

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("Expected 404 but got %v", rr.Code)
	}
}

func Test_UpdateLanguageHandler_ShouldReturnErrorMessageOnNotFoundError(t *testing.T) {
	expected := "No language found with that id to update"

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

	reqBody, err := json.Marshal(lang)
	if err != nil {
		t.Error(err)
	}

	req, err := http.NewRequest(http.MethodPatch, "/1", bytes.NewReader(reqBody))
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	handler := ctrl.UpdateLanguageHandler(mockRepository{err: models.ErrNotFound})

	handler.ServeHTTP(rr, req)

	respBody := rr.Body.String()

	if !reflect.DeepEqual(respBody, expected) {
		t.Errorf("Expected %+v but got %+v", expected, respBody)
	}
}

func Test_UpdateLanguageHandler_ShouldHaveContentTypeHeaderOnInternalError(t *testing.T) {
	expected := "text/plain; charset=utf-8"

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

	reqBody, err := json.Marshal(lang)
	if err != nil {
		t.Error(err)
	}

	req, err := http.NewRequest(http.MethodPatch, "/1", bytes.NewReader(reqBody))
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	handler := ctrl.UpdateLanguageHandler(mockRepository{err: errors.New("internal server error")})

	handler.ServeHTTP(rr, req)

	ct := rr.Header().Get("Content-Type")

	if ct != expected {
		t.Errorf("Expected Content-Type of %s, but got %v", expected, ct)
	}
}

func Test_UpdateLanguageHandler_ShouldReturnStatus500OnInternalError(t *testing.T) {
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

	reqBody, err := json.Marshal(lang)
	if err != nil {
		t.Error(err)
	}

	req, err := http.NewRequest(http.MethodPatch, "/1", bytes.NewReader(reqBody))
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	handler := ctrl.UpdateLanguageHandler(mockRepository{err: errors.New("internal server error")})

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Expected 500 but got %v", rr.Code)
	}
}

func Test_UpdateLanguageHandler_ShouldReturnErrorMessageOnInternalError(t *testing.T) {
	expected := "An error occurred processing this request"

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

	reqBody, err := json.Marshal(lang)
	if err != nil {
		t.Error(err)
	}

	req, err := http.NewRequest(http.MethodPatch, "/1", bytes.NewReader(reqBody))
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	handler := ctrl.UpdateLanguageHandler(mockRepository{err: errors.New("internal server error")})

	handler.ServeHTTP(rr, req)

	respBody := rr.Body.String()

	if !reflect.DeepEqual(respBody, expected) {
		t.Errorf("Expected %+v but got %+v", expected, respBody)
	}
}

func Test_UpdateLanguageHandler_ShouldReturnStatus200OnSuccess(t *testing.T) {
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

	reqBody, err := json.Marshal(lang)
	if err != nil {
		t.Error(err)
	}

	req, err := http.NewRequest(http.MethodPatch, "/1", bytes.NewReader(reqBody))
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	handler := ctrl.UpdateLanguageHandler(mockRepository{})

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected 200 but got %v", rr.Code)
	}
}

func Test_UpdateLanguageHandler_ShouldReturnNoMessageOnSuccess(t *testing.T) {
	expected := ""

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

	reqBody, err := json.Marshal(lang)
	if err != nil {
		t.Error(err)
	}

	req, err := http.NewRequest(http.MethodPatch, "/1", bytes.NewReader(reqBody))
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	handler := ctrl.UpdateLanguageHandler(mockRepository{})

	handler.ServeHTTP(rr, req)

	respBody := rr.Body.String()

	if !reflect.DeepEqual(respBody, expected) {
		t.Errorf("Expected %+v but got %+v", expected, respBody)
	}
}

func Test_DeleteLanguageHandler_ShouldHaveContentTypeHeaderOnInvalidIdError(t *testing.T) {
	expected := "text/plain; charset=utf-8"

	req, err := http.NewRequest(http.MethodDelete, "/1", nil)
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	handler := ctrl.DeleteLanguageHandler(mockRepository{err: models.ErrInvalidId})

	handler.ServeHTTP(rr, req)

	ct := rr.Header().Get("Content-Type")

	if ct != expected {
		t.Errorf("Expected Content-Type of %s, but got %v", expected, ct)
	}
}

func Test_DeleteLanguageHandler_ShouldReturnStatus400OnInvalidIdError(t *testing.T) {
	req, err := http.NewRequest(http.MethodDelete, "/1", nil)
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	handler := ctrl.DeleteLanguageHandler(mockRepository{err: models.ErrInvalidId})

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 but got %v", rr.Code)
	}
}

func Test_DeleteLanguageHandler_ShouldReturnErrorMessageOnInvalidIdError(t *testing.T) {
	expected := "The given id is not a valid id"

	req, err := http.NewRequest(http.MethodDelete, "/1", nil)
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	handler := ctrl.DeleteLanguageHandler(mockRepository{err: models.ErrInvalidId})

	handler.ServeHTTP(rr, req)

	respBody := rr.Body.String()

	if !reflect.DeepEqual(respBody, expected) {
		t.Errorf("Expected %+v but got %+v", expected, respBody)
	}
}

func Test_DeleteLanguageHandler_ShouldHaveContentTypeHeaderOnNotFoundError(t *testing.T) {
	expected := "text/plain; charset=utf-8"

	req, err := http.NewRequest(http.MethodDelete, "/1", nil)
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	handler := ctrl.DeleteLanguageHandler(mockRepository{err: models.ErrNotFound})

	handler.ServeHTTP(rr, req)

	ct := rr.Header().Get("Content-Type")

	if ct != expected {
		t.Errorf("Expected Content-Type of %s, but got %v", expected, ct)
	}
}

func Test_DeleteLanguageHandler_ShouldReturnStatus404OnNotFoundError(t *testing.T) {
	req, err := http.NewRequest(http.MethodDelete, "/1", nil)
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	handler := ctrl.DeleteLanguageHandler(mockRepository{err: models.ErrNotFound})

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("Expected 404 but got %v", rr.Code)
	}
}

func Test_DeleteLanguageHandler_ShouldReturnErrorMessageOnNotFoundError(t *testing.T) {
	expected := "No language found with that id to delete"

	req, err := http.NewRequest(http.MethodDelete, "/1", nil)
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	handler := ctrl.DeleteLanguageHandler(mockRepository{err: models.ErrNotFound})

	handler.ServeHTTP(rr, req)

	respBody := rr.Body.String()

	if !reflect.DeepEqual(respBody, expected) {
		t.Errorf("Expected %+v but got %+v", expected, respBody)
	}
}

func Test_DeleteLanguageHandler_ShouldHaveContentTypeHeaderOnInternalError(t *testing.T) {
	expected := "text/plain; charset=utf-8"

	req, err := http.NewRequest(http.MethodDelete, "/1", nil)
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	handler := ctrl.DeleteLanguageHandler(mockRepository{err: errors.New("internal server error")})

	handler.ServeHTTP(rr, req)

	ct := rr.Header().Get("Content-Type")

	if ct != expected {
		t.Errorf("Expected Content-Type of %s, but got %v", expected, ct)
	}
}

func Test_DeleteLanguageHandler_ShouldReturnStatus500OnInternalError(t *testing.T) {
	req, err := http.NewRequest(http.MethodDelete, "/1", nil)
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	handler := ctrl.DeleteLanguageHandler(mockRepository{err: errors.New("internal server error")})

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Expected 500 but got %v", rr.Code)
	}
}

func Test_DeleteLanguageHandler_ShouldReturnErrorMessageOnInternalError(t *testing.T) {
	expected := "An error occurred processing this request"

	req, err := http.NewRequest(http.MethodDelete, "/1", nil)
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	handler := ctrl.DeleteLanguageHandler(mockRepository{err: errors.New("internal server error")})

	handler.ServeHTTP(rr, req)

	respBody := rr.Body.String()

	if !reflect.DeepEqual(respBody, expected) {
		t.Errorf("Expected %+v but got %+v", expected, respBody)
	}
}

func Test_DeleteLanguageHandler_ShouldReturnStatus204OnSuccess(t *testing.T) {
	req, err := http.NewRequest(http.MethodDelete, "/1", nil)
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	handler := ctrl.DeleteLanguageHandler(mockRepository{})

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Errorf("Expected 204 but got %v", rr.Code)
	}
}

func Test_DeleteLanguageHandler_ShouldReturnNoMessageOnSuccess(t *testing.T) {
	expected := ""

	req, err := http.NewRequest(http.MethodDelete, "/1", nil)
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	handler := ctrl.DeleteLanguageHandler(mockRepository{})

	handler.ServeHTTP(rr, req)

	respBody := rr.Body.String()

	if !reflect.DeepEqual(respBody, expected) {
		t.Errorf("Expected %+v but got %+v", expected, respBody)
	}
}

func Test_NotFoundPageHandler_ShouldHaveContentTypeHeader(t *testing.T) {
	expected := "text/plain; charset=utf-8"

	req, err := http.NewRequest(http.MethodPost, "/1", nil)
	if err != nil {
		t.Error(err)
	}

	mrw := mockResponseWriter{header: http.Header{}}
	ctrl.NotFoundPageHandler(&mrw, req)

	ct := mrw.Header().Get("Content-Type")

	if ct != expected {
		t.Errorf("Expected Content-Type of %s, but got %v", expected, ct)
	}
}

func Test_NotFoundPageHandler_ShouldReturnStatus404(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "/1", nil)
	if err != nil {
		t.Error(err)
	}

	mrw := mockResponseWriter{header: http.Header{}, err: errors.New("404 page not found")}
	ctrl.NotFoundPageHandler(&mrw, req)

	if mrw.statusCode != http.StatusNotFound {
		t.Errorf("Expected 404 but got %v", mrw.statusCode)
	}
}

func Test_NotFoundPageHandler_ShouldWriteInvalidURLMessage(t *testing.T) {
	expected := "You have accessed an invalid URL"

	req, err := http.NewRequest(http.MethodPost, "/1", nil)
	if err != nil {
		t.Error(err)
	}

	mrw := mockResponseWriter{header: http.Header{}}
	ctrl.NotFoundPageHandler(&mrw, req)

	if !reflect.DeepEqual(mrw.message, expected) {
		t.Errorf("Expected %+v but got %+v", expected, mrw.message)
	}
}
