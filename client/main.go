package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
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
	dirDB           = "/SD/boltdbs"
	nameDB          = "updatefs"
	fileserverdir   = "updatevoc/static"
	pathupdatefile  = "/SD/update/migracion.zip"
	pathfirmwareRef = "/usr/include/firmware-ne"
	pathenvfile     = "/usr/include/serial-dev"
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

	envdev := make(map[string]string)
	if fileenv, err := os.Open(pathenvfile); err != nil {
		log.Printf("error: reading file env, %s", err)
	} else {
		scanner := bufio.NewScanner(fileenv)
		for scanner.Scan() {
			line := scanner.Text()
			log.Println(line)
			split := strings.Split(line, "=")
			if len(split) > 1 {
				envdev[split[0]] = split[1]
			}
		}
	}

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

	// hostname, err := os.Hostname()
	// if err != nil {
	// 	log.Fatalf("Error: there is not hostname! %s", err)
	// }

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
	filedata.Ref = -1
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

	//TODO: is this necesary?
	// var refSystem int
	// if firmwareV, err := ioutil.ReadFile(pathfirmwareRef); err != nil {
	// 	refSystem = -1
	// } else {
	// 	if refSystem, err = strconv.Atoi(string(firmwareV)); err != nil {
	// 		refSystem = -1
	// 	}
	// }
	// if filedata.Ref < refSystem {
	// 	filedata.Ref = refSystem
	// }

	var groupname string
	var client *http.Client

	keycloakconn := func() error {
		defer func() {
			if r := recover(); r != nil {
				log.Println("Recovered in f", r)
			}
		}()
		if serverkey == nil {
			if err := keycloakinit(); err != nil {
				log.Printf("ERROR keycloak init: %s", err)
				return err
			}
		}
		token, err := keycloakNewToken()
		if err != nil {
			log.Printf("ERROR keycloak token request: %s", err)
			return err
		}
		ts := keycloakTokenSource(token)
		atts, err := keycloakinfo(ts)
		if err != nil {
			log.Printf("ERROR keycloak token request attribs: %s", err)
			return err
		}
		log.Printf("attrs: %+v", atts)
		if v, ok := atts["group_name"]; ok {
			groupname = fmt.Sprintf("%s", v)
		}
		client, err = keycloakclient(ts)
		if err != nil {
			return err
		}
		return nil
	}

	tick := time.NewTicker(10 * time.Minute)
	start := time.NewTimer(3 * time.Second)

	loopfunc := func(client *http.Client) {
		if client == nil {
			return
		}
		defer func() {
			if r := recover(); r != nil {
				log.Println("Recovered in f", r)
			}
			client.CloseIdleConnections()
		}()
		hostname, err := os.Hostname()
		if err != nil {
			log.Fatalf("Error: there is not hostname! %s", err)
		}
		if v, ok := envdev["sn-dev"]; ok {
			if len(v) > 0 {
				hostname = v
			}
		}

		store, err := NewRequestFilesByDevicename(client, urlin, hostname, int(filedata.Date), 1, 0)
		if err != nil {
			log.Printf("ERROR NewRequestFilesByDevicename: %s", err)
			return
		}
		if store == nil || len(*store) <= 0 {
			if len(groupname) > 0 {
				store, err = NewRequestFilesByDevicename(client, urlin, groupname, int(filedata.Date), 1, 0)
				if err != nil {
					log.Printf("ERROR NewRequestFilesByDevicename all: %s", err)
					return
				}
				if store == nil || len(*store) <= 0 {
					store, err = NewRequestFilesByDevicename(client, urlin, "all", int(filedata.Date), 1, 0)
					if err != nil {
						log.Printf("ERROR NewRequestFilesByDevicename all: %s", err)
						return
					}
					if store == nil || len(*store) <= 0 {
						return
					}
				}
			}
		}

		filedatanow := (*store)[0]
		fmt.Printf("%+v, %+v\n", filedatanow, filedata)
		if filedatanow.Date > filedata.Date &&
			filedatanow.Ref > filedata.Ref &&
			(len(filedata.Md5) <= 0 ||
				!strings.Contains(filedatanow.Md5, filedata.Md5)) {

			fileurl := fmt.Sprintf("%s/%s/%s", urlin, fileserverdir, filedatanow.FilePath)
			err := DownloadFile(fileurl, pathupdatefile)
			if err != nil {
				log.Printf("ERROR DownloadFile: %s", err)
				return
			}
			log.Printf("UPDATE FILE DOWNLOAD: %+v", filedatanow)

			if filedatanows, err := json.Marshal(filedatanow); err == nil {
				if err := NewUpdateByDevicename(
					client,
					urlin,
					hostname,
					filedatanow.Md5,
					string(filedatanows),
					int(time.Now().Unix()),
				); err != nil {
					log.Printf("error NewUpdateByDevicename: %s", err)
				}
			}
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
			// keycloakconn()
			loopfunc(client)
		case <-start.C:
			if err := keycloakconn(); err != nil {
				start.Reset(30 * time.Second)
				break
			}
			loopfunc(client)
		}

	}
}
