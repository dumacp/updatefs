package loader

import (
	"encoding/json"
	"github.com/boltdb/bolt"
	"reflect"
	"testing"
)

func TestLoadData(t *testing.T) {
	type args struct {
		db *bolt.DB
	}
	tests := []struct {
		name string
		args args
		want *[]*FileData
	}{
		{
			// TODO: Add test cases.
			"test1",
			args{
				nil,
			},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LoadData(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadData() = %v, want %v", got, tt.want)
				for _, v := range *got {
					t.Errorf("data: %v", v)
				}
			}
		})
	}
}

func TestData(t *testing.T) {
	type args struct {
		data1 []byte
	}
	tests := []struct {
		name string
		args args
		want *FileData
	}{
		{
			// TODO: Add test cases.
			"test1",
			args{
				[]byte(`{
					"id": "1",
					"name": "Prueba",
					"devicename": ["testdev1"],
					"md5": "12345677788",
					"date": 0,
					"filepath": "/tmp/all",
					"desc": "",
					"ref": "v1"
				}`),
			},
			&FileData{
				ID:         "1",
				Name:       "Prueba",
				DeviceName: []string{"testdev1"},
				Md5:        "12345677788",
				Date:       0,
				FilePath:   "/tmp/all",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var1 := new(FileData)
			if err := json.Unmarshal(tt.args.data1, var1); err != nil {
				t.Fatalf("error: %v", err)
			}
			if got := var1; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UnMarshall = %v, want %v", got, tt.want)
			} else {
				t.Errorf("UnMarshall = %v", got)
			}
		})
	}
}
