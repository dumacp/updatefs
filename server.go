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
	pathdb      string
	pathfilesdb string
	files       datastore.FileStore
	updates     updatedata.UpdateStore
	socket      string
)

const (
	listensocket = "0.0.0.0:8000"
	formFiledata = `<html>
    <head>
    <title></title>
    </head>
    <body>
		<form action="/updatevoc/api/v2/files" enctype="multipart/form-data" method="post">
			<div>
				<label>description:</label>
				<input type="text" name="description">
			</div>	
			<div>
				<label>version:</label>
				<input type="text" name="version">
			</div>
			<div>
				<label>reference:</label>
				<input type="number" name="reference">
			</div>
			<div>
				<label>path:</label>
				<input type="text" name="path">
			</div>
			<div>
				<input type="file" name="fileToUpload" id="fileToUpload">
			</div>
			<div>
				<input type="submit" value="Upload">
			</div>
        </form>
    </body>
</html>`
	formDeleteFile = `<!DOCTYPE html>
	<html>
	<body>
	
	<h1>Show currents files:</h1>
	
	<form action="/updatevoc/api/v2/files/delete" enctype="multipart/form-data" method="post">
	{{range .Store}}

		<input type="checkbox" name="role[]" value="{{.ID}}">{{.Name}}<br>	  
	{{end}}
		<input type="submit" value="Submit">
	</form>
	
	</body>
	</html>`
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
	flag.StringVar(&socket, "listensocket", listensocket, "socket to listen")
	files = &datastore.Files{}
	updates = &updatedata.UpdateData{}

}

func main() {
	flag.Parse()
	files.Initialize(pathfilesdb)
	updates.Initialize(pathdb)

	if err := os.MkdirAll(dir, 0755); err != nil {
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

	uploadfile := r.PathPrefix("/updatevoc/upload").Subrouter()
	uploadfile.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		templateForm, _ := template.New("uploadfile").Parse(formFiledata)
		if err := templateForm.Execute(w, nil); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
	deletefiles := r.PathPrefix("/updatevoc/delete").Subrouter()
	deletefiles.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		templateForm, _ := template.New("deletefiles").Parse(formDeleteFile)
		if err := templateForm.Execute(w, *files.AllData()); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})

	srv := &http.Server{
		Handler: r,
		Addr:    socket,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatalln(srv.ListenAndServe())
}
