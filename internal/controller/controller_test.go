package controller

import (
	"languages-api/internal/config"

	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
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
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := ctrl.HealthCheckHandler(mockRepository{err: errors.New("ping")})

	handler.ServeHTTP(rr, req)

	if rr.Header()["Content-Type"][0] != expected {
		t.Errorf("Expected Content-Type of %s, but got %v", expected, rr.Header()["Content-Type"][0])
	}
}

func Test_HealthCheckHandler_ShouldReturnStatus500OnPingError(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/health", nil)
	if err != nil {
		t.Fatal(err)
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
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := ctrl.HealthCheckHandler(mockRepository{err: errors.New("ping")})

	handler.ServeHTTP(rr, req)

	var respBody HealthCheck

	err = json.Unmarshal(rr.Body.Bytes(), &respBody)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(respBody, expected) {
		t.Errorf("Expected %+v but got %+v", expected, respBody)
	}
}

func Test_HealthCheckHandler_ShouldReturnStatus200OnSuccess(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/health", nil)
	if err != nil {
		t.Fatal(err)
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
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := ctrl.HealthCheckHandler(mockRepository{err: nil})

	handler.ServeHTTP(rr, req)

	var respBody HealthCheck

	err = json.Unmarshal(rr.Body.Bytes(), &respBody)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(respBody, expected) {
		t.Errorf("Expected %+v but got %+v", expected, respBody)
	}
}
