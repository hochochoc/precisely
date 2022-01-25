package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
	"precisely/model"
	"precisely/service"
	"precisely/utils"
	"strconv"
)

func CreateHandler(w http.ResponseWriter, r *http.Request) {
	var newDocument model.Document
	err := json.NewDecoder(r.Body).Decode(&newDocument)
	if err != nil {
		utils.JsonRespond(w, false, http.StatusBadRequest, err, nil)
		return
	}
	document, err := service.DocumentService.Create(newDocument)
	if err != nil {
		if errors.Is(err, model.TitleInvalidValue) || errors.Is(err, model.SigneeInvalidValue) {
			utils.JsonRespond(w, false, http.StatusUnprocessableEntity, err, nil)
			return
		}
		utils.JsonRespond(w, false, http.StatusInternalServerError, err, nil)
		return
	}
	utils.JsonRespond(w, true, http.StatusCreated, err, document)
	return
}

func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	var updatedDocument model.Document
	err := json.NewDecoder(r.Body).Decode(&updatedDocument)
	if err != nil {
		utils.JsonRespond(w, false, http.StatusBadRequest, err, nil)
		return
	}
	idStr, _ := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.JsonRespond(w, false, http.StatusBadRequest, err, nil)
		return
	}

	updatedDocument.ID = id
	document, err := service.DocumentService.Update(updatedDocument)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.JsonRespond(w, false, http.StatusNotFound, err, nil)
			return
		} else if errors.Is(err, model.TitleInvalidValue) || errors.Is(err, model.SigneeInvalidValue) {
			utils.JsonRespond(w, false, http.StatusUnprocessableEntity, err, nil)
			return
		}
		utils.JsonRespond(w, false, http.StatusInternalServerError, err, nil)
		return
	}

	utils.JsonRespond(w, true, http.StatusOK, err, document)
	return
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	idStr, _ := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.JsonRespond(w, false, http.StatusBadRequest, err, nil)
		return
	}

	err = service.DocumentService.Delete(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.JsonRespond(w, false, http.StatusNotFound, err, nil)
			return
		}
		utils.JsonRespond(w, false, http.StatusInternalServerError, err, nil)
		return
	}

	utils.JsonRespond(w, true, http.StatusOK, err, nil)
	return
}

func GetByIdHandler(w http.ResponseWriter, r *http.Request) {
	idStr, _ := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.JsonRespond(w, false, http.StatusBadRequest, err, nil)
		return
	}

	document, err := service.DocumentService.Get(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.JsonRespond(w, false, http.StatusNotFound, err, nil)
			return
		}
		utils.JsonRespond(w, false, http.StatusInternalServerError, err, nil)
		return
	}

	utils.JsonRespond(w, true, http.StatusOK, err, document)
	return
}

func GetAllHandler(w http.ResponseWriter, r *http.Request) {
	documents, err := service.DocumentService.GetAll()
	if err != nil {
		utils.JsonRespond(w, false, http.StatusInternalServerError, err, nil)
		return
	}
	utils.JsonRespond(w, true, http.StatusOK, err, documents)
	return
}
