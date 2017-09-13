package common

import (
	"fmt"
	"net/http"

	"github.com/ugorji/go/codec"
)

// ParseContent looks for content-type header and does cool stuff
func ParseContent(w http.ResponseWriter, r *http.Request, obj interface{}) (interface{}, error) {
	content := r.Header.Get("Content-Type")

	fmt.Println("Content-Type: " + content)

	var err error
	json := &codec.JsonHandle{}
	cbor := &codec.CborHandle{}

	switch content {
	case "application/cbor":
		dec := codec.NewDecoder(r.Body, cbor)
		err = dec.Decode(&obj)
	case "application/json", "application/json;charset=UTF-8":
		dec := codec.NewDecoder(r.Body, json)
		err = dec.Decode(&obj)
	}
	if err != nil {
		return nil, err
	}
	return obj, nil
}
