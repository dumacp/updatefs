package loader

import (
	"encoding/json"
	"log"
	"sort"

	"github.com/boltdb/bolt"
)

// //LoadData Loda initial metadata
// func LoadData(dir string) *[]*FileData {

// 	ret := make([]*FileData, 0, 0)

// 	var findfiles func(string) *[]*FileData
// 	findfiles = func(dirt string) *[]*FileData {
// 		reti := make([]*FileData, 0, 0)
// 		files, err := ioutil.ReadDir(dirt)
// 		if err != nil {
// 			log.Println(err)
// 			return nil
// 		}
// 		for _, filename := range files {
// 			// fmt.Printf("filename: %v\n", filename.Name())
// 			el := &FileData{}
// 			if filename.IsDir() {
// 				// fmt.Printf("dir recursive: %s/%s\n", dirt, filename.Name())
// 				elt := findfiles(fmt.Sprintf("%s/%s", dirt, filename.Name()))
// 				// fmt.Printf("dir recursive: %+v\n", *elt)
// 				if elt != nil {
// 					for _, v := range *elt {
// 						// fmt.Printf("dir add: %+v, name: %v\n", filename.Name(), v.Name)
// 						v.DeviceName = append(v.DeviceName, filename.Name())
// 						// v.FilePath = fmt.Sprintf("%s/%s", filename.Name(), v.Name)
// 					}
// 					reti = append(reti, *elt...)
// 				}
// 				continue
// 			}
// 			el.ID = uuid.New().String()
// 			el.Name = filename.Name()
// 			el.Date = filename.ModTime().Unix()

// 			pathfile := fmt.Sprintf("%s/%s", dirt, filename.Name())
// 			el.FilePath, _ = filepath.Rel(dir, pathfile)

// 			content, err := ioutil.ReadFile(pathfile)
// 			if err != nil {
// 				continue
// 			}

// 			md5sum := md5.Sum(content)
// 			el.Md5 = hex.EncodeToString(md5sum[0:])

// 			reti = append(reti, el)
// 		}
// 		return &reti
// 	}
// 	retii := findfiles(dir)
// 	if retii != nil {
// 		ret = append(ret, *retii...)
// 		for _, v := range ret {
// 			v.DeviceName = append(v.DeviceName, filepath.Base(dir))
// 			// v.FilePath = fmt.Sprintf("%s/%s", filepath.Base(dir), v.Name)
// 		}

// 	}

// 	return &ret
// }

//FileData info about a file
type FileData struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	DeviceName  []string `json:"devicename"`
	Md5         string   `json:"md5"`
	Date        int64    `json:"date"`
	FilePath    string   `json:"filepath"`
	Description string   `json:"desc"`
	Ref         int      `json:"ref"`
	Version     string   `json:"version"`
	ForceReboot bool     `json:"reboot"`
	ForceApply  bool     `json:"apply"`
	Override    bool     `json:"override"`
}

const (
	Bucketfiles = "migrationfiles"
)

//LoadData Loda initial metadata
func LoadData(db *bolt.DB) *[]*FileData {

	ret := make([]*FileData, 0)
	storeMap := make(map[int]*FileData)

	if err := db.View(func(tx *bolt.Tx) error {
		bk := tx.Bucket([]byte(Bucketfiles))
		if bk == nil {
			return bolt.ErrBucketNotFound
		}
		if err := bk.ForEach(func(k []byte, v []byte) error {
			if k != nil {
				filed := new(FileData)
				if err := json.Unmarshal(v, filed); err != nil {
					return err
				}
				log.Printf("filed: %+v", filed)
				storeMap[int(filed.Date)] = filed
				// ret = append(ret, filed)
			}
			return nil
		}); err != nil {
			log.Println(err)
			return err
		}

		return nil

	}); err != nil {
		log.Println(err)
	}

	keys := make([]int, 0, len(storeMap))
	for k := range storeMap {
		keys = append(keys, k)
	}

	sort.Ints(keys)

	for _, k := range keys {
		ret = append(ret, storeMap[k])
	}

	return &ret
}

// //CreateFile new fileData
// func CreateFile(db *bolt.DB, fileupload *multipart.File, desc, version, base, path string, ref int) *FileData {
// 	filePath := filepath.Clean(fmt.Sprintf("%s/%s/migracion_%s.zip", base, filepath.Clean(path), version))
// 	data, err := ioutil.ReadAll(*fileupload)
// 	if err != nil {
// 		if err := ioutil.WriteFile(filePath, data, 0644); err != nil {
// 			return nil
// 		}
// 	}
// 	filed := new(FileData)

// 	filed.FilePath = filePath
// 	filed.Name = filepath.Base(filePath)
// 	filed.ID = uuid.New().String()
// 	filed.Date = time.Now().Unix()
// 	md5sum := md5.Sum(data)
// 	filed.Md5 = hex.EncodeToString(md5sum[0:])
// 	filed.Description = desc
// 	filed.Ref = ref

// 	filev, err := json.Marshal(filed)
// 	if err != nil {
// 		log.Print("error: dont parse file")
// 		return nil
// 	}

// 	if err := db.Update(func(tx *bolt.Tx) error {
// 		bk := tx.Bucket([]byte(bucketfiles))
// 		if bk == nil {
// 			return bolt.ErrBucketNotFound
// 		}
// 		if err := bk.Put([]byte(filed.Md5), filev); err != nil {
// 			return err
// 		}
// 		return nil
// 	}); err != nil {
// 		return nil
// 	}

// 	return filed
// }
