package data

import (
	"strings"

	"github.com/stewartwebb/filestore/src/models"
)

// SaveFile attempts to save files in the database through INSERT or UPDATE
func (db *DB) SaveFile(f models.File) (models.File, error) {
	var query string

	tx, err := db.Begin()
	if err != nil {
		return f, err
	}

	var queryParams = `rApplicationID = ?, dCreated = NOW(), sTitle = ?, sURL = ?, sMimeType = ?, sSize = ?, sValidateMimeTypes = ?`
	var params = []interface{}{
		f.ApplicationID,
		f.Title,
		f.URL,
		f.MimeType,
		f.Size,
		strings.Join(f.Validate.MimeType, ","),
	}
	if f.ID == 0 {
		query = "INSERT INTO tFiles SET " + queryParams
	} else {
		query = "UPDATE tFiles SET " + queryParams + " WHERE pFileID = ?"
		params = append(params, f.ID)
	}
	result, err := tx.Exec(query, params...)
	if err != nil {
		tx.Rollback()
		return f, err
	}
	f.ID, err = result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return f, err
	}

	query = "DELETE FROM tFileTags WHERE rFileID = ?"
	_, err = tx.Exec(query, f.ID)
	if err != nil {
		tx.Rollback()
		return f, err
	}

	if len(f.Tags) > 0 {
		params = []interface{}{}

		query = "INSERT INTO tFileTags (rFileID, sKey, sValue) VALUES "
		for index, element := range f.Tags {
			query += "(?, ?, ?), "
			params = append(params, f.ID)
			params = append(params, index)
			params = append(params, element)
		}
		query = query[:len(query)-2]

		result, err = tx.Exec(query, params...)
		if err != nil {
			tx.Rollback()
			return f, err
		}
	}
	tx.Commit()
	return f, nil
}
