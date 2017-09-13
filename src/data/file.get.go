package data

import (
	"database/sql"
	"strings"

	"github.com/stewartwebb/filestore/src/models"
)

// GetFile reuturns a single file to the user.
func (db *DB) GetFile(id int64) (models.File, error) {
	var f = models.File{Tags: make(map[string]string, 0)}
	query := `SELECT pFileID, rApplicationID, dCreated, dUploaded, sTitle, sURL, sMimeType, sSize, sValidateMimeTypes, sKey, sValue
	FROM tFiles JOIN tFileTags ON rFileID = pFileID
	WHERE pFileID = ?`

	rows, err := db.Query(query, id)
	if err != nil {
		return f, err
	}
	defer rows.Close()
	for rows.Next() {
		var mimeType string
		var key string
		var value string
		err = rows.Scan(&f.ID, &f.ApplicationID, &f.Created, &f.Uploaded, &f.Title, &f.URL, &f.MimeType, &f.Size,
			&mimeType, &key, &value)
		if err != nil {
			return f, err
		}
		f.Validate.MimeType = strings.Split(mimeType, ",")
		f.Tags[key] = value
	}
	err = rows.Err()
	if err != nil {
		return f, err
	}
	if f.ID == 0 {
		return f, sql.ErrNoRows
	}
	return f, nil
}
