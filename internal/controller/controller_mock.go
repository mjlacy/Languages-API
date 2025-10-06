package controller

import (
	"languages-api/internal/models"

	"net/http"
)

type mockResponseWriter struct {
	header     http.Header
	message    string
	num        int
	statusCode int
	err        error
}

func (w *mockResponseWriter) Header() http.Header {
	return w.header
}

func (w *mockResponseWriter) Write(data []byte) (int, error) {
	w.message = string(data)
	return w.num, w.err
}

func (w *mockResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	return
}

type mockRepository struct {
	err        error
	errs       []error
	id         string
	isUpserted bool
	ls         models.Languages
	l          models.Language
}

func (r mockRepository) Ping() error {
	return r.err
}

func (r mockRepository) Close() error {
	return r.err
}

func (r mockRepository) GetLanguages(_ models.Language) (models.Languages, []error) {
	return r.ls, r.errs
}

func (r mockRepository) GetLanguage(_ string) (models.Language, error) {
	return r.l, r.err
}

func (r mockRepository) PostLanguage(_ models.Language) (string, error) {
	return r.id, r.err
}

func (r mockRepository) PutLanguage(_ string, _ models.Language) (isUpserted bool, err error) {
	return r.isUpserted, r.err
}

func (r mockRepository) PatchLanguage(_ string, _ models.Language) (err error) {
	return r.err
}

func (r mockRepository) DeleteLanguage(_ string) (err error) {
	return r.err
}
