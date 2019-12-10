package loader

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
)

//FileData info about a file
type FileData struct {
	Name       string   `json:"name"`
	DeviceName []string `json:"devicename"`
	Md5        string   `json:"md5"`
	Date       int64    `json:"date"`
}

//LoadData Loda initial metadata
func LoadData(dir string) *[]*FileData {

	ret := make([]*FileData, 0, 0)

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Println(err)
		return nil
	}
	findfiles := func(files)
	for _, filename := range files {
		el := &FileData{}
		if filename.IsDir() {
			dataTemp := LoadData(fmt.Sprintf("%s/%s", dir, filename.Name()))
			for _, v := range *dataTemp {
				v.DeviceName = append(v.DeviceName, filename.Name())
			}
			return dataTemp
		}
		el.Name = filename.Name()
		el.Date = filename.ModTime().Unix()

		pathfile := fmt.Sprintf("%s/%s", dir, filename.Name())

		content, err := ioutil.ReadFile(pathfile)
		if err != nil {
			return nil
		}

		md5sum := md5.Sum(content)
		el.Md5 = hex.EncodeToString(md5sum[0:])

		ret = append(ret, el)
	}
	return &ret
}
