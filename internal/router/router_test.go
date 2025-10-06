package router

import (
	"languages-api/internal/config"
	"languages-api/internal/controller"

	"reflect"
	"testing"

	"github.com/gorilla/mux"
)

func Test_CreateHandler_ShouldReturnRouter(t *testing.T) {
	r := CreateHandler(controller.New(config.Config{}), nil)

	if reflect.TypeOf(r) != reflect.TypeOf(mux.NewRouter()) {
		t.Errorf("CreateHandler returned %+v, expected %+v", reflect.TypeOf(r), reflect.TypeOf(mux.NewRouter()))
	}
}
