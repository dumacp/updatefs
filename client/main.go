package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"github.com/dumacp/updatefs/loader"

	"github.com/boltdb/bolt"
)

var urlin string

type devicedata struct {
	Name         string `json:"name"`
	Bussinesunit string `json:"bussinesunit"`
	Route        string `json:"route"`
	Organization string `json:"org"`
	Type         string `json:"type"`
}

const (
	dirDB          = "/tmp/SD/boltdbs"
	nameDB         = "updatefs"
	fileserverdir  = "static"
	pathupdatefile = "/tmp/SD/update/migracion.zip"
)

func init() {
	flag.StringVar(&urlin, "url", "http://127.0.0.1:8000", "url server")
}

func main() {
	flag.Parse()
	if err := os.MkdirAll(dirDB, 0755); err != nil {
		log.Fatalln(err)
	}
	if err := os.MkdirAll(path.Dir(pathupdatefile), 0755); err != nil {
		log.Fatalln(err)
	}

	db, err := bolt.Open(fmt.Sprintf("%s/%s", dirDB, nameDB), 0666, nil)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	//TODO: load device data
	// device := new(devicedata)
	// if db.View(func(tx *bolt.Tx) error {

	// 	bk := tx.Bucket([]byte("devicedata"))
	// 	if bk == nil {
	// 		return nil
	// 	}
	// 	val := bk.Get([]byte("data"))
	// 	if val == nil {
	// 		return nil
	// 	}
	// 	if err := json.Unmarshal(val, device); err != nil {
	// 		return err
	// 	}
	// 	return nil
	// }); err != nil {
	// 	log.Fatalln(err)
	// }

	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalf("Error: there is not hostname! %s", err)
	}

	//TODO: load device data
	// flagupdateDB :=  false
	// if len(hostname) > 0 {
	// 	if !strings.Contains(device.Name, hostname)  {
	// 		device.Name = hostname
	// 		flagupdateDB = true
	// 	}
	// }

	// if flagupdateDB {
	// 	if db.Update(func(tx *bolt.Tx) error {
	// 		return nil

	// 	}); err != nil {
	// 		log.Fatalln(err)
	// 	}
	// }

	filedata := new(loader.FileData)
	if err := db.View(func(tx *bolt.Tx) error {

		bk := tx.Bucket([]byte("updates"))
		if bk == nil {
			return bolt.ErrBucketNotFound
		}
		val := bk.Get([]byte("lastupdate"))
		if val == nil {
			return nil
		}
		if err := json.Unmarshal(val, filedata); err != nil {
			return err
		}
		return nil
	}); err != nil {
		if err == bolt.ErrBucketNotFound {
			if err := db.Update(func(tx *bolt.Tx) error {
				_, err := tx.CreateBucketIfNotExists([]byte("updates"))
				if err != nil {
					return err
				}
				return nil
			}); err != nil {
				log.Fatalf("ERROR create bucket update: %s", err)
			}
		} else {
			log.Fatalf("ERROR import lastupdate data: %s", err)
		}
	}

	tick := time.NewTicker(30 * time.Minute)
	start := time.After(3 * time.Second)

	loopfunc := func() {
		store, err := NewRequestFilesByDevicename(urlin, hostname, int(filedata.Date), 1, 0)
		if err != nil {
			log.Printf("ERROR NewRequestFilesByDevicename: %s", err)
			return
		}
		if store == nil || len(*store) <= 0 {
			store, err = NewRequestFilesByDevicename(urlin, "all", int(filedata.Date), 1, 0)
			if err != nil {
				log.Printf("ERROR NewRequestFilesByDevicename all: %s", err)
				return
			}
			if store == nil || len(*store) <= 0 {
				return
			}
		}

		filedatanow := (*store)[0]
		fmt.Printf("%v, %v\n", filedatanow, filedata)
		if filedatanow.Date > filedata.Date &&
			(len(filedata.Md5) <= 0 ||
				!strings.Contains(filedatanow.Md5, filedata.Md5)) {

			fileurl := fmt.Sprintf("%s/%s/%s", urlin, fileserverdir, filedatanow.FilePath)
			err := DownloadFile(fileurl, pathupdatefile)
			if err != nil {
				log.Printf("ERROR DownloadFile: %s", err)
				return
			}
			log.Printf("UPDATE FILE DOWNLOAD: %+v", filedatanow)
			if err := db.Update(func(tx *bolt.Tx) error {
				bk := tx.Bucket([]byte("updates"))
				if bk == nil {
					return bolt.ErrBucketNotFound
				}
				val, err := json.Marshal(filedatanow)
				if err != nil {
					return err
				}
				return bk.Put([]byte("lastupdate"), val)
			}); err != nil {
				log.Printf("ERROR in update data lastupdate: %s", err)
			} else {
				log.Printf("update data lastupdate!")
			}
		}
	}
	for {

		select {
		case <-tick.C:
			loopfunc()
		case <-start:
			loopfunc()
		}

	}
}
