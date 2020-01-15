package updatedata

//UpdateStore database with of info about updates
type UpdateStore interface {
	Initialize(pathdb string) error
	NewUpdateDataDevice(devicename, key []byte, value *Updatedatadevice) error
	NewUpdateDataFile(filemd5, key []byte, value *Updatedatafile) error
	GetUpdateDataDevice(devicename, key []byte) (*Updatedatadevice, error)
	GetUpdateDataFile(filemd5, key []byte) (*Updatedatafile, error)
	SearchUpdateDataDevice(key []byte, date, limit, skip int) (*[]*Updatedatadevice, error)
	SearchUpdateDataFile(key []byte, date, limit, skip int) (*[]*Updatedatafile, error)
	GetLastDataDevices() *[]*Updatedatadevice
}
