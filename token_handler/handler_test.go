package tokenhandler

import (
	"io/ioutil"
	"os"
	"testing"
)

func Test_dump(t *testing.T) {
	file, err := ioutil.TempFile("", "*.txt")
	if err != nil {
		t.Errorf("dump() error creating temp file: %v", err)
	}
	defer os.Remove(file.Name())
	file1, err := ioutil.TempFile("", "*.txt")
	if err != nil {
		t.Errorf("dump() error creating temp file: %v", err)
	}
	type args struct {
		m   map[string]int
		out string
	}
	tests := []struct {
		name            string
		args            args
		wantErr         bool
		wantFileContent string
	}{
		{
			name: "only_unique_tokens",
			args: args{
				m: map[string]int{
					"test":  1,
					"test1": 1,
				},
				out: file.Name(),
			},
			wantErr:         false,
			wantFileContent: "{}",
		},
		{
			name: "not_only_unique_tokens",
			args: args{
				m: map[string]int{
					"test":  4,
					"test1": 1,
					"test2": 22,
					"test3": 1,
				},
				out: file1.Name(),
			},
			wantErr:         false,
			wantFileContent: `{"test":4,"test2":22}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := dump(tt.args.m, tt.args.out); (err != nil) != tt.wantErr {
				t.Errorf("dump() error = %v, wantErr %v", err, tt.wantErr)
			}
			content, err := ioutil.ReadFile(tt.args.out)
			if err != nil {
				t.Errorf("dump() error reading temp file: %v", err)
			}
			if string(content) != tt.wantFileContent {
				t.Errorf("dump() content = %v, wantFileContent %v", string(content), tt.wantFileContent)
			}
		})
	}
}
