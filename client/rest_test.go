package main

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/dumacp/updatefs/loader"
)

func TestNewRequestFilesByDeviname(t *testing.T) {
	type args struct {
		client     *http.Client
		urlin      string
		devicename string
		date       int
		limit      int
		skip       int
	}
	tests := []struct {
		name    string
		args    args
		want    *[]*loader.FileData
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			"test1",
			args{
				&http.Client{},
				"http://127.0.0.1:8000",
				"all",
				0,
				2,
				1,
			},
			nil,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewRequestFilesByDevicename(tt.args.client, tt.args.urlin, tt.args.devicename, tt.args.date, tt.args.limit, tt.args.skip)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewRequestFilesByDeviname() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRequestFilesByDeviname() = %v, want %v", got, tt.want)
				for _, v := range got {
					t.Errorf("NewRequestFiles() = %+v", v)
				}
			}
		})
	}
}
