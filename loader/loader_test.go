package loader

import (
	"reflect"
	"testing"
)

func TestLoadData(t *testing.T) {
	type args struct {
		dir string
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
				"/tmp/all",
			},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LoadData(tt.args.dir); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadData() = %v, want %v", got, tt.want)
				for _, v := range *got {
					t.Errorf("data: %v", v)
				}
			}
		})
	}
}
