package loader

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/google/uuid"
)

//FileData info about a file
type FileData struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	DeviceName []string `json:"devicename"`
	Md5        string   `json:"md5"`
	Date       int64    `json:"date"`
	FilePath   string   `json:"filepath"`
}

//LoadData Loda initial metadata
func LoadData(dir string) *[]*FileData {

	ret := make([]*FileData, 0, 0)

	var findfiles func(string) *[]*FileData
	findfiles = func(dirt string) *[]*FileData {
		reti := make([]*FileData, 0, 0)
		files, err := ioutil.ReadDir(dirt)
		if err != nil {
			log.Println(err)
			return nil
		}
		for _, filename := range files {
			// fmt.Printf("filename: %v\n", filename.Name())
			el := &FileData{}
			if filename.IsDir() {
				// fmt.Printf("dir recursive: %s/%s\n", dirt, filename.Name())
				elt := findfiles(fmt.Sprintf("%s/%s", dirt, filename.Name()))
				// fmt.Printf("dir recursive: %+v\n", *elt)
				if elt != nil {
					for _, v := range *elt {
						// fmt.Printf("dir add: %+v, name: %v\n", filename.Name(), v.Name)
						v.DeviceName = append(v.DeviceName, filename.Name())
						// v.FilePath = fmt.Sprintf("%s/%s", filename.Name(), v.Name)
					}
					reti = append(reti, *elt...)
				}
				continue
			}
			el.ID = uuid.New().String()
			el.Name = filename.Name()
			el.Date = filename.ModTime().Unix()

			pathfile := fmt.Sprintf("%s/%s", dirt, filename.Name())
			el.FilePath, _ = filepath.Rel(dir, pathfile)

			content, err := ioutil.ReadFile(pathfile)
			if err != nil {
				continue
			}

			md5sum := md5.Sum(content)
			el.Md5 = hex.EncodeToString(md5sum[0:])

			reti = append(reti, el)
		}
		return &reti
	}
	retii := findfiles(dir)
	if retii != nil {
		ret = append(ret, *retii...)
		for _, v := range ret {
			v.DeviceName = append(v.DeviceName, filepath.Base(dir))
			// v.FilePath = fmt.Sprintf("%s/%s", filepath.Base(dir), v.Name)
		}

	}

	return &ret
}
