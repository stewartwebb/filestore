package common

import (
	"log"
	"net/http"

	"github.com/stewartwebb/filestore/src/models"
	"github.com/ugorji/go/codec"
)

// RespondOk is used to return data to the client with a 200
func RespondOk(w http.ResponseWriter, r *http.Request, response interface{}) {
	accept := r.Header.Get("Accept")

	var err error
	json := &codec.JsonHandle{}
	json.Canonical = true
	cbor := &codec.CborHandle{}

	switch accept {
	case "application/cbor":
		w.Header().Set("Content-Type", "application/cbor")
		enc := codec.NewEncoder(w, cbor)
		err = enc.Encode(response)
	default:
		w.Header().Set("Content-Type", "application/json")
		enc := codec.NewEncoder(w, json)
		err = enc.Encode(response)
	}
	if err != nil {
		log.Printf("[RespondOk] Encode Error: %v, ", err)
		RespondError(w, r, http.StatusInternalServerError)
		return
	}
}

// RespondError is used to return errors to the client in default format
func RespondError(w http.ResponseWriter, r *http.Request, status int, detail ...interface{}) {
	accept := r.Header.Get("Accept")

	msg := http.StatusText(status)
	doc := "TBA"
	var fields []models.ErrorField

	/*
	  var actualError error = nil
	  var request *http.Request
	*/

	for _, v := range detail {
		if str, ok := v.(string); ok {
			if str[0] == '#' {
				doc = str
			} else {
				msg = str
			}
		} else if f, ok := v.([]models.ErrorField); ok {
			fields = f
		} /*else if e, ok := v.(error); ok {
		    actualError = e
		  } else if r, ok := v.(*http.Request); ok {
		    request = r
		  }*/
	}

	response := models.Error{
		StatusCode:    status,
		Message:       msg,
		Documentation: AppConfig.WebAddress + doc,
		Fields:        fields,
	}

	var err error
	json := &codec.JsonHandle{}
	json.Canonical = true
	cbor := &codec.CborHandle{}

	switch accept {
	case "application/cbor":
		w.Header().Set("Content-Type", "application/cbor")
		w.WriteHeader(status)
		enc := codec.NewEncoder(w, cbor)
		err = enc.Encode(response)
	default:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		enc := codec.NewEncoder(w, json)
		err = enc.Encode(response)
	}
	if err != nil {
		log.Printf("[RespondOk] Encode Error: %v, ", err)
		RespondError(w, r, http.StatusInternalServerError)
		return
	}
}
