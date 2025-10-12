package config

import (
	"errors"
	"io/fs"
	"reflect"
	"testing"

	"github.com/spf13/viper"
)

func Test_New_ShouldReturnErrorOnReadInConfigError(t *testing.T) {
	_, err := New()
	if err == nil {
		t.Error("No error thrown where one was expected")
	}

	var pathError *fs.PathError

	if !errors.As(err, &pathError) {
		t.Errorf("Error should be of type fs.PathError, got %v", err)
	}
}

func Test_New_ShouldReturnConfigOnSuccess(t *testing.T) {
	expected := Config{
		AppName:    AppName,
		ConfigPath: "../../config.json",
		Table:      "languages",
		Database:   "languages",
		DBURL:      "root:password@tcp(host.docker.internal:3306)/languages?parseTime=true",
		Port:       "8080",
		Version:    Version,
	}

	viper.Set("ConfigPath", "../../config.json")

	cfg, err := New()
	if err != nil {
		t.Errorf("Unexpected error while creating new config: %s", err)
	}

	if !reflect.DeepEqual(cfg, expected) {
		t.Errorf("Expected %v, got %v", expected, cfg)
	}
}
