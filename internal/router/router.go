package router

import (
	"languages-api/internal/controller"
	"languages-api/internal/mgo"

	"net/http"

	"github.com/gorilla/mux"
)

func CreateHandler(ctrl controller.APIController, repo *mgo.Repo) http.Handler {
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/health", ctrl.HealthCheckHandler(repo)).Methods(http.MethodGet)
	r.HandleFunc("/", ctrl.GetLanguagesHandler(repo)).Methods(http.MethodGet)
	r.HandleFunc("/{id}", ctrl.GetLanguageHandler(repo)).Methods(http.MethodGet)
	r.HandleFunc("/", ctrl.CreateLanguageHandler(repo)).Methods(http.MethodPost)
	r.HandleFunc("/{id}", ctrl.UpsertLanguageHandler(repo)).Methods(http.MethodPut)
	r.HandleFunc("/{id}", ctrl.UpdateLanguageHandler(repo)).Methods(http.MethodPatch)
	r.HandleFunc("/{id}", ctrl.DeleteLanguageHandler(repo)).Methods(http.MethodDelete)
	r.NotFoundHandler = http.HandlerFunc(ctrl.NotFoundPageHandler)

	return r
}
