package repo

import (
	"languages-api/internal/models"
)

type MockRepo struct {
	languages  models.Languages
	language   *models.Language
	id         string
	isUpserted bool
	Err        error
}

func (m *MockRepo) Ping() error {
	return m.Err
}

func (m *MockRepo) GetLanguages(language models.Language) (languages models.Languages, err error) {
	return m.languages, m.Err
}

func (m *MockRepo) GetLanguage(id string) (language *models.Language, err error) {
	return m.language, m.Err
}

func (m *MockRepo) PostLanguage(language *models.Language) (string, error) {
	return m.id, m.Err
}

func (m *MockRepo) PutLanguage(id string, language *models.Language) (bool, error) {
	return m.isUpserted, m.Err
}

func (m *MockRepo) PatchLanguage(id string, update models.Language) (err error) {
	return m.Err
}

func (m *MockRepo) DeleteLanguage(id string) (err error) {
	return m.Err
}

func (m *MockRepo) Close() error {
	return m.Err
}
