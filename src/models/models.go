package models

// File was the type
type File struct {
	ID    int64             `codec:"id"`
	Title string            `codec:"title"`
	Meta  map[string]string `codec:"meta"`
	URL   string            `codec:"url"`
}
