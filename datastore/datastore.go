package datastore

import "github.com/dumacp/updatefs/loader"

//FileStore database with of info about files
type FileStore interface {
	Initialize(dir string)
	SearchMd5(md5 string, date, limit, skip int) *[]*loader.FileData
	SearchDeviceName(deviname string, date, limit, skip int) *[]*loader.FileData
	CreateFile(file *loader.FileData) bool
	DeleteFileByMD5(md5 string) bool
	DeleteFileByDeviceName(devicename string) bool
	Updatefile(md5 string, book *loader.FileData) bool
}
