package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dumacp/updatefs/datastore"
	"github.com/gorilla/mux"
)

var (
	dir   string
	files datastore.FileStore
)

func init() {
	flag.StringVar(&dir, "dir", ".", "the directory to serve files from. Defaults to the current dir")
}

func main() {
	flag.Parse()

	r := mux.NewRouter()
	log.Println("bookdata api")
	fileserver := r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(dir))))
	if methods, err := fileserver.GetMethods(); err != nil {
		for i, v := range methods {
			log.Printf("Method %d: %s", i, v)
		}
	}
	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "api v1")
	})
	api.HandleFunc("/filedata/{deviceName}", searchByDeviceName).Methods(http.MethodGet)
	api.HandleFunc("/file", createFile).Methods(http.MethodPost)
	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatalln(srv.ListenAndServe())
}
