package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

var listenAddr string

func updateFile(w http.ResponseWriter, r *http.Request) {
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

func home(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, "Welcome home!")
}

func main() {
	flag.StringVar(&listenAddr, "listen-addr", ":8080", "server listen address")
	flag.Parse()
	logger := log.New(os.Stdout, "http: ", log.LstdFlags)
	logger.Println("Server is starting...")

	router := mux.NewRouter().StrictSlash(true)
	catalog := router.PathPrefix("/v1/catalog").Subrouter()

	catalog.Path("").Methods(http.MethodGet).HandlerFunc(home)
	catalog.Path("").Methods(http.MethodPost).HandlerFunc(updateFile)
	catalog.Path("").Methods(http.MethodPatch).HandlerFunc(updateFile)

	router.HandleFunc("/", home).Methods("GET")
	router.HandleFunc("/file", updateFile).Methods("POST")
	router.HandleFunc("/file", updateFile).Methods("PATCH")

	log.Fatal(http.ListenAndServe(":8080", router))
}
