package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

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

func createFile(w http.ResponseWriter, r *http.Request) {}

func createUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "body not parsed"}`))
		return
	}

	date, _ := strconv.Atoi(r.FormValue("date"))
	filemd5 := r.FormValue("date")
	devicename := r.FormValue("devicename")
	ipclient := r.RemoteAddr
	filedata := r.FormValue("filedata")

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
