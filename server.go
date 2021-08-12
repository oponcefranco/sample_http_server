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

func patchUpdateCsvFile(w http.ResponseWriter, r *http.Request) {
	// constructor.io accepts the following:
	// r.FormFile("items")  or r.FormFile("variations")
	file, handler, err := r.FormFile("variations")
	fileName := r.FormValue("file_name")

	if err != nil {
		panic(err)
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
		}
	}(file)
	//CSV file permission: read & write (0666)
	f, err := os.OpenFile(handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
		}
	}(f)
	message := "PATCH request w/ file "+fileName+".csv was sent successfully"
	_, _ = io.WriteString(w, message)
	_, _ = io.Copy(f, file)
	logger := log.New(os.Stdout, "API Request: ", log.LstdFlags)
	logger.Println(message)
}

func getIndex(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, "Welcome Home!")
}

func main() {
	flag.StringVar(&listenAddr, "listen-addr", ":8080", "server listen address")
	flag.Parse()
	logger := log.New(os.Stdout, "http: ", log.LstdFlags)
	logger.Println("\nServer is starting on port"+listenAddr+"...")

	router := mux.NewRouter().StrictSlash(true)
	catalog := router.PathPrefix("/v1/catalog").Subrouter()

	catalog.Path("").Methods(http.MethodGet).HandlerFunc(getIndex)
	catalog.Path("").Methods(http.MethodPatch).HandlerFunc(patchUpdateCsvFile)

	router.HandleFunc("/", getIndex).Methods("GET")
	router.HandleFunc("/file", patchUpdateCsvFile).Methods("PATCH")

	log.Fatal(http.ListenAndServe(":8080", router))
}
