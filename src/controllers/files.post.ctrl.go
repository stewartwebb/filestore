package controllers

import (
	"log"
	"net/http"

	"github.com/stewartwebb/filestore/src/common"
	"github.com/stewartwebb/filestore/src/data"
	"github.com/stewartwebb/filestore/src/models"
)

// CreateFile allows the user to create the meta data for a file
func CreateFile(w http.ResponseWriter, r *http.Request) {
	db := r.Context().Value("conn").(*data.DB)

	var file models.File
	input, err := common.ParseContent(w, r, file)
	if err != nil {
		common.RespondError(w, r, http.StatusBadRequest, "Invalid file meta")
		return
	}
	file = input.(models.File)

	ok, errors := file.ValidateInput()
	if !ok {
		common.RespondError(w, r, http.StatusBadRequest, errors)
		return
	}

	file.ID = 0

	file, err = db.SaveFile(file, false)
	if err != nil {
		log.Printf("[CreateFile] SQL Error: %v", err)
		common.RespondError(w, r, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Location", "*")

	common.RespondOk(w, r, file)
}

// UpdateFile returns a single file meta data
func UpdateFile(w http.ResponseWriter, r *http.Request) {
	common.RespondOk(w, r, "Hello World")
}
