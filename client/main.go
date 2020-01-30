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
	"regexp"
	"strings"
	"syscall"
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

	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalf("Error: there is not hostname! %s", err)
	}
	if v, ok := envdev["sn-dev"]; ok {
		// if len(v) > 0 && v[0] == '"' {
		// 	v = v[1:]
		// }
		// if len(v) > 0 && v[len(v)-1] == '"' {
		// 	v = v[:len(v)-1]
		// }
		// if len(v) > 0 && v[0] != '"' {
		// 	hostname = v
		// }
		reg, err := regexp.Compile("[^a-zA-Z0-9\\-_\\.]+")
		if err != nil {
			log.Println(err)
		}
		processdString := reg.ReplaceAllString(v, "")
		log.Println(processdString)
		if len(processdString) > 0 {
			hostname = processdString
		}
	}
	log.Printf("hostname: %s", hostname)

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
		token, err := keycloakNewToken(hostname)
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

		store, err := NewRequestFilesByDevicename(client, urlin, hostname, int(filedata.Date), 1, 0)
		if err != nil {
			log.Printf("ERROR NewRequestFilesByDevicename: %s", err)
			return
		}
		if store == nil || len(store) <= 0 {
			if len(groupname) > 0 {
				store, err = NewRequestFilesByDevicename(client, urlin, groupname, int(filedata.Date), 1, 0)
				if err != nil {
					log.Printf("ERROR NewRequestFilesByDevicename all: %s", err)
					return
				}
				if store == nil || len(store) <= 0 {
					store, err = NewRequestFilesByDevicename(client, urlin, "all", int(filedata.Date), 1, 0)
					if err != nil {
						log.Printf("ERROR NewRequestFilesByDevicename all: %s", err)
						return
					}
					if store == nil || len(store) <= 0 {
						return
					}
				}
			}
		}

		var lastFiledata *loader.FileData
		for _, v := range store {
			filedatanow := v
			if !filedatanow.Override && lastFiledata != nil {
				break
			}
			lastFiledata = v
			fmt.Printf("%+v, %+v\n", filedatanow, filedata)
			if filedatanow.Date > filedata.Date && (filedatanow.Override || filedatanow.Ref > filedata.Ref) {
				if len(filedatanow.Md5) > 0 &&
					(!strings.Contains(filedatanow.Md5, filedata.Md5) || filedatanow.Override) {

					if _, err := os.Stat(pathupdatefile); err == nil {
						if !filedatanow.Override {
							log.Print("ERROR old pathupdatefile exits")
							return
						}
						log.Printf("override pathupdatefile!")
					}

					fileurl := fmt.Sprintf("%s/%s/%s", urlin, fileserverdir, filedatanow.FilePath)
					var errorDownload error
					for range []int{1, 2, 3} {
						md5sum, err := DownloadFile(client, fileurl, pathupdatefile)
						if err != nil {
							log.Printf("ERROR DownloadFile: %s", err)
							errorDownload = err
							continue
						}
						if !strings.Contains(filedatanow.Md5, md5sum) {
							log.Println("ERROR DownloadFile: md5 failed")
							errorDownload = err
							continue
						}
						errorDownload = nil
						break
					}
					if errorDownload != nil {
						log.Printf("ERROR DownloadFile: %s", err)
						return
					}
					log.Printf("UPDATE FILE DOWNLOAD: %+v", filedatanow)
				}

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
					if filedatanow.ForceReboot {
						log.Printf("force reboot!")
						syscall.Sync()
						syscall.Reboot(syscall.LINUX_REBOOT_CMD_RESTART)
						os.Exit(0)
					}

				}
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
