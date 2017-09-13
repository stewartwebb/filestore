package controllers

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"syscall"

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

	defer r.Body.Close()
	f, err := os.Create("/tmp/" + strconv.FormatInt(id, 10))
	if err != nil {
		log.Printf("[UploadFile] Error Opening File: %v", err)
		common.RespondError(w, r, http.StatusInternalServerError)
		return
	}
	defer f.Close()

	if _, err = io.Copy(f, r.Body); err != nil {
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
	ok, err = uploadFile(id)
	if err != nil {
		log.Printf("[UploadFile] Upload Error: %v", err)
		common.RespondError(w, r, http.StatusInternalServerError)
		return
	}
	common.RespondOk(w, r, ok)

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

func uploadFile(id int64) (bool, error) {
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

	arn := "arn:aws:s3:::pinnacleweb/" + strconv.FormatInt(id, 10)

	keywrap := s3crypto.NewKMSKeyGenerator(kms.New(sess), arn)
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
		return false, err
	}
	return true, nil
}
