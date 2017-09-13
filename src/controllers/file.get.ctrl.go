package controllers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/stewartwebb/filestore/src/common"
	"github.com/stewartwebb/filestore/src/data"
)

// GetFile returns a single file meta data
func GetFile(w http.ResponseWriter, r *http.Request) {
	db := r.Context().Value("conn").(*data.DB)
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		common.RespondError(w, r, http.StatusBadRequest)
		return
	}

	file, err := db.GetFile(id)
	fmt.Println(err)
	if err != nil {
		if err == sql.ErrNoRows {
			common.RespondError(w, r, http.StatusNotFound)
			return
		}
		log.Printf("[GetFile] SQL Error: %v", err)
		common.RespondError(w, r, http.StatusInternalServerError)
		return
	}

	common.RespondOk(w, r, file)
}
