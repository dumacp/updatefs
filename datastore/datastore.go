package datastore

import "github.com/dumacp/updatefs/loader"

//FileStore database with of info about files
type FileStore interface {
	Initialize(pathdb, pathfiles string)
	SearchID(id string) *loader.FileData
	SearchMD5(md5sum string) *loader.FileData
	AllData() *[]*loader.FileData
	SearchDeviceName(deviname string, date, limit, skip int) *[]*loader.FileData
	CreateFile(file *loader.FileData) bool
	DeleteFile(id string) bool
	UpdateFile(id string, file *loader.FileData) bool
}
