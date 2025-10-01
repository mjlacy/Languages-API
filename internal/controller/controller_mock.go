package controller

import "languages-api/internal/models"

type mockRepository struct {
	err        error
	errs       []error
	id         string
	isUpserted bool
	l          models.Languages
}

func (r mockRepository) Ping() error {
	return r.err
}

func (r mockRepository) Close() error {
	return r.err
}

func (r mockRepository) GetLanguages(_ models.Language) (models.Languages, []error) {
	return r.l, r.errs
}

func (r mockRepository) GetLanguage(_ string) (models.Language, error) {
	return r.l.Languages[0], r.err
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
