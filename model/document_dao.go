package model

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
)

var (
	DocumentRepository documentRepositoryInterface = &documentRepository{}
)

type documentRepositoryInterface interface {
	Get(int64) (*Document, error)
	Create(Document) (*Document, error)
	Update(Document) (*Document, error)
	Delete(int64) error
	GetAll() ([]*Document, error)
	Init(string, string, string, string, string, string) *sql.DB
}

type documentRepository struct {
	db *sql.DB
}

func NewDocumentRepository(db *sql.DB) documentRepositoryInterface {
	return &documentRepository{db: db}
}

func (r *documentRepository) Init(driver, username, password, port, host, database string) *sql.DB {
	var err error
	psqlInfo := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		username, password, host, port, database)
	r.db, err = sql.Open(driver, psqlInfo)
	if err != nil {
		log.Fatalf("failed opening connection to mysql: %v", err)
	}
	err = r.db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	return r.db
}

func (r *documentRepository) Get(id int64) (*Document, error) {
	stmt, err := r.db.Prepare("SELECT * FROM documents WHERE id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var doc Document
	result := stmt.QueryRow(id)
	if err := result.Scan(&doc.ID, &doc.Title, &doc.Content, &doc.Signee); err != nil {
		return nil, err
	}
	return &doc, nil
}

func (r *documentRepository) Create(newDoc Document) (*Document, error) {
	stmt, err := r.db.Prepare("INSERT INTO documents(title, content, signee) VALUES(?, ?, ?);")
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	contentJson, _ := json.Marshal(newDoc.Content)
	insertResult, err := stmt.Exec(newDoc.Title, contentJson, newDoc.Signee)
	if err != nil {
		return nil, err
	}
	id, err := insertResult.LastInsertId()
	if err != nil {
		return nil, err
	}
	newDoc.ID = id

	return &newDoc, nil
}

func (r *documentRepository) Update(upDoc Document) (*Document, error) {

	stmt, err := r.db.Prepare("UPDATE documents SET title = ?, content = ?, signee = ? WHERE id = ?")
	if err != nil {
		return nil, err
	}

	defer stmt.Close()
	contentJson, _ := json.Marshal(upDoc.Content)
	result, err := stmt.Exec(
		upDoc.Title,
		contentJson,
		upDoc.Signee,
		upDoc.ID)
	if err != nil {
		return nil, err
	}
	if _, err := result.RowsAffected(); err != nil {
		return nil, err
	}
	return &upDoc, nil
}

func (r *documentRepository) GetAll() ([]*Document, error) {
	stmt, err := r.db.Prepare("SELECT * FROM documents")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]*Document, 0)

	for rows.Next() {
		var doc Document
		if getError := rows.Scan(&doc.ID, &doc.Title, &doc.Content, &doc.Signee); getError != nil {
			return nil, getError
		}
		results = append(results, &doc)
	}
	return results, nil
}

func (r *documentRepository) Delete(id int64) error {
	stmt, err := r.db.Prepare("DELETE FROM documents WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(id); err != nil {
		return err
	}
	return nil
}
