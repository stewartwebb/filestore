package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/stewartwebb/filestore/src/common"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/stewartwebb/filestore/src/data"

	"github.com/stewartwebb/filestore/src/controllers"
)

func Home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!") // send data to client side
}

type middleware func(http.HandlerFunc) http.HandlerFunc

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	db, err := data.NewDB()
	if err != nil {
		log.Fatalf("Database is bad: %v", err)
	}
	a := alice.New(
		common.LoggingHandler,
		common.RecoveryHandler,
		data.AttachDatabaseOld(db),
	)

	r := mux.NewRouter()

	r.HandleFunc("/", Home).Methods("GET")
	r.HandleFunc("/files", controllers.GetFiles).Methods("GET")
	r.HandleFunc("/files", controllers.CreateFile).Methods("POST")
	r.HandleFunc("/files/{id:[0-9]+}", controllers.GetFile).Methods("GET")
	r.HandleFunc("/files/{id:[0-9]+}", controllers.UpdateFile).Methods("PUT")
	r.HandleFunc("/files/{id:[0-9]+}/upload", controllers.UploadFile).Methods("PUT")

	err = http.ListenAndServe(":9090", a.Then(r)) // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
