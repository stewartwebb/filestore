package controllers

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"syscall"

	"github.com/stewartwebb/filestore/src/data"

	"github.com/gorilla/mux"
	"github.com/stewartwebb/filestore/src/common"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3crypto"
)

// UploadFile returns a single file meta data
func UploadFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		common.RespondError(w, r, http.StatusBadRequest)
		return
	}

	db := r.Context().Value("conn").(*data.DB)
	f, err := db.GetFile(id)
	if err != nil {
		if err == sql.ErrNoRows {
			common.RespondError(w, r, http.StatusNotFound)
			return
		}
		common.RespondError(w, r, http.StatusInternalServerError)
		return
	}

	if f.Size > 0 {
		common.RespondError(w, r, http.StatusConflict, "Cannot upload file again")
		return
	}

	defer r.Body.Close()
	file, err := os.Create("/tmp/" + strconv.FormatInt(id, 10))
	if err != nil {
		log.Printf("[UploadFile] Error Opening File: %v", err)
		common.RespondError(w, r, http.StatusInternalServerError)
		return
	}
	defer file.Close()

	if _, err = io.Copy(file, r.Body); err != nil {
		log.Printf("[UploadFile] Error saving file: %v", err)
		common.RespondError(w, r, http.StatusInternalServerError)
		return
	}

	ok, err := scanFile(id)
	if err != nil {
		log.Printf("[UploadFile] Error scanning file: %v", err)
		common.RespondError(w, r, http.StatusInternalServerError)
		return
	}

	if !ok {
		common.RespondError(w, r, http.StatusBadRequest, "Failed virus scan")
		return
	}
	size, fileType, err := uploadFile(id)
	if err != nil {
		log.Printf("[UploadFile] Upload Error: %v", err)
		common.RespondError(w, r, http.StatusInternalServerError)
		return
	}

	f.Size = size
	f.MimeType = fileType

	f, err = db.SaveFile(f, true)
	if err != nil {
		log.Printf("[UploadFile] SQL Error: %v", err)
		common.RespondError(w, r, http.StatusInternalServerError)
		return
	}

	common.RespondOk(w, r, f)

}

// scanFile will virus scan a file to check it is legit.
func scanFile(id int64) (bool, error) {
	cmd := exec.Command("clamdscan", "/tmp/"+strconv.FormatInt(id, 10))

	var waitStatus syscall.WaitStatus
	if err := cmd.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			waitStatus = exitError.Sys().(syscall.WaitStatus)
		}
	} else {
		// Success
		waitStatus = cmd.ProcessState.Sys().(syscall.WaitStatus)
	}

	if waitStatus.ExitStatus() == 0 {
		return true, nil
	} else if waitStatus.ExitStatus() == 1 {
		return false, nil
	} else {
		return false, errors.New("Clamdscan reported an error")
	}
}

func uploadFile(id int64) (int64, string, error) {
	file, err := os.Open("/tmp/" + strconv.FormatInt(id, 10))
	if err != nil {
		fmt.Printf("err opening file: %s", err)
	}
	defer file.Close()

	fileInfo, _ := file.Stat()
	var size int64 = fileInfo.Size()

	buffer := make([]byte, size)
	file.Read(buffer)
	fileBytes := bytes.NewReader(buffer)
	fileType := http.DetectContentType(buffer)

	sess := session.Must(session.NewSession())

	keywrap := s3crypto.NewKMSKeyGenerator(kms.New(sess), common.AppConfig.KmsARN)
	builder := s3crypto.AESGCMContentCipherBuilder(keywrap)
	client := s3crypto.NewEncryptionClient(sess, builder)

	params := &s3.PutObjectInput{
		Bucket:      aws.String(common.AppConfig.AwsBucket),
		Key:         aws.String(strconv.FormatInt(id, 10)),
		Body:        fileBytes,
		ContentType: aws.String(fileType),
	}
	_, err = client.PutObject(params)
	if err != nil {
		return 0, "", err
	}
	return size, fileType, nil
}
