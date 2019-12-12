package main

import "testing"

func TestDownloadFile(t *testing.T) {
	type args struct {
		url      string
		filepath string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			"test1",
			args{
				"http://127.0.0.1:8000/static/file0",
				"/tmp/fileDownload0",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DownloadFile(tt.args.url, tt.args.filepath); (err != nil) != tt.wantErr {
				t.Errorf("DownloadFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
