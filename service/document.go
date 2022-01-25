package service

import (
	"precisely/model"
)

var (
	DocumentService documentServiceInterface = &documentService{}
)

type documentService struct{}

type documentServiceInterface interface {
	Get(int64) (*model.Document, error)
	Create(model.Document) (*model.Document, error)
	Update(model.Document) (*model.Document, error)
	Delete(int64) error
	GetAll() ([]*model.Document, error)
}

func (s *documentService) Create(newDocument model.Document) (*model.Document, error) {
	if err := newDocument.Validate(); err != nil {
		return nil, err
	}
	return model.DocumentRepository.Create(newDocument)
}

func (s *documentService) Update(inputDocument model.Document) (*model.Document, error) {
	if err := inputDocument.Validate(); err != nil {
		return nil, err
	}
	_, err := model.DocumentRepository.Get(inputDocument.ID)
	if err != nil {
		return nil, err
	}
	return model.DocumentRepository.Update(inputDocument)
}

func (s *documentService) Delete(id int64) error {
	_, err := model.DocumentRepository.Get(id)
	if err != nil {
		return err
	}
	return model.DocumentRepository.Delete(id)
}

func (s *documentService) Get(id int64) (*model.Document, error) {
	return model.DocumentRepository.Get(id)
}

func (s *documentService) GetAll() ([]*model.Document, error) {
	return model.DocumentRepository.GetAll()
}
