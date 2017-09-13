package controllers

import (
	"net/http"

	"github.com/stewartwebb/filestore/src/common"
)

// GetFiles returns a list of files
func GetFiles(w http.ResponseWriter, r *http.Request) {
	common.RespondOk(w, r, "Hello World")
}
