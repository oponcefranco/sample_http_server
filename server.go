package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

func UploadFile(w http.ResponseWriter, r *http.Request) {
	file, handler, err := r.FormFile("file")
	fileName := r.FormValue("file_name")
	if err != nil {
		panic(err)
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
		}
	}(file)

	f, err := os.OpenFile(handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {

		}
	}(f)
	_, _ = io.WriteString(w, "File "+fileName+" Uploaded successfully")
	_, _ = io.Copy(f, file)
}

func homeLink(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, "Welcome home!")
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	catalog := router.PathPrefix("/v1/catalog").Subrouter()
	//usersR.Path("").Methods(http.MethodGet).HandlerFunc(getHome)
	catalog.Path("").Methods(http.MethodGet).HandlerFunc(homeLink)
	catalog.Path("").Methods(http.MethodPost).HandlerFunc(UploadFile)
	catalog.Path("").Methods(http.MethodPatch).HandlerFunc(UploadFile)

	router.HandleFunc("/", homeLink)
	router.HandleFunc("/file", UploadFile).Methods("POST")
	router.HandleFunc("/file", UploadFile).Methods("PATCH")
	log.Fatal(http.ListenAndServe(":8081", router))
}
