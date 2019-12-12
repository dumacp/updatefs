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

const (
	listensocket = "0.0.0.0:8000"
)

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

func init() {
	defer timeTrack(time.Now(), "file load")
	flag.StringVar(&dir, "dir", "/data/all", "the directory to serve files from. Defaults to the current dir")
	files = &datastore.Files{}

}

func main() {
	flag.Parse()
	files.Initialize(dir)

	r := mux.NewRouter()
	log.Println("filedata api")
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
	api.HandleFunc("/filedata/device/{devicename}", searchByDeviceName).Methods(http.MethodGet)
	api.HandleFunc("/filedata/device", allDevices).Methods(http.MethodGet)
	api.HandleFunc("/file", createFile).Methods(http.MethodPost)
	srv := &http.Server{
		Handler: r,
		Addr:    listensocket,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatalln(srv.ListenAndServe())
}
