package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

var listenAddr string

type LogResponseWriter struct {
	http.ResponseWriter
	statusCode int
	buf        bytes.Buffer
}

func NewLogResponseWriter(w http.ResponseWriter) *LogResponseWriter {
	return &LogResponseWriter{ResponseWriter: w}
}

func (w *LogResponseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *LogResponseWriter) Write(body []byte) (int, error) {
	w.buf.Write(body)
	return w.ResponseWriter.Write(body)
}

type LogMiddleware struct {
	logger *log.Logger
}

func NewLogMiddleware(logger *log.Logger) *LogMiddleware {
	return &LogMiddleware{logger: logger}
}

func (m *LogMiddleware) Func() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			startTime := time.Now()

			logRespWriter := NewLogResponseWriter(w)
			next.ServeHTTP(logRespWriter, r)

			m.logger.Printf(
				"duration=%s status=%d body=%s",
				time.Since(startTime).String(),
				logRespWriter.statusCode,
				logRespWriter.buf.String())
		})
	}
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprintf(w, "Welcome Home!")
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprintf(w, "Healthy!")
}

func CsvFileHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	parseCsvFile(w, r)
}

func parseCsvFile(w http.ResponseWriter, r *http.Request) {
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
	message := "PATCH request was sent successfully: " + fileName + ".csv"
	_, _ = io.WriteString(w, message)
	_, _ = io.Copy(f, file)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}

func main() {
	logger := log.New(os.Stdout, "", log.Lmicroseconds)

	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration the server gracefully waits for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()
	flag.StringVar(&listenAddr, "listen-addr", ":8080", "server listen address")
	flag.Parse()

	logger.Println("\nServer is starting on port" + listenAddr + "...")

	router := mux.NewRouter().StrictSlash(true)
	router.Use(loggingMiddleware)

	srv := &http.Server{
		Addr:         "127.0.0.1:8080",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	catalog := router.PathPrefix("/v1/catalog").Subrouter()

	catalog.Path("/healthcheck").Methods(http.MethodGet).HandlerFunc(HealthCheckHandler)
	catalog.Path("").Methods(http.MethodGet).HandlerFunc(IndexHandler)
	catalog.Path("").Methods(http.MethodPatch).HandlerFunc(CsvFileHandler)

	router.HandleFunc("/", IndexHandler).Methods("GET")
	router.HandleFunc("/healthcheck", HealthCheckHandler).Methods("GET")
	router.HandleFunc("/", CsvFileHandler).Methods("PATCH")

	logMiddleware := NewLogMiddleware(logger)
	router.Use(logMiddleware.Func())

	logger.Fatalln(http.ListenAndServe(":8080", router))
}
