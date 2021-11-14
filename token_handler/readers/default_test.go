package readers

import (
	"reflect"
	"testing"
)

func TestDefaultReader_Read(t *testing.T) {
	type fields struct {
		hooks []Hook
	}
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[string]int
		wantErr bool
	}{
		{
			name: "non_existing_file",
			fields: fields{
				hooks: []Hook{},
			},
			args: args{
				path: ";)",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "non_existing_file",
			fields: fields{
				hooks: []Hook{},
			},
			args: args{
				path: t2,
			},
			want: map[string]int{
				"FtyolcM": 1,
				"yCmslXb": 1,
				"znZsAoS": 3,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &DefaultReader{
				hooks: tt.fields.hooks,
			}
			got, err := r.Read(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("DefaultReader.Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DefaultReader.Read() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultReader_AppendHook(t *testing.T) {
	type fields struct {
		hooks []Hook
	}
	type args struct {
		h Hook
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		wantHookCount int
	}{
		{
			name: "new_hook",
			fields: fields{
				hooks: []Hook{},
			},
			args: args{
				h: nil,
			},
			wantHookCount: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &DefaultReader{
				hooks: tt.fields.hooks,
			}
			if l := len(r.AppendHook(tt.args.h).hooks); l != tt.wantHookCount {
				t.Errorf("hooks count after calling DefaultReader.AppendHook() = %v, want %v", l, tt.wantHookCount)
			}
		})
	}
}
