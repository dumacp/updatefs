package datastore

import (
	"reflect"
	"testing"

	"github.com/dumacp/updatefs/loader"
)

func TestFiles_Initialize(t *testing.T) {
	type fields struct {
		Store *[]*loader.FileData
	}
	type args struct {
		dir string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
		{
			"test1",
			fields{},
			args{
				"/tmp/testloader",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Files{
				Store: tt.fields.Store,
			}
			b.Initialize(tt.args.dir)
			for _, v := range *b.Store {
				t.Log(v)
			}
		})
	}
}

func TestFiles_SearchDeviceName(t *testing.T) {
	type fields struct {
		Store *[]*loader.FileData
	}
	type args struct {
		dir        string
		devicename string
		date       int
		limit      int
		skip       int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *[]*loader.FileData
	}{
		// TODO: Add test cases.
		{
			"test1",
			fields{},
			args{
				"/tmp/testloader",
				"dir1",
				0,
				10,
				0,
			},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Files{
				Store: tt.fields.Store,
			}
			b.Initialize(tt.args.dir)
			if got := b.SearchDeviceName(tt.args.devicename, tt.args.date, tt.args.limit, tt.args.skip); !reflect.DeepEqual(got, tt.want) {
				for _, v := range *got {
					t.Logf("Files.SearchDeviceName() = %v", v)
				}
			}
		})
	}
}

func TestFiles_SearchID(t *testing.T) {
	type fields struct {
		Store *[]*loader.FileData
	}
	type args struct {
		id string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *loader.FileData
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Files{
				Store: tt.fields.Store,
			}
			if got := b.SearchID(tt.args.id); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Files.SearchID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFiles_Createfile(t *testing.T) {
	type fields struct {
		Store *[]*loader.FileData
	}
	type args struct {
		file *loader.FileData
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Files{
				Store: tt.fields.Store,
			}
			if got := b.CreateFile(tt.args.file); got != tt.want {
				t.Errorf("Files.Createfile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFiles_UpdateFile(t *testing.T) {
	type fields struct {
		Store *[]*loader.FileData
	}
	type args struct {
		id   string
		book *loader.FileData
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Files{
				Store: tt.fields.Store,
			}
			if got := b.UpdateFile(tt.args.id, tt.args.book); got != tt.want {
				t.Errorf("Files.UpdateFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFiles_DeleteFile(t *testing.T) {
	type fields struct {
		Store *[]*loader.FileData
	}
	type args struct {
		id string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Files{
				Store: tt.fields.Store,
			}
			if got := b.DeleteFile(tt.args.id); got != tt.want {
				t.Errorf("Files.DeleteFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilter(t *testing.T) {
	type args struct {
		vs *[]*loader.FileData
		f  func(*loader.FileData) bool
	}
	tests := []struct {
		name string
		args args
		want *[]*loader.FileData
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Filter(tt.args.vs, tt.args.f); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Filter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFiles_AllData(t *testing.T) {
	type fields struct {
		Store *[]*loader.FileData
	}
	tests := []struct {
		name   string
		fields fields
		want   *[]*loader.FileData
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Files{
				Store: tt.fields.Store,
			}
			if got := b.AllData(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Files.AllData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFiles_SearchMD5(t *testing.T) {
	type fields struct {
		Store *[]*loader.FileData
	}
	type args struct {
		md5sum string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *loader.FileData
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Files{
				Store: tt.fields.Store,
			}
			if got := b.SearchMD5(tt.args.md5sum); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Files.SearchMD5() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFiles_CreateFile(t *testing.T) {
	type fields struct {
		Store *[]*loader.FileData
	}
	type args struct {
		file *loader.FileData
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Files{
				Store: tt.fields.Store,
			}
			if got := b.CreateFile(tt.args.file); got != tt.want {
				t.Errorf("Files.CreateFile() = %v, want %v", got, tt.want)
			}
		})
	}
}
