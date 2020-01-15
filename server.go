package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dumacp/updatefs/datastore"
	"github.com/dumacp/updatefs/updatedata"
	"github.com/gorilla/mux"
)

var (
	dir         string
	pathcss     string
	pathdb      string
	pathfilesdb string
	files       datastore.FileStore
	updates     updatedata.UpdateStore
	socket      string
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
	flag.StringVar(&dir, "dir", "/data/all/files", "the directory to serve files from. Defaults to the current dir")
	flag.StringVar(&pathdb, "pathupdatesdb", "/data/all/updates.db", "path to updates database")
	flag.StringVar(&pathfilesdb, "pathfilesdb", "/data/all/files.db", "path to files database")
	flag.StringVar(&pathcss, "pathcss", "/data/all/css", "path to css")
	flag.StringVar(&socket, "listensocket", listensocket, "socket to listen")
	files = &datastore.Files{}
	updates = &updatedata.UpdateData{}

}

func main() {
	flag.Parse()
	files.Initialize(pathfilesdb, dir)
	updates.Initialize(pathdb)

	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Fatalf("error: cannot create dir %q", dir)
	}
	if err := os.MkdirAll(pathcss, 0755); err != nil {
		log.Fatalf("error: cannot create dir %q", dir)
	}

	r := mux.NewRouter()
	log.Println("filedata api")
	fileserver1 := r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(dir))))
	if methods, err := fileserver1.GetMethods(); err != nil {
		for i, v := range methods {
			log.Printf("Method %d: %s", i, v)
		}
	}
	fileserver2 := r.PathPrefix("/updatevoc/static/").Handler(http.StripPrefix("/updatevoc/static/", http.FileServer(http.Dir(dir))))
	if methods, err := fileserver2.GetMethods(); err != nil {
		for i, v := range methods {
			log.Printf("Method %d: %s", i, v)
		}
	}
	filecss := r.PathPrefix("/updatevoc/css/").Handler(http.StripPrefix("/updatevoc/css/", http.FileServer(http.Dir(pathcss))))
	if methods, err := filecss.GetMethods(); err != nil {
		for i, v := range methods {
			log.Printf("Method %d: %s", i, v)
		}
	}
	apiv1 := r.PathPrefix("/updatevoc/api/v1").Subrouter()
	apiv1.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "api v1")
	})

	apiv1.HandleFunc("/device/{devicename}", searchByDeviceName).Methods(http.MethodGet)
	apiv1.HandleFunc("/device", allDevices).Methods(http.MethodGet)

	apiv2 := r.PathPrefix("/updatevoc/api/v2").Subrouter()
	apiv2.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "api v2")
	})
	apiv2.HandleFunc("/files/device/{devicename}", searchByDeviceName).Methods(http.MethodGet)
	apiv2.HandleFunc("/files/delete", deleteFiles).Methods(http.MethodPost)
	apiv2.HandleFunc("/files/md5/{md5}", searchByMD5).Methods(http.MethodGet)
	apiv2.HandleFunc("/files", allDevices).Methods(http.MethodGet)
	apiv2.HandleFunc("/files", createFile).Methods(http.MethodPost)
	apiv2.HandleFunc("/updates/device/{devicename}", searchUpdateByDeviceName).Methods(http.MethodGet)
	apiv2.HandleFunc("/updates/file/{md5}", searchUpdateByFile).Methods(http.MethodGet)
	apiv2.HandleFunc("/updates", createUpdate).Methods(http.MethodPost)

	datasite := r.PathPrefix("/updatevoc/data").Subrouter()
	datasite.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		templateForm, _ := template.New("uploadfile").Parse(formCreateFile)
		if err := templateForm.Execute(w, nil); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
	datasite.HandleFunc("/delete", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		templateForm, _ := template.New("deletefiles").Parse(formDeleteFile)
		store := files.AllData()
		log.Printf("%v", *store)
		if err := templateForm.Execute(w, *store); err != nil {
			log.Printf("error: tmeplate delete, %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})

	datasite.HandleFunc("/updates/{devicename}", func(w http.ResponseWriter, r *http.Request) {
		pathParams := mux.Vars(r)
		devicename := pathParams["devicename"]
		templateForm, _ := template.New("lastUpdates").Parse(viewDeviceUpdate)
		store, err := updates.SearchUpdateDataDevice([]byte(devicename), 0, 100, 0)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("error: %s", err)))
			return
		}
		data := struct {
			Name    string
			Updates []*updatedata.Updatedatadevice
		}{
			devicename,
			*store,
		}
		if err := templateForm.Execute(w, data); err != nil {
			log.Printf("error: template lastUpdates, %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}).Methods(http.MethodGet)

	datasite.HandleFunc("/updates/{devicename}", func(w http.ResponseWriter, r *http.Request) {
		pathParams := mux.Vars(r)
		devicename := pathParams["devicename"]
		templateForm, _ := template.New("lastUpdates").Parse(viewDeviceUpdate)
		store, err := updates.SearchUpdateDataDevice([]byte(devicename), 0, 100, 0)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("error: %s", err)))
			return
		}
		data := struct {
			Name   string
			Update []*updatedata.Updatedatadevice
		}{
			devicename,
			*store,
		}
		if err := templateForm.Execute(w, data); err != nil {
			log.Printf("error: template lastUpdates, %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}).Methods(http.MethodGet)

	srv := &http.Server{
		Handler: r,
		Addr:    socket,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatalln(srv.ListenAndServe())
}
