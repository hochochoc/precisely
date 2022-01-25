package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"precisely/model"
	"precisely/service"
	"precisely/utils"
	"testing"
)

var (
	getMessageService    func(id int64) (*model.Document, error)
	createMessageService func(doc model.Document) (*model.Document, error)
	updateMessageService func(doc model.Document) (*model.Document, error)
	deleteMessageService func(id int64) error
	getAllMessageService func() ([]*model.Document, error)
)

type serviceMock struct{}

func (m *serviceMock) Get(id int64) (*model.Document, error) {
	return getMessageService(id)
}

func (m *serviceMock) GetAll() ([]*model.Document, error) {
	return getAllMessageService()
}

func (m *serviceMock) Delete(id int64) error {
	return deleteMessageService(id)
}

func (m *serviceMock) Create(doc model.Document) (*model.Document, error) {
	return createMessageService(doc)
}

func (m *serviceMock) Update(doc model.Document) (*model.Document, error) {
	return updateMessageService(doc)
}

func TestGetAllHandler(t *testing.T) {
	service.DocumentService = &serviceMock{}
	testData := []*model.Document{
		{
			ID:     1,
			Title:  "title",
			Signee: "signee",
		},
	}
	getAllMessageService = func() ([]*model.Document, error) {
		return testData, nil
	}
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/documents", nil)
	handler := http.HandlerFunc(GetAllHandler)
	handler.ServeHTTP(rr, req)

	var res utils.HttpResponse
	err := json.Unmarshal(rr.Body.Bytes(), &res)
	if err != nil {
		t.Error(err)
	}
	assert.Nil(t, err)
	assert.EqualValues(t, 1, len(res.Data.([]interface{})))
	assert.EqualValues(t, http.StatusOK, res.Code)
}

func TestCreateHandler_Success(t *testing.T) {
	service.DocumentService = &serviceMock{}
	createMessageService = func(doc model.Document) (*model.Document, error) {
		return &model.Document{
			ID:     1,
			Title:  "title",
			Signee: "signee",
		}, nil
	}
	jsonBody := `{"title": "title", "signee": "signee"}`

	req, err := http.NewRequest(http.MethodPost, "/documents", bytes.NewBufferString(jsonBody))
	if err != nil {
		t.Error(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateHandler)
	handler.ServeHTTP(rr, req)
	var res utils.HttpResponse
	err = json.Unmarshal(rr.Body.Bytes(), &res)
	if err != nil {
		t.Error(err)
	}
	assert.Nil(t, err)
	assert.EqualValues(t, http.StatusCreated, res.Code)
	assert.EqualValues(t, 1, res.Data.(map[string]interface{})["id"])
	assert.EqualValues(t, "title", res.Data.(map[string]interface{})["title"])
	assert.EqualValues(t, "signee", res.Data.(map[string]interface{})["signee"])
}

func TestCreateHandler_UnprocessableEntity(t *testing.T) {
	service.DocumentService = &serviceMock{}
	createMessageService = func(doc model.Document) (*model.Document, error) {
		return nil, model.TitleInvalidValue
	}
	jsonBody := `{"title": "", "signee": "signee"}`

	req, err := http.NewRequest(http.MethodPost, "/documents", bytes.NewBufferString(jsonBody))
	if err != nil {
		t.Error(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateHandler)
	handler.ServeHTTP(rr, req)

	var res utils.HttpResponse
	err = json.Unmarshal(rr.Body.Bytes(), &res)
	if err != nil {
		t.Error(err)
	}
	assert.Nil(t, err)
	assert.EqualValues(t, http.StatusUnprocessableEntity, res.Code)
}

func TestDeleteHandler_Success(t *testing.T) {
	service.DocumentService = &serviceMock{}
	deleteMessageService = func(id int64) error {
		return nil
	}

	req, err := http.NewRequest(http.MethodDelete, "/documents", nil)
	if err != nil {
		t.Error(err)
	}
	req = mux.SetURLVars(req, map[string]string{
		"id": "1",
	})
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(DeleteHandler)
	handler.ServeHTTP(rr, req)

	var res utils.HttpResponse
	err = json.Unmarshal(rr.Body.Bytes(), &res)
	if err != nil {
		t.Error(err)
	}
	assert.EqualValues(t, http.StatusOK, res.Code)
}

func TestDeleteHandler_NotFound(t *testing.T) {
	service.DocumentService = &serviceMock{}
	deleteMessageService = func(id int64) error {
		return sql.ErrNoRows
	}
	req, err := http.NewRequest(http.MethodDelete, "/documents", nil)
	if err != nil {
		t.Error(err)
	}
	req = mux.SetURLVars(req, map[string]string{
		"id": "1",
	})
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(DeleteHandler)
	handler.ServeHTTP(rr, req)

	var res utils.HttpResponse
	err = json.Unmarshal(rr.Body.Bytes(), &res)
	if err != nil {
		t.Error(err)
	}
	assert.EqualValues(t, http.StatusNotFound, res.Code)
	assert.EqualValues(t, sql.ErrNoRows.Error(), res.Error)
}

func TestUpdateHandler_Success(t *testing.T) {
	service.DocumentService = &serviceMock{}
	updateMessageService = func(doc model.Document) (*model.Document, error) {
		return &model.Document{
			ID:     1,
			Title:  "updated title",
			Signee: "updated signee",
			Content: model.Content{
				Header: "updated header",
				Data:   "updated data",
			},
		}, nil
	}

	jsonBody := `
		{
			"title": "updated title",
			"content": {
				"header": "updated header",
				"data": "updated data"
			},
			"signee": "updated signee"
		}
`

	req, err := http.NewRequest(http.MethodPut, "/documents", bytes.NewBufferString(jsonBody))
	if err != nil {
		t.Error(err)
	}
	req = mux.SetURLVars(req, map[string]string{
		"id": "1",
	})
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(UpdateHandler)
	handler.ServeHTTP(rr, req)

	var res utils.HttpResponse
	err = json.Unmarshal(rr.Body.Bytes(), &res)
	if err != nil {
		t.Error(err)
	}
	doc := res.Data.(map[string]interface{})
	content := doc["content"].(map[string]interface{})

	assert.Nil(t, err)
	assert.EqualValues(t, http.StatusOK, res.Code)
	assert.EqualValues(t, "", res.Error)
	assert.EqualValues(t, "updated title", doc["title"])
	assert.EqualValues(t, "updated signee", doc["signee"])
	assert.EqualValues(t, "updated header", content["header"])
	assert.EqualValues(t, "updated data", content["data"])
}

func TestUpdateHandler_NotFound(t *testing.T) {
	service.DocumentService = &serviceMock{}
	updateMessageService = func(doc model.Document) (*model.Document, error) {
		return nil, sql.ErrNoRows
	}

	jsonBody := `
		{
			"title": "updated title",
			"content": {
				"header": "updated header",
				"data": "updated data"
			},
			"signee": "updated signee"
		}
`

	req, err := http.NewRequest(http.MethodPut, "/documents", bytes.NewBufferString(jsonBody))
	if err != nil {
		t.Error(err)
	}
	req = mux.SetURLVars(req, map[string]string{
		"id": "1",
	})
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(UpdateHandler)
	handler.ServeHTTP(rr, req)

	var res utils.HttpResponse
	err = json.Unmarshal(rr.Body.Bytes(), &res)
	if err != nil {
		t.Error(err)
	}

	assert.Nil(t, err)
	assert.EqualValues(t, http.StatusNotFound, res.Code)
}

func TestUpdateHandler_InvalidJsonBody(t *testing.T) {
	jsonBody := `
		{
			"title": 123,
			"content": {
				"header": "updated header",
				"data": "updated data"
			},
			"signee": "updated signee"
		}
`
	req, err := http.NewRequest(http.MethodPut, "/documents", bytes.NewBufferString(jsonBody))
	if err != nil {
		t.Error(err)
	}
	req = mux.SetURLVars(req, map[string]string{
		"id": "1",
	})
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(UpdateHandler)
	handler.ServeHTTP(rr, req)

	var res utils.HttpResponse
	err = json.Unmarshal(rr.Body.Bytes(), &res)
	if err != nil {
		t.Error(err)
	}

	assert.Nil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, res.Code)
}

func TestGetByIdHandler_Success(t *testing.T) {
	service.DocumentService = &serviceMock{}
	getMessageService = func(id int64) (*model.Document, error) {
		return &model.Document{
			ID:    1,
			Title: "title",
			Content: model.Content{
				Header: "header",
				Data:   "data",
			},
			Signee: "signee",
		}, nil
	}
	req, _ := http.NewRequest(http.MethodGet, "/messages", nil)
	req = mux.SetURLVars(req, map[string]string{
		"id": "1",
	})
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetByIdHandler)
	handler.ServeHTTP(rr, req)

	var res utils.HttpResponse
	err := json.Unmarshal(rr.Body.Bytes(), &res)
	if err != nil {
		t.Error(err)
	}
	doc := res.Data.(map[string]interface{})
	content := doc["content"].(map[string]interface{})

	assert.Nil(t, err)
	assert.EqualValues(t, http.StatusOK, res.Code)
	assert.EqualValues(t, "", res.Error)
	assert.EqualValues(t, "title", doc["title"])
	assert.EqualValues(t, "signee", doc["signee"])
	assert.EqualValues(t, "header", content["header"])
	assert.EqualValues(t, "data", content["data"])
}

func TestGetByIdHandler_NotFound(t *testing.T) {
	service.DocumentService = &serviceMock{}
	getMessageService = func(id int64) (*model.Document, error) {
		return nil, sql.ErrNoRows
	}
	req, _ := http.NewRequest(http.MethodGet, "/messages", nil)
	req = mux.SetURLVars(req, map[string]string{
		"id": "1",
	})
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetByIdHandler)
	handler.ServeHTTP(rr, req)

	var res utils.HttpResponse
	err := json.Unmarshal(rr.Body.Bytes(), &res)
	if err != nil {
		t.Error(err)
	}

	assert.Nil(t, err)
	assert.EqualValues(t, http.StatusNotFound, res.Code)
}
