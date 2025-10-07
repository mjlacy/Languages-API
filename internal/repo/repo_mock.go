package repo

import (
	"languages-api/internal/models"
)

type MockRepo struct {
	languages  models.Languages
	language   models.Language
	id         string
	isUpserted bool
	Err        error
}

func (m *MockRepo) Ping() error {
	return m.Err
}

func (m *MockRepo) GetLanguages(_ models.Language) (languages models.Languages, err error) {
	return m.languages, m.Err
}

func (m *MockRepo) GetLanguage(_ string) (language models.Language, err error) {
	return m.language, m.Err
}

func (m *MockRepo) PostLanguage(_ models.Language) (string, error) {
	return m.id, m.Err
}

func (m *MockRepo) PutLanguage(_ string, _ models.Language) (bool, error) {
	return m.isUpserted, m.Err
}

func (m *MockRepo) PatchLanguage(_ string, _ models.Language) (err error) {
	return m.Err
}

func (m *MockRepo) DeleteLanguage(_ string) (err error) {
	return m.Err
}

func (m *MockRepo) Close() error {
	return m.Err
}
