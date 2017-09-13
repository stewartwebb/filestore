package models

import "time"

// ErrorField is for request errors
type ErrorField struct {
	Field   string `json:"field"`
	Message string `json:"message:"`
}

// Error is the core error object returned in all events.
type Error struct {
	StatusCode    int          `codec:"status_code"`
	Message       string       `codec:"message"`
	Documentation string       `codec:"documentation"`
	Fields        []ErrorField `codec:"fields,omitempty"`
}

type mimeType string

var TypePDF mimeType = "application/pdf"
var TypeJPEG mimeType = "image/jpeg"
var TypePNG mimeType = "image/png"
var TypeGIF mimeType = "image/gif"

//MimeGroup is an array of mime types
type mimeGroup []mimeType

// GroupDocument is a collection of mime types considered to be documents
var GroupDocument = mimeGroup{
	TypePDF,
}

// GroupImage a collection of mime types considered to be images
var GroupImage = mimeGroup{
	TypeJPEG,
	TypePNG,
	TypeGIF,
}

// File was the type
type File struct {
	ID            int64             `codec:"id"`
	ApplicationID int64             `codec:"-"`
	Created       *time.Time        `codec:"created,omitempty"`
	Uploaded      *time.Time        `codec:"uploaded,omitempty"`
	Title         string            `codec:"title"`
	Tags          map[string]string `codec:"meta"`
	URL           string            `codec:"url,omitempty"`
	MimeType      string            `codec:"mime_type,omitempty"`
	Size          int               `codec:"size,omitempty"`
	Validate      struct {
		MimeType []string `codec:"mime_type,omitempty"`
	} `codec:"validate,omitempty"`
}

// ValidateInput will return aq list of error fields if the thing is broke.
func (f *File) ValidateInput() (bool, []ErrorField) {
	return true, nil
}
