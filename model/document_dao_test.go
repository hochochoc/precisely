package model

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"log"
	"reflect"
	"testing"
)

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	client, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("")
	}
	return client, mock
}

func TestDocumentRepository_Get(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	r := NewDocumentRepository(db)

	tests := []struct {
		name    string
		r       documentRepositoryInterface
		id      int64
		mock    func()
		want    *Document
		wantErr bool
	}{
		{
			name: "Ok",
			r:    r,
			id:   1,
			mock: func() {
				contentBytes, _ := json.Marshal(Content{Header: "header", Data: "data"})
				rows := sqlmock.NewRows([]string{"id", "title", "content", "signee"}).
					AddRow(1, "Document 1", contentBytes, "Signee 1")
				mock.ExpectPrepare("SELECT (.+) FROM documents").ExpectQuery().WithArgs(1).
					WillReturnRows(rows)
			},
			want: &Document{
				ID:      1,
				Title:   "Document 1",
				Content: Content{Header: "header", Data: "data"},
				Signee:  "Signee 1",
			},
		},
		{
			name: "Document Not Found",
			r:    r,
			id:   1,
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "content", "signee"})
				mock.ExpectPrepare("SELECT (.+) FROM documents").ExpectQuery().WithArgs(1).
					WillReturnRows(rows)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := tt.r.Get(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error new = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDocumentRepository_Create(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()
	r := NewDocumentRepository(db)
	tests := []struct {
		name    string
		r       documentRepositoryInterface
		doc     Document
		mock    func()
		want    *Document
		wantErr bool
	}{
		{
			name: "Ok",
			r:    r,
			doc: Document{
				Title: "title",
				Content: Content{
					Header: "header",
					Data:   "data",
				},
				Signee: "signee",
			},
			mock: func() {
				contentBytes, _ := json.Marshal(Content{Header: "header", Data: "data"})
				mock.ExpectPrepare("INSERT INTO documents").ExpectExec().
					WithArgs("title", contentBytes, "signee").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},

			want: &Document{
				ID:    1,
				Title: "title",
				Content: Content{
					Header: "header",
					Data:   "data",
				},
				Signee: "signee",
			},
		},
		{
			name: "Empty content",
			r:    r,
			doc: Document{
				Title:  "title",
				Signee: "signee",
			},
			mock: func() {
				contentBytes, _ := json.Marshal(Content{Header: "header", Data: "data"})
				mock.ExpectPrepare("INSERT INTO messages").ExpectExec().WithArgs("title", contentBytes, "signee").WillReturnError(errors.New("empty content"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := tt.r.Create(tt.doc)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDocumentRepository_Update(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	r := NewDocumentRepository(db)

	tests := []struct {
		name    string
		r       documentRepositoryInterface
		doc     Document
		mock    func()
		want    *Document
		wantErr bool
	}{
		{
			name: "Ok",
			r:    r,
			doc: Document{
				ID:    1,
				Title: "title",
				Content: Content{
					Header: "header",
					Data:   "data",
				},
				Signee: "signee",
			},
			mock: func() {
				contentBytes, _ := json.Marshal(Content{Header: "header", Data: "data"})
				mock.ExpectPrepare("UPDATE documents").ExpectExec().
					WithArgs("title", contentBytes, "signee", 1).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			want: &Document{
				ID:    1,
				Title: "title",
				Content: Content{
					Header: "header",
					Data:   "data",
				},
				Signee: "signee",
			},
		},
		{
			name: "Invalid query id",
			r:    r,
			doc: Document{
				ID:    1,
				Title: "title",
				Content: Content{
					Header: "header",
					Data:   "data",
				},
				Signee: "signee",
			},
			mock: func() {
				contentBytes, _ := json.Marshal(Content{Header: "header", Data: "data"})
				mock.ExpectPrepare("UPDATE documents").ExpectExec().
					WithArgs("update title", contentBytes, "update signee", 0).
					WillReturnError(errors.New("invalid update id"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := tt.r.Update(tt.doc)
			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Update() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDocumentRepository_GetAll(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	r := NewDocumentRepository(db)

	tests := []struct {
		name    string
		r       documentRepositoryInterface
		mock    func()
		want    []*Document
		wantErr bool
	}{
		{
			name: "Ok",
			r:    r,
			mock: func() {
				contentBytes, _ := json.Marshal(Content{
					Header: "header",
					Data:   "data",
				})
				rows := sqlmock.NewRows([]string{"id", "title", "content", "signee"}).
					AddRow(1, "first title", contentBytes, "first signee").
					AddRow(2, "second title", contentBytes, "second signee")
				mock.ExpectPrepare("SELECT (.+) FROM documents").ExpectQuery().WillReturnRows(rows)
			},
			want: []*Document{
				{
					ID:    1,
					Title: "first title",
					Content: Content{
						Header: "header",
						Data:   "data",
					},
					Signee: "first signee",
				},
				{
					ID:    2,
					Title: "second title",
					Content: Content{
						Header: "header",
						Data:   "data",
					},
					Signee: "second signee",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := tt.r.GetAll()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAll() error new = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDocumentRepository_Delete(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	r := NewDocumentRepository(db)

	tests := []struct {
		name    string
		r       documentRepositoryInterface
		id      int64
		mock    func()
		want    *Document
		wantErr bool
	}{
		{
			name: "Ok",
			r:    r,
			id:   1,
			mock: func() {
				mock.ExpectPrepare("DELETE FROM documents").ExpectExec().
					WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: false,
		},
		{
			name: "Not Found Id",
			r:    r,
			id:   1,
			mock: func() {
				mock.ExpectPrepare("DELETE FROM documents").ExpectExec().
					WithArgs(100).WillReturnResult(sqlmock.NewResult(0, 0))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			err := tt.r.Delete(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error new = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
