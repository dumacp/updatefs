package main

import (
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
)

func TestNewRequestFiles(t *testing.T) {
	type args struct {
		url        string
		devicename string
	}
	tests := []struct {
		name    string
		args    args
		want    *http.Response
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			"test1",
			args{
				"http://127.0.0.1:8000/api/v1/filedata",
				"dir1",
			},
			nil,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewRequestFiles(tt.args.url, tt.args.devicename)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewRequestFiles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRequestFiles() = %+v, url %v want %v", got, tt.args.url, tt.want)
				body, _ := ioutil.ReadAll(got.Body)
				got.Body.Close()
				t.Errorf("NewRequestFiles() = %s", body)
			}
		})
	}
}
