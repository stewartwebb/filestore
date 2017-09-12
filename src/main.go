package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func Home(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "Hello World!") // send data to client side
}

func GetFiles(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "Get some files")
}

func CreateFile(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "Get some files")
}

func GetFile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "hello, %s", ps.ByName("file"))
}

func UploadFile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "hello, go get amazon s3 file, %s", ps.ByName("file"))
}

func main() {
	router := httprouter.New()
	router.GET("/", Home)
	router.GET("/files", GetFiles)
	router.POST("/files", CreateFile)
	router.GET("/files/:file", GetFile)
	router.PUT("/files/:file", UploadFile)
	err := http.ListenAndServe(":9090", router) // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
