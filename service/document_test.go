package service

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"precisely/model"
	"reflect"
	"testing"
)

var (
	getMessageDAO    func(id int64) (*model.Document, error)
	createMessageDAO func(doc model.Document) (*model.Document, error)
	updateMessageDAO func(doc model.Document) (*model.Document, error)
	deleteMessageDAO func(id int64) error
	getAllMessageDAO func() ([]*model.Document, error)
)

type dBMock struct{}

func (m *dBMock) Get(id int64) (*model.Document, error) {
	return getMessageDAO(id)
}

func (m *dBMock) GetAll() ([]*model.Document, error) {
	return getAllMessageDAO()
}

func (m *dBMock) Delete(id int64) error {
	return deleteMessageDAO(id)
}

func (m *dBMock) Create(doc model.Document) (*model.Document, error) {
	return createMessageDAO(doc)
}

func (m *dBMock) Update(doc model.Document) (*model.Document, error) {
	return updateMessageDAO(doc)
}

func (m *dBMock) Init(string, string, string, string, string, string) *sql.DB {
	return nil
}

func TestDocumentService_Get_Success(t *testing.T) {
	model.DocumentRepository = &dBMock{}
	mockData := &model.Document{
		ID:    1,
		Title: "title",
		Content: model.Content{
			Header: "header",
			Data:   "data",
		},
		Signee: "signee",
	}
	getMessageDAO = func(id int64) (*model.Document, error) {
		return mockData, nil
	}
	doc, err := DocumentService.Get(1)
	assert.NotNil(t, doc)
	assert.Nil(t, err)
	assert.True(t, reflect.DeepEqual(doc, mockData))
}

func TestDocumentService_Get_NotFound(t *testing.T) {
	model.DocumentRepository = &dBMock{}
	getMessageDAO = func(id int64) (*model.Document, error) {
		return nil, sql.ErrNoRows
	}
	doc, err := DocumentService.Get(1)
	assert.Nil(t, doc)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "sql: no rows in result set")
}

func TestDocumentService_Create_Success(t *testing.T) {
	model.DocumentRepository = &dBMock{}
	mockData := &model.Document{
		ID:    1,
		Title: "title",
		Content: model.Content{
			Header: "header",
			Data:   "data",
		},
		Signee: "signee",
	}
	createMessageDAO = func(doc model.Document) (*model.Document, error) {
		return mockData, nil
	}
	doc := model.Document{
		Title: "title",
		Content: model.Content{
			Header: "header",
			Data:   "data",
		},
		Signee: "signee",
	}
	savedDoc, err := DocumentService.Create(doc)
	assert.Nil(t, err)
	assert.NotNil(t, savedDoc)
	assert.EqualValues(t, doc.Title, savedDoc.Title)
	assert.EqualValues(t, doc.Content, savedDoc.Content)
	assert.EqualValues(t, doc.Signee, savedDoc.Signee)
}

func TestDocumentService_Create_BadRequest(t *testing.T) {
	tests := []struct {
		doc        model.Document
		errMessage string
	}{
		{
			doc: model.Document{
				Signee: "signee",
			},
			errMessage: model.TitleInvalidValue.Error(),
		},
		{
			doc: model.Document{
				Title: "title",
			},
			errMessage: model.SigneeInvalidValue.Error(),
		},
	}
	for _, tt := range tests {
		doc, err := DocumentService.Create(tt.doc)
		assert.Nil(t, doc)
		assert.NotNil(t, err)
		assert.EqualValues(t, tt.errMessage, err.Error())
	}
}

func TestDocumentService_Create_InternalFailure(t *testing.T) {
	model.DocumentRepository = &dBMock{}
	createMessageDAO = func(doc model.Document) (*model.Document, error) {
		return nil, &mysql.MySQLError{Number: 1062, Message: "Duplicate entry 'Title' for key 'title'"}
	}

	doc := model.Document{
		Title:  "title",
		Signee: "signee",
	}
	savedDoc, err := DocumentService.Create(doc)
	assert.Nil(t, savedDoc)
	assert.NotNil(t, err)
	assert.EqualValues(t, 1062, err.(*mysql.MySQLError).Number)
}

func TestDocumentService_Update_Success(t *testing.T) {
	model.DocumentRepository = &dBMock{}
	getMessageDAO = func(id int64) (*model.Document, error) {
		return &model.Document{
			ID:     1,
			Title:  "title",
			Signee: "signee",
		}, nil
	}
	updateMessageDAO = func(doc model.Document) (*model.Document, error) {
		return &model.Document{
			ID:     1,
			Title:  "update title",
			Signee: "update signee",
		}, nil
	}
	doc := model.Document{
		Title:  "update title",
		Signee: "update signee",
	}
	updatedDoc, err := DocumentService.Update(doc)
	assert.NotNil(t, updatedDoc)
	assert.Nil(t, err)
	assert.EqualValues(t, doc.Title, updatedDoc.Title)
	assert.EqualValues(t, doc.Signee, updatedDoc.Signee)
	assert.EqualValues(t, doc.Content, updatedDoc.Content)
}

func TestDocumentService_Delete_Success(t *testing.T) {
	model.DocumentRepository = &dBMock{}
	getMessageDAO = func(id int64) (*model.Document, error) {
		return &model.Document{
			ID:     1,
			Title:  "title",
			Signee: "signee",
		}, nil
	}

	deleteMessageDAO = func(id int64) error {
		return nil
	}
	err := DocumentService.Delete(1)
	assert.Nil(t, err)
}

func TestDocumentService_Delete_NotFound(t *testing.T) {
	model.DocumentRepository = &dBMock{}
	getMessageDAO = func(id int64) (*model.Document, error) {
		return nil, sql.ErrNoRows
	}
	err := DocumentService.Delete(1)
	assert.NotNil(t, err)
}

func TestDocumentService_GetAll_Success(t *testing.T) {
	model.DocumentRepository = &dBMock{}
	getAllMessageDAO = func() ([]*model.Document, error) {
		return []*model.Document{
			{
				ID:     1,
				Title:  "title1",
				Signee: "signee1",
			},
		}, nil
	}
	docs, err := DocumentService.GetAll()
	assert.Nil(t, err)
	assert.NotNil(t, docs)
	assert.EqualValues(t, len(docs), 1)
}
