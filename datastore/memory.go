package datastore

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/boltdb/bolt"
	"github.com/dumacp/updatefs/loader"
)

//Files struct abour store
type Files struct {
	Store *[]*loader.FileData `json:"store"`
	db    *bolt.DB
}

//Initialize init database
func (b *Files) Initialize(pathdb string) {

	var err error
	b.db, err = bolt.Open(pathdb, 0644, nil)
	if err != nil {
		log.Print("error: dont open files database")
	}

	b.Store = loader.LoadData(b.db)
	for i, v := range *b.Store {
		log.Printf("file %d: %v", i, v)
	}
}

func (b *Files) SearchDeviceName(devicename string, date, limit, skip int) *[]*loader.FileData {
	ret := Filter(b.Store, func(v *loader.FileData) bool {
		for _, vi := range v.DeviceName {
			if strings.Contains(strings.ToLower(vi), strings.ToLower(devicename)) && int(v.Date) > date {
				return true
			}
		}
		return false
	})
	if limit == 0 || limit > len(*ret) {
		limit = len(*ret)
	}
	data := (*ret)[skip:limit]
	return &data
}

func (b *Files) SearchUpdate(devicename string, date, limit, skip int) *[]*loader.FileData {
	ret := Filter(b.Store, func(v *loader.FileData) bool {
		for _, vi := range v.DeviceName {
			if strings.Contains(strings.ToLower(vi), strings.ToLower(devicename)) && int(v.Date) > date {
				return true
			}
		}
		return false
	})
	if limit == 0 || limit > len(*ret) {
		limit = len(*ret)
	}
	data := (*ret)[skip:limit]
	return &data
}

func (b *Files) AllData() *[]*loader.FileData {
	return b.Store
}

func (b *Files) SearchID(id string) *loader.FileData {
	ret := Filter(b.Store, func(v *loader.FileData) bool {
		return strings.ToLower(v.ID) == strings.ToLower(id)
	})
	if len(*ret) > 0 {
		return (*ret)[0]
	}
	return nil
}

func (b *Files) SearchMD5(md5sum string) *loader.FileData {
	ret := Filter(b.Store, func(v *loader.FileData) bool {
		return strings.ToLower(v.Md5) == strings.ToLower(md5sum)
	})
	if len(*ret) > 0 {
		return (*ret)[0]
	}
	return nil
}

func (b *Files) CreateFile(file *loader.FileData) bool {
	filev, err := json.Marshal(file)
	if err != nil {
		log.Print("error: dont parse file")
		return false
	}

	if err := b.db.Update(func(tx *bolt.Tx) error {
		bk, err := tx.CreateBucketIfNotExists([]byte(loader.Bucketfiles))
		if err != nil {
			return bolt.ErrBucketNotFound
		}
		if err := bk.Put([]byte(file.Md5), filev); err != nil {
			return err
		}
		return nil
	}); err != nil {
		log.Printf("error: create file data in database, %s", err)
		return false
	}

	*b.Store = append(*b.Store, file)
	return true
}

func (b *Files) UpdateFile(id string, book *loader.FileData) bool {
	return false
}

func (b *Files) DeleteFile(id string) bool {
	return false
}

/**
func (b *Books) DeleteBook(isbn string) bool {
	indexToDelete := -1
	for i, v := range *b.Store {
		if v.ISBN == isbn {
			indexToDelete = i
			break
		}
	}
	if indexToDelete >= 0 {
		(*b.Store)[indexToDelete], (*b.Store)[len(*b.Store)-1] = (*b.Store)[len(*b.Store)-1], (*b.Store)[indexToDelete]
		*b.Store = (*b.Store)[:len(*b.Store)-1]
		return true
	}
	return false
}

func (b *Books) UpdateBook(isbn string, book *loader.BookData) bool {
	for _, v := range *b.Store {
		if v.ISBN == isbn {
			v = book
			return true
		}
	}
	return false
}
/**/

func Filter(vs *[]*loader.FileData, f func(*loader.FileData) bool) *[]*loader.FileData {
	vsf := make([]*loader.FileData, 0)
	for _, v := range *vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return &vsf
}
