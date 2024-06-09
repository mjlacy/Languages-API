package controller

import (
	"languages-api/internal/config"
	"languages-api/internal/mgo"
	"languages-api/internal/models"

	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/rs/zerolog/log"
)

type APIController interface {
	HealthCheckHandler(repo mgo.Repository) http.HandlerFunc
	GetLanguagesHandler(repo mgo.Repository) http.HandlerFunc
	GetLanguageHandler(repo mgo.Repository) http.HandlerFunc
	CreateLanguageHandler(repo mgo.Repository) http.HandlerFunc
	UpsertLanguageHandler(repo mgo.Repository) http.HandlerFunc
	UpdateLanguageHandler(repo mgo.Repository) http.HandlerFunc
	DeleteLanguageHandler(repo mgo.Repository) http.HandlerFunc
	NotFoundPageHandler(w http.ResponseWriter, r *http.Request)
}

type Info struct {
	ApplicationName string `json:"ApplicationName"`
	Version         string `json:"Version"`
}

type HealthCodes struct {
	Application     string `json:"Application"`
	MongoConnection string `json:"MongoConnection"`
}

type HealthCheck struct {
	Info        Info        `json:"Info"`
	HealthCodes HealthCodes `json:"HealthCodes"`
}

type Controller struct {
	Config config.Config
}

func New(cfg config.Config) *Controller {
	return &Controller{
		Config: cfg,
	}
}

func (ctrl *Controller) HealthCheckHandler(repo mgo.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mongoStatus := http.StatusOK
		err := repo.Ping()
		if err != nil {
			mongoStatus = http.StatusInternalServerError
		}

		overallStatus := HealthCheck{
			Info: Info{
				ApplicationName: ctrl.Config.AppName,
				Version:         ctrl.Config.Version,
			},
			HealthCodes: HealthCodes{
				Application:     "OK",
				MongoConnection: http.StatusText(mongoStatus),
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(mongoStatus)

		if err := json.NewEncoder(w).Encode(overallStatus); err != nil {
			log.Error().Err(err).Msg("Failed to write response")
		}
	}
}

func (ctrl *Controller) GetLanguagesHandler(repo mgo.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var queryStrings models.QueryString
		err := schema.NewDecoder().Decode(&queryStrings, r.URL.Query())
		if err != nil {
			log.Error().Err(err).Msg("Failed to decode query string")
			http.Error(w, "Invalid query string", http.StatusBadRequest)
			return
		}

		rawQueryString := r.URL.Query().Encode()
		formattedQueryString := rawQueryString

		if len(queryStrings.Creators) > 0 {
			queryStrings.Creators = strings.Split(queryStrings.Creators[0], ",")
		}

		if len(queryStrings.Extensions) > 0 {
			queryStrings.Extensions = strings.Split(queryStrings.Extensions[0], ",")
		}

		if queryStrings.Size != nil {
			start := "size="
    		stop := strconv.Itoa(*queryStrings.Size)
			startIndex := strings.Index(formattedQueryString, start)
			stopIndex := strings.Index(formattedQueryString, stop) + len(stop)
			formattedQueryString = formattedQueryString[:startIndex] + formattedQueryString[stopIndex:]

			if *queryStrings.Size == -1 && queryStrings.Page != nil {
				log.Error().Err(err).Msg("Invalid query string given")
				http.Error(w, "Invalid query string", http.StatusBadRequest)
				return
			}

			if *queryStrings.Size != -1 && *queryStrings.Size < 1 {
				log.Error().Err(err).Msg("Invalid query string given")
				http.Error(w, "Invalid query string", http.StatusBadRequest)
				return
			}
		} else {
			queryStrings.Size = new(int)
			*queryStrings.Size = 10
		}

		if queryStrings.Page != nil {
			start := "page="
    		stop := strconv.Itoa(*queryStrings.Page)
			startIndex := strings.Index(formattedQueryString, start)
			stopIndex := strings.Index(formattedQueryString, stop) + len(stop)
			formattedQueryString = formattedQueryString[:startIndex] + formattedQueryString[stopIndex:]

			if *queryStrings.Page < 1 {
				log.Error().Err(err).Msg("Invalid query string given")
				http.Error(w, "Invalid query string", http.StatusBadRequest)
				return
			}
		} else {
			queryStrings.Page = new(int)
			*queryStrings.Page = 1
		}

		languages, total, filteredTotal, err := repo.GetLanguages(queryStrings)
		if err != nil {
			log.Error().Err(err).Msg("Failed to get languages")
			http.Error(w, "An error occurred processing this request", http.StatusInternalServerError)
			return
		}

		languagesResp := []models.LanguageResponse{}

		for _, language := range languages.Languages {
			languagesResp = append(languagesResp, models.LanguageResponse{
				Language: language,
				Links: []models.Links{
					models.Links{
						Rel:  "self",
						Href: fmt.Sprintf("/%s", language.Id),
					},
				},
			})
		}

		paginationLinks := []models.Links{}

		if len(formattedQueryString) > 0 {
			if string(formattedQueryString[0]) == "&" {
				formattedQueryString = formattedQueryString[1:]
			}

			formattedQueryString = formattedQueryString + "&"
		}

		if *queryStrings.Size != -1 {
			if total > (*queryStrings.Page * *queryStrings.Size) && len(languages.Languages) > 0 {
				if filteredTotal == 0 || filteredTotal > (*queryStrings.Page * *queryStrings.Size) {
					paginationLinks = append(paginationLinks, models.Links{
						Rel:  "next",
						Href: fmt.Sprintf("/?%spage=%d&size=%d", formattedQueryString, *queryStrings.Page+1, *queryStrings.Size),
					})
					if filteredTotal == 0 {
						if total % *queryStrings.Size != 0 {
							paginationLinks = append(paginationLinks, models.Links{
								Rel:  "last",
								Href: fmt.Sprintf("/?%spage=%d&size=%d", formattedQueryString, (total / *queryStrings.Size) + 1, *queryStrings.Size),
							})
						} else {
							paginationLinks = append(paginationLinks, models.Links{
								Rel:  "last",
								Href: fmt.Sprintf("/?%spage=%d&size=%d", formattedQueryString, (total / *queryStrings.Size), *queryStrings.Size),
							})
						}
						
						
					} else {
						if filteredTotal % *queryStrings.Size != 0 {
							paginationLinks = append(paginationLinks, models.Links{
								Rel:  "last",
								Href: fmt.Sprintf("/?%spage=%d&size=%d", formattedQueryString, (filteredTotal / *queryStrings.Size) + 1, *queryStrings.Size),
							})
						} else {
							paginationLinks = append(paginationLinks, models.Links{
								Rel:  "last",
								Href: fmt.Sprintf("/?%spage=%d&size=%d", formattedQueryString, (filteredTotal / *queryStrings.Size), *queryStrings.Size),
							})
						}
					}
				}
			}

			if *queryStrings.Page > 1 {
				paginationLinks = append(paginationLinks, models.Links{
					Rel:  "prev",
					Href: fmt.Sprintf("/?%spage=%d&size=%d", formattedQueryString, *queryStrings.Page-1, *queryStrings.Size),
				})

				paginationLinks = append(paginationLinks, models.Links{
					Rel:  "first",
					Href: fmt.Sprintf("/?%spage=%d&size=%d", formattedQueryString, 1, *queryStrings.Size),
				})
			}
		}

		resp := models.LanguagesResponse{
			Languages: languagesResp,
			Links:     paginationLinks,
			Total:     total,
			FilteredTotal: filteredTotal,
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Error().Err(err).Msg("Failed to write response")
		}
	}
}

func (ctrl *Controller) GetLanguageHandler(repo mgo.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]

		output, err := repo.GetLanguage(id)
		if err != nil {
			if errors.Is(err, models.ErrInvalidId) {
				http.Error(w, "The given id is not a valid id", http.StatusBadRequest)
				return
			}

			if errors.Is(err, models.ErrNotFound) {
				http.Error(w, "No language found with that id", http.StatusNotFound)
				return
			}

			log.Error().Err(err).Msg("Failed to get language")
			http.Error(w, "An error occurred processing this request", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(output); err != nil {
			log.Error().Err(err).Msg("Failed to write response")
		}
	}
}

func (ctrl *Controller) CreateLanguageHandler(repo mgo.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var language = models.Language{}

		err := json.NewDecoder(r.Body).Decode(&language)
		if err != nil {
			log.Error().Err(err).Msg("Failed to decode request body")
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		id, err := repo.PostLanguage(&language)
		if err != nil {
			log.Error().Err(err).Msg("Failed to create language")
			http.Error(w, "An error occurred processing this request", http.StatusInternalServerError)
			return
		}

		w.Header().Add("Location", "/"+url.PathEscape(id))
		w.WriteHeader(http.StatusCreated)
	}
}

func (ctrl *Controller) UpsertLanguageHandler(repo mgo.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]

		var language = models.Language{}
		err := json.NewDecoder(r.Body).Decode(&language)
		if err != nil {
			log.Error().Err(err).Msg("Failed to decode request body")
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		isUpserted, err := repo.PutLanguage(id, &language)
		if err != nil {
			if errors.Is(err, models.ErrInvalidId) {
				http.Error(w, "The given id is not a valid id", http.StatusBadRequest)
				return
			}

			log.Error().Err(err).Msg("Failed to upsert language")
			http.Error(w, "An error occurred processing this request", http.StatusInternalServerError)
			return
		}

		w.Header().Add("Location", "/"+url.PathEscape(id))

		if isUpserted {
			w.WriteHeader(http.StatusCreated)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}
}

func (ctrl *Controller) UpdateLanguageHandler(repo mgo.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]

		var update models.Language

		err := json.NewDecoder(r.Body).Decode(&update)
		if err != nil {
			log.Error().Err(err).Msg("Failed to decode request body")
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if len(update.Creators) > 0 {
			update.Creators = strings.Split(update.Creators[0], ",")
		}

		if len(update.Extensions) > 0 {
			update.Extensions = strings.Split(update.Extensions[0], ",")
		}

		err = repo.PatchLanguage(id, update)
		if err != nil {
			if errors.Is(err, models.ErrInvalidId) {
				http.Error(w, "The given id is not a valid id", http.StatusBadRequest)
				return
			}

			if errors.Is(err, models.ErrNotFound) {
				http.Error(w, "No language with that id found to update", http.StatusNotFound)
				return
			}

			log.Error().Err(err).Msg("Failed to update language")
			http.Error(w, "An error occurred processing this request", http.StatusInternalServerError)
			return
		}

		w.Header().Add("Location", "/"+url.PathEscape(id))
		w.WriteHeader(http.StatusOK)
	}
}

func (ctrl *Controller) DeleteLanguageHandler(repo mgo.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]

		err := repo.DeleteLanguage(id)
		if err != nil {
			if errors.Is(err, models.ErrInvalidId) {
				http.Error(w, "The given id is not a valid id", http.StatusBadRequest)
				return
			}

			if errors.Is(err, models.ErrNotFound) {
				http.Error(w, "No language with that id found to delete", http.StatusNotFound)
				return
			}

			log.Error().Err(err).Msg("Failed to delete language")
			http.Error(w, "An error occurred processing this request", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func (ctrl *Controller) NotFoundPageHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "You have accessed an invalid URL", http.StatusNotFound)
}
