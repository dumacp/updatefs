package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/dumacp/updatefs/loader"
	"github.com/dumacp/updatefs/updatedata"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func searchByMD5(w http.ResponseWriter, r *http.Request) {
	queries := mux.Vars(r)

	w.Header().Set("Content-Type", "application/json")
	if val, ok := queries["md5"]; ok {
		data := files.SearchMD5(val)
		if data != nil {
			b, err := json.Marshal(data)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"error": "error marshalling data"}`))
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write(b)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{"error": "not found"}`))
}

func searchByID(w http.ResponseWriter, r *http.Request) {
	queries := mux.Vars(r)

	w.Header().Set("Content-Type", "application/json")
	if val, ok := queries["id"]; ok {
		data := files.SearchID(val)
		if data != nil {
			b, err := json.Marshal(data)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"error": "error marshalling data"}`))
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write(b)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{"error": "not found"}`))
}

func searchByDeviceName(w http.ResponseWriter, r *http.Request) {
	pathParams := mux.Vars(r)
	//log.Printf("%+v", pathParams)
	w.Header().Set("Content-Type", "application/json")
	limit, err := getLimitParam(r)
	skip, err := getSkipParam(r)
	date, err := getDateParam(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "invalid datatype for parameter"}`))
		return
	}
	if val, ok := pathParams["devicename"]; ok {
		//og.Printf("%+v", val)
		data := *files.SearchDeviceName(val, date, limit, skip)
		b, err := json.Marshal(data)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "error marshalling data"}`))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(b)
		return
	}
	w.WriteHeader(http.StatusNotFound)
}

func allDevices(w http.ResponseWriter, r *http.Request) {
	mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")
	data := *files.AllData()
	b, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "error marshalling data"}`))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func searchUpdateByDeviceName(w http.ResponseWriter, r *http.Request) {
	pathParams := mux.Vars(r)
	//log.Printf("%+v", pathParams)
	w.Header().Set("Content-Type", "application/json")
	limit, err := getLimitParam(r)
	skip, err := getSkipParam(r)
	date, err := getDateParam(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "invalid datatype for parameter"}`))
		return
	}
	if val, ok := pathParams["devicename"]; ok {
		//og.Printf("%+v", val)
		data, err := updates.SearchUpdateDataDevice([]byte(val), date, limit, skip)
		if err != nil {
			log.Printf("error SearchUpdateDataDevice: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "error collenting data"}`))
			return
		}
		b, err := json.Marshal(data)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "error marshalling data"}`))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(b)
		return
	}
	w.WriteHeader(http.StatusNotFound)
}

func searchUpdateByFile(w http.ResponseWriter, r *http.Request) {
	pathParams := mux.Vars(r)
	log.Printf("%+v", pathParams)
	w.Header().Set("Content-Type", "application/json")
	limit, err := getLimitParam(r)
	skip, err := getSkipParam(r)
	date, err := getDateParam(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "invalid datatype for parameter"}`))
		return
	}
	if val, ok := pathParams["md5"]; ok {
		log.Printf("%+v", val)
		data, err := updates.SearchUpdateDataFile([]byte(val), date, limit, skip)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "error collenting data"}`))
			return
		}
		b, err := json.Marshal(data)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "error marshalling data"}`))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(b)
		return
	}
	w.WriteHeader(http.StatusNotFound)
}

func allUpdateDevices(w http.ResponseWriter, r *http.Request) {
	mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")
	data := *files.AllData()
	b, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "error marshalling data"}`))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func createFile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "body not parsed"}`))
		return
	}

	fileupload, _, err := r.FormFile("fileToUpload")
	if err != nil {
		log.Printf("error fileupload: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "upload file data"}`))
		return
	}
	defer fileupload.Close()
	desc := r.FormValue("description")
	ref, _ := strconv.Atoi(r.FormValue("reference"))
	version := r.FormValue("version")
	path := r.FormValue("path")

	filePath := filepath.Clean(fmt.Sprintf("%s/migracion_%s.zip", dir, version))
	if _, err := os.Stat(filePath); err == os.ErrNotExist {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(`"error": "file upload already exist, %q"}`, filePath)))
		return
	}
	data, err := ioutil.ReadAll(fileupload)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "read file upload"}`))
		return
	}
	if err := os.MkdirAll(filepath.Base(filePath), 0755); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "write file data"}`))
		return
	}
	if err := ioutil.WriteFile(filePath, data, 0644); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "write file data"}`))
		return
	}
	filed := new(loader.FileData)
	if len(path) <= 0 {
		filed.DeviceName = append(filed.DeviceName, "all")
	} else {
		pathtrim := strings.TrimSpace(path)
		paths := strings.Split(pathtrim, ",")
		for _, v := range paths {
			filed.DeviceName = append(filed.DeviceName, v)
		}
	}

	filed.FilePath, _ = filepath.Rel(dir, filePath)
	filed.Name = filepath.Base(filePath)
	filed.ID = uuid.New().String()
	filed.Date = time.Now().Unix()
	md5sum := md5.Sum(data)
	filed.Md5 = hex.EncodeToString(md5sum[0:])
	filed.Description = desc
	filed.Ref = ref
	filed.Version = version

	filev, err := json.Marshal(filed)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "dont parse file""}`))
		return
	}

	if !files.CreateFile(filed) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "dont persist file""}`))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(filev)

}

func deleteFiles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "body not parsed"}`))
		return
	}

	roleb := r.FormValue("role")
	var rolei interface{}

	if err := json.Unmarshal([]byte(roleb), rolei); err != nil {
		log.Panicln(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "value array files not parsed"}`))
		return
	}

	switch vt := rolei.(type) {
	case []string:
		for _, v := range vt {

		}
	}

	filed := new(loader.FileData)

	if err := json.Unmarshal([]byte(filedata), filed); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "value filedata not parsed"}`))
		return
	}
	upDevice := &updatedata.Updatedatadevice{
		ID:        uuid.New().String(),
		Date:      date,
		Filedata:  filed,
		IPRequest: ipclient,
	}

	if err := updates.NewUpdateDataDevice([]byte(devicename), []byte(filemd5), upDevice); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "NewUpdateDataDevice not created"}`))
		return
	}
	b, err := json.Marshal(upDevice)
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

/**
func createBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := r.ParseForm()
	if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error": "body not parsed"}`))
			return
	}

	avgRating, _ := strconv.ParseFloat(r.FormValue("AverageRating"), 64)
	numPages, _ := strconv.Atoi(r.FormValue("NumPages"))
	ratings, _ := strconv.Atoi(r.FormValue("Ratings"))
	reviews, _ := strconv.Atoi(r.FormValue("Reviews"))

	ok := books.CreateBook(&loader.BookData{
			BookID:        r.FormValue("BookID"),
			Title:         r.FormValue("Title"),
			Authors:       r.FormValue("Authors"),
			AverageRating: avgRating,
			ISBN:          r.FormValue("ISBN"),
			ISBN13:        r.FormValue("ISBN13"),
			LanguageCode:  r.FormValue("LanguageCode"),
			NumPages:      numPages,
			Ratings:       ratings,
			Reviews:       reviews,
	})
	if ok {
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(`{"success": "created"}`))
			return
	}
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{"error": "not created"}`))
}
/**/

func getDateParam(r *http.Request) (int, error) {
	limit := 0
	queryParams := r.URL.Query()
	l := queryParams.Get("date")
	if l != "" {
		val, err := strconv.Atoi(l)
		if err != nil {
			return limit, err
		}
		limit = val
	}
	return limit, nil
}

func getLimitParam(r *http.Request) (int, error) {
	limit := 0
	queryParams := r.URL.Query()
	l := queryParams.Get("limit")
	if l != "" {
		val, err := strconv.Atoi(l)
		if err != nil {
			return limit, err
		}
		limit = val
	}
	return limit, nil
}

func getSkipParam(r *http.Request) (int, error) {
	skip := 0
	queryParams := r.URL.Query()
	l := queryParams.Get("skip")
	if l != "" {
		val, err := strconv.Atoi(l)
		if err != nil {
			return skip, err
		}
		skip = val
	}
	return skip, nil
}
