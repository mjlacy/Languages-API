package router

import (
	"languages-api/internal/config"
	"languages-api/internal/controller"

	"reflect"
	"testing"
)

func Test_CreateHandler_ShouldReturnRouter(t *testing.T) {
	r := CreateHandler(controller.New(config.Config{}), nil)

	if reflect.TypeOf(r).String() != "*mux.Router" {
		t.Errorf("CreateHandler returned %s, expected *mux.Router", reflect.TypeOf(r).String())
	}
}
