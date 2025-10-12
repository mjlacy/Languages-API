package mariadb

import (
	"languages-api/internal/config"
	"languages-api/internal/models"

	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

// Client is for wrappers of *sql.DB
type Client interface {
	Ping() error
	Disconnect() error
	Find(filter interface{}) (languages models.Languages, errors []error)
	FindOne(id string) (language models.Language, err error)
	InsertOne(document interface{}) (insertedId string, err error)
	ReplaceOne(id string, document interface{}) (isUpserted bool, err error)
	UpdateOne(id string, update interface{}) (err error)
	DeleteOne(id string) (err error)
}

// MariaClient implements the Client interface
type MariaClient struct {
	*sql.DB
	DatabaseName string
	TableName    string
}

//func New(c Config) (repo MariaRepository, err error) {
//	db, err := sql.Open("mysql", c.URL)
//	if err != nil {
//		log.Printf("Error establishing connection to MariaDB: %v", err)
//		return
//	}
//
//	err = db.Ping()
//	if err != nil {
//		log.Printf("Error while pinging Maria: %v", err)
//	}
//
//	repo = MariaRepository{DB: db, DatabaseName: c.Database, TableName: c.Table}
//
//	return
//}

func (r MariaClient) Ping() error {
	return r.DB.Ping()
}

func (r MariaClient) Disconnect() error {
	return r.DB.Close()
}

func (r MariaClient) Find(filter interface{}) (languages models.Languages, errs []error) {
	language := filter.(models.Language)

	conditions := ""
	var values []interface{}

	if language.Name != "" {
		conditions += " AND name=?"
		values = append(values, language.Name)
	}

	if len(language.Creators) > 0 {
		conditions += " AND creators=?"
		values = append(values, strings.Join(language.Creators, ","))
	}

	if len(language.Extensions) > 0 {
		conditions += " AND extensions=?"
		values = append(values, strings.Join(language.Extensions, ","))
	}

	if language.FirstAppeared != nil {
		conditions += " AND firstAppeared=?"
		values = append(values, language.FirstAppeared)
	}

	if language.Year != 0 {
		conditions += " AND year=?"
		values = append(values, language.Year)
	}

	if language.Wiki != "" {
		conditions += " AND wiki=?"
		values = append(values, language.Wiki)
	}

	stmt, err := r.DB.Prepare(fmt.Sprintf("SELECT * FROM %s WHERE 1=1%v;", r.TableName, conditions)) // prevents something like ?origin=' OR color='pale
	if err != nil {
		errs = append(errs, err)
	}

	results, err := stmt.Query(values...)
	if err != nil {
		errs = append(errs, err)
	}

	defer func() {
		if resultsCloseErr := results.Close(); resultsCloseErr != nil && err == nil {
			err = resultsCloseErr
		}

		if stmtCloseErr := stmt.Close(); stmtCloseErr != nil && err == nil {
			err = stmtCloseErr
		}
	}()

	languages.Languages = []models.Language{}

	for results.Next() {
		var lang models.Language
		creators := ""
		extensions := ""

		err = results.Scan(&lang.Id, &lang.Name, &creators, &extensions, &lang.FirstAppeared, &lang.Year, &lang.Wiki)
		if err != nil {
			errs = append(errs, err)
		}
		lang.Creators = strings.Split(creators, ",")
		lang.Extensions = strings.Split(extensions, ",")
		languages.Languages = append(languages.Languages, lang)
	}

	return
}

func (r MariaClient) FindOne(id string) (language models.Language, err error) {
	idNum, err := strconv.Atoi(id)
	if err != nil || idNum <= 0 {
		err = models.ErrInvalidId
		return
	}

	stmt, err := r.DB.Prepare(fmt.Sprintf("SELECT * FROM %s WHERE _id=?;", r.TableName)) // prevents something like /0 OR name='Michelob' if id didn't get turned into a number
	if err != nil {
		return
	}

	defer func() {
		if closeErr := stmt.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	creators := ""
	extensions := ""

	err = stmt.QueryRow(idNum).Scan(&language.Id, &language.Name, &creators, &extensions, &language.FirstAppeared, &language.Year, &language.Wiki)
	if errors.Is(err, sql.ErrNoRows) {
		err = models.ErrNotFound
	}

	language.Creators = strings.Split(creators, ",")
	language.Extensions = strings.Split(extensions, ",")

	return
}

func (r MariaClient) InsertOne(document interface{}) (insertedId string, err error) {
	language := document.(models.Language)

	creators := strings.Join(language.Creators, ",")
	extensions := strings.Join(language.Extensions, ",")

	stmt, err := r.DB.Prepare(fmt.Sprintf("INSERT INTO %s (name, creators, extensions, firstAppeared, year, wiki) VALUES (?, ?, ?, ?, ?, ?);", r.TableName))
	if err != nil {
		return
	}

	defer func() {
		if closeErr := stmt.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	result, err := stmt.Exec(language.Name, creators, extensions, language.FirstAppeared, language.Year, language.Wiki)
	if err != nil {
		return
	}

	idStr, err := result.LastInsertId()
	if err != nil {
		return
	}

	insertedId = strconv.FormatInt(idStr, 10)
	return
}

func (r MariaClient) ReplaceOne(id string, document interface{}) (isUpserted bool, err error) {
	idNum, err := strconv.Atoi(id)
	if err != nil || idNum <= 0 {
		err = models.ErrInvalidId
		return
	}

	stmt, err := r.DB.Prepare(fmt.Sprintf("REPLACE INTO %s (_id, name, creators, extensions, firstAppeared, year, wiki) VALUES (?, ?, ?, ?, ?, ?, ?);", r.TableName))
	if err != nil {
		return
	}

	defer func() {
		if closeErr := stmt.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	language := document.(models.Language)
	creators := strings.Join(language.Creators, ",")
	extensions := strings.Join(language.Extensions, ",")

	result, err := stmt.Exec(idNum, language.Name, creators, extensions, language.FirstAppeared, language.Year, language.Wiki)
	if err != nil {
		return
	}

	numRowsAffected, err := result.RowsAffected()
	if err != nil {
		return
	}

	return numRowsAffected == 1, err
}

func (r MariaClient) UpdateOne(id string, update interface{}) (err error) {
	idNum, err := strconv.Atoi(id)
	if err != nil || idNum <= 0 {
		err = models.ErrInvalidId
		return
	}

	_, err = r.FindOne(id)
	if err != nil {
		return
	}

	lang := update.(models.Language)
	creators := strings.Join(lang.Creators, ",")
	extensions := strings.Join(lang.Extensions, ",")

	keys := ""

	v := reflect.ValueOf(lang)

	var values []interface{}

	for i := 0; i < v.NumField(); i++ {
		if !v.Field(i).IsZero() && strings.ToLower(v.Type().Field(i).Name) != "id" {
			keys += fmt.Sprintf(" %v = ?,", strings.ToLower(v.Type().Field(i).Name))
			if strings.ToLower(v.Type().Field(i).Name) == "creators" {
				values = append(values, creators)
			} else if strings.ToLower(v.Type().Field(i).Name) == "extensions" {
				values = append(values, extensions)
			} else {
				values = append(values, v.Field(i).Interface())
			}
		}
	}

	keys = strings.TrimSuffix(keys, ",")

	if len(values) == 0 {
		return models.ErrInvalidRequestBody
	}

	stmt, err := r.DB.Prepare(fmt.Sprintf("UPDATE %s SET%v WHERE _id=%d;", r.TableName, keys, idNum))
	if err != nil {
		return
	}

	defer func() {
		if closeErr := stmt.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	_, err = stmt.Exec(values...)

	return
}

func (r MariaClient) DeleteOne(id string) (err error) {
	idNum, err := strconv.Atoi(id)
	if err != nil || idNum <= 0 {
		err = models.ErrInvalidId
		return
	}

	stmt, err := r.DB.Prepare(fmt.Sprintf("DELETE FROM %s WHERE _id=?;", r.TableName))
	if err != nil {
		return
	}

	defer func() {
		if closeErr := stmt.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	result, err := stmt.Exec(idNum)
	if err != nil {
		return
	}

	numRowsAffected, err := result.RowsAffected()
	if err != nil {
		return
	} else if numRowsAffected == 0 {
		err = models.ErrNotFound
	}

	return
}

// Connector specifies the methods needed to connect to mongo
type Connector interface {
	Connect(cfg config.Config) (Client, error)
}

// MariaConnector implements the Connector interface
type MariaConnector struct{}

// Connect establishes the connection to maria
func (mc MariaConnector) Connect(cfg config.Config) (Client, error) {
	db, err := sql.Open("mysql", cfg.DBURL)

	return &MariaClient{DB: db, DatabaseName: cfg.Database, TableName: cfg.Table}, err
}
