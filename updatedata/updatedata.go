package updatedata

import (
	"encoding/json"
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/dumacp/updatefs/loader"
)

//UpdateData type updatedata
type UpdateData struct {
	db *bolt.DB
}

//Updatedatadevice type info about of update in devices
type Updatedatadevice struct {
	ID        string           `json:"id"`
	Name      string           `json:Name`
	Date      int              `json:"date"`
	Filedata  *loader.FileData `json:"filedata"`
	IPRequest string           `json:"iprequest"`
}

//Updatedatafile type info about of update in files
type Updatedatafile struct {
	ID        string `json:"id"`
	Name      string `json:"Name"`
	Date      int    `json:"date"`
	Device    string `json:"device"`
	IPRequest string `json:"iprequest"`
}

const (
	bucketupdates            = "updatedata"
	bucketupdatesdevices     = "updatesdevices"
	bucketupdatesfiles       = "updatesfiles"
	bucketupdatesdevicesDate = "updatesdevicesDates"
	bucketupdatesfilesDate   = "updatesfilesDates"
)

//Initialize Open Database with buckets
func (up *UpdateData) Initialize(pathdb string) error {
	var err error
	up.db, err = bolt.Open(pathdb, 0644, nil)
	if err != nil {
		return err
	}
	return nil
}

//NewUpdateDataDevice create new entry for "bucketupdatesdevices"
func (up *UpdateData) NewUpdateDataDevice(devicename, key []byte, value *Updatedatadevice) error {

	val1, err := json.Marshal(value)
	if err != nil {
		return err
	}
	if err := up.db.Update(func(tx *bolt.Tx) error {
		bk, err := tx.CreateBucketIfNotExists([]byte(bucketupdatesdevices))
		if err != nil {
			return err
		}
		bkDate, err := tx.CreateBucketIfNotExists([]byte(bucketupdatesdevicesDate))
		if err != nil {
			return err
		}
		bkdevices, err := bk.CreateBucketIfNotExists(devicename)
		if err != nil {
			return err
		}
		bkdevicesDate, err := bkDate.CreateBucketIfNotExists(devicename)
		if err != nil {
			return err
		}
		if err := bkdevices.Put(key, val1); err != nil {
			return err
		}
		if err := bkdevicesDate.Put([]byte(fmt.Sprintf("%v", value.Date)), val1); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

//NewUpdateDataFile create new entry data for "bucketupdatesfiles"
func (up *UpdateData) NewUpdateDataFile(filemd5, key []byte, value *Updatedatafile) error {
	val1, err := json.Marshal(value)
	if err != nil {
		return err
	}
	if err := up.db.Update(func(tx *bolt.Tx) error {
		bk, err := tx.CreateBucketIfNotExists([]byte(bucketupdatesfiles))
		if err != nil {
			return err
		}
		bkDate, err := tx.CreateBucketIfNotExists([]byte(bucketupdatesfilesDate))
		if err != nil {
			return err
		}
		bkfiles, err := bk.CreateBucketIfNotExists(filemd5)
		if err != nil {
			return err
		}
		bkfilesDate, err := bkDate.CreateBucketIfNotExists(filemd5)
		if err != nil {
			return err
		}
		if err := bkfiles.Put(key, val1); err != nil {
			return err
		}
		if err := bkfilesDate.Put([]byte(fmt.Sprintf("%v", value.Date)), val1); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

//GetUpdateDataDevice get value for key in bucket bucketupdatesdevices
func (up *UpdateData) GetUpdateDataDevice(devicename, key []byte) (*Updatedatadevice, error) {
	var deviced *Updatedatadevice
	if err := up.db.View(func(tx *bolt.Tx) error {
		bk := tx.Bucket([]byte(bucketupdatesdevices))
		if bk == nil {
			return bolt.ErrBucketNotFound
		}
		bkdevices := bk.Bucket(devicename)
		if bkdevices == nil {
			return bolt.ErrBucketNotFound
		}
		value := bkdevices.Get(key)
		if value == nil {
			return nil
		}
		deviced = new(Updatedatadevice)
		if err := json.Unmarshal(value, deviced); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return deviced, nil
}

//GetLastDataDevices get value for key in bucket bucketupdatesdevices
func (up *UpdateData) GetLastDataDevices() *[]*Updatedatadevice {
	updatesData := make([]*Updatedatadevice, 0)
	if err := up.db.View(func(tx *bolt.Tx) error {
		bk := tx.Bucket([]byte(bucketupdatesdevicesDate))
		it := bk.Cursor()
		ki, vi := it.First()
		if ki == nil {
			return nil
		}
		appendFunc := func(key, value []byte) error {
			if value != nil {
				return nil
			}
			bkDevice := tx.Bucket([]byte(key))
			itlast := bkDevice.Cursor()
			_, vii := itlast.Last()
			deviced := new(Updatedatadevice)
			if err := json.Unmarshal(vii, deviced); err != nil {
				return err
			}
			return nil
		}
		if err := appendFunc(ki, vi); err != nil {
			return err
		}
		for {
			if kii, vii := it.Prev(); kii != nil {
				if err := appendFunc(kii, vii); err != nil {
					return err
				}
				continue
			}
			break
		}

		return nil
	}); err != nil {
		return nil
	}
	return &updatesData
}

//GetUpdateDataFile get value for key in bucket bucketupdatesfiles
func (up *UpdateData) GetUpdateDataFile(filemd5, key []byte) (*Updatedatafile, error) {
	var filed *Updatedatafile
	if err := up.db.View(func(tx *bolt.Tx) error {
		bk := tx.Bucket([]byte(bucketupdatesfiles))
		if bk == nil {
			return bolt.ErrBucketNotFound
		}
		bkfiles := bk.Bucket(filemd5)
		if bkfiles == nil {
			return bolt.ErrBucketNotFound
		}
		value := bkfiles.Get(key)
		if value == nil {
			return nil
		}
		filed = new(Updatedatafile)
		if err := json.Unmarshal(value, filed); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return filed, nil
}

//SearchUpdateDataDevice search value parameters in bucket bucketupdatesdevicesDate
func (up *UpdateData) SearchUpdateDataDevice(devicename []byte, date, limit, skip int) (*[]*Updatedatadevice, error) {
	updatesData := make([]*Updatedatadevice, 0)
	if err := up.db.View(func(tx *bolt.Tx) error {
		bk := tx.Bucket([]byte(bucketupdatesdevicesDate))
		if bk == nil {
			return bolt.ErrBucketNotFound
		}
		bkdevices := bk.Bucket(devicename)
		if bkdevices == nil {
			return bolt.ErrBucketNotFound
		}

		it := bkdevices.Cursor()
		ki, vi := it.Last()
		if ki == nil {
			return nil
		}
		appendFunc := func(value []byte) error {
			deviced := new(Updatedatadevice)
			if err := json.Unmarshal(value, deviced); err != nil {
				return err
			}
			if deviced.Date < date {
				return nil
			}
			updatesData = append(updatesData, deviced)
			return nil
		}
		if err := appendFunc(vi); err != nil {
			return err
		}
		for {
			if kii, vii := it.Prev(); kii != nil {
				if err := appendFunc(vii); err != nil {
					return err
				}
				continue
			}
			break
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return &updatesData, nil
}

//SearchUpdateDataFile search value for parameters in bucket bucketupdatesfilesDate
func (up *UpdateData) SearchUpdateDataFile(filemd5 []byte, date, limit, skip int) (*[]*Updatedatafile, error) {
	updatesData := make([]*Updatedatafile, 0)
	if err := up.db.View(func(tx *bolt.Tx) error {
		bk := tx.Bucket([]byte(bucketupdatesfilesDate))
		if bk == nil {
			return bolt.ErrBucketNotFound
		}
		bkdevices := bk.Bucket(filemd5)
		if bkdevices == nil {
			return bolt.ErrBucketNotFound
		}

		it := bkdevices.Cursor()
		ki, vi := it.Last()
		if ki == nil {
			return nil
		}
		appendFunc := func(value []byte) error {
			deviced := new(Updatedatafile)
			if err := json.Unmarshal(value, deviced); err != nil {
				return err
			}
			if deviced.Date < date {
				return nil
			}
			updatesData = append(updatesData, deviced)
			return nil
		}
		if err := appendFunc(vi); err != nil {
			return err
		}
		for {
			if kii, vii := it.Prev(); kii != nil {
				if err := appendFunc(vii); err != nil {
					return err
				}
				continue
			}
			break
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return &updatesData, nil
}
