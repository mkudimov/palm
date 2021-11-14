package readers

import (
	"io/ioutil"
	"os"
	"reflect"
	"runtime"
	"sync"
	"testing"
)

func TestFastReader_process(t *testing.T) {
	type args struct {
		buffer []byte
		store  *cmap
	}
	tests := []struct {
		name    string
		args    args
		wantMap map[string]int
	}{
		{
			name: "only_unique_tokens",
			args: args{
				buffer: []byte("FtyolcM\nyCmslXb\nznZsAoS"),
				store:  newCmap(),
			},
			wantMap: map[string]int{
				"FtyolcM": 1,
				"yCmslXb": 1,
				"znZsAoS": 1,
			},
		},
		{
			name: "only_unique_tokens",
			args: args{
				buffer: []byte("FtyolcM\nyCmslXb\nznZsAoS\nznZsAoS"),
				store:  newCmap(),
			},
			wantMap: map[string]int{
				"FtyolcM": 1,
				"yCmslXb": 1,
				"znZsAoS": 2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &FastReader{
				chunckSize: 1024 * 128,
				hooks:      []Hook{},
				numCPU:     runtime.NumCPU(),
				chunksPool: sync.Pool{New: func() interface{} {
					return make([]byte, 1024*128)
				}},
			}
			r.process(tt.args.buffer, tt.args.store)
			if !reflect.DeepEqual(tt.args.store.s, tt.wantMap) {
				t.Errorf("map content after FastReader.process() = %v, want %v", tt.args.store.s, tt.wantMap)
			}
		})
	}
}

func TestFastReader_Read(t *testing.T) {
	type fields struct {
		chunckSize int
		hooks      []Hook
		numCPU     int
		chunksPool sync.Pool
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
				chunckSize: 8,
				hooks:      []Hook{},
				numCPU:     runtime.NumCPU(),
				chunksPool: sync.Pool{New: func() interface{} {
					return make([]byte, 1024*128)
				}},
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
				chunckSize: 8,
				hooks:      []Hook{},
				numCPU:     runtime.NumCPU(),
				chunksPool: sync.Pool{New: func() interface{} {
					return make([]byte, 8)
				}},
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
			r := &FastReader{
				chunckSize: tt.fields.chunckSize,
				hooks:      tt.fields.hooks,
				numCPU:     tt.fields.numCPU,
				chunksPool: tt.fields.chunksPool,
			}
			got, err := r.Read(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("FastReader.Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FastReader.Read() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFastReader_read(t *testing.T) {
	file, err := ioutil.TempFile("", "*.txt")
	if err != nil {
		t.Errorf("FastReader.read() error creating temp file: %v", err)
	}
	defer os.Remove(file.Name())
	err = ioutil.WriteFile(file.Name(), []byte("testred\neyes\nmusic"), 0644)
	if err != nil {
		t.Errorf("FastReader.read() error creating temp file: %v", err)
	}
	br := newBufioReader(file)
	type fields struct {
		chunckSize int
		hooks      []Hook
		numCPU     int
		chunksPool sync.Pool
	}
	type args struct {
		br *bufioReader
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "test_until_new_line",
			fields: fields{
				chunckSize: 10,
				hooks:      []Hook{},
				numCPU:     runtime.NumCPU(),
				chunksPool: sync.Pool{New: func() interface{} {
					return make([]byte, 10)
				}},
			},
			args: args{
				br: br,
			},
			want:    []byte("testred\neyes\n"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &FastReader{
				chunckSize: tt.fields.chunckSize,
				hooks:      tt.fields.hooks,
				numCPU:     tt.fields.numCPU,
				chunksPool: tt.fields.chunksPool,
			}
			got, err := r.read(tt.args.br)
			if (err != nil) != tt.wantErr {
				t.Errorf("FastReader.read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FastReader.read() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFastReader_AppendHook(t *testing.T) {
	type fields struct {
		chunckSize int
		hooks      []Hook
		numCPU     int
		chunksPool sync.Pool
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
				chunckSize: 1024 * 128,
				hooks:      []Hook{},
				numCPU:     runtime.NumCPU(),
				chunksPool: sync.Pool{New: func() interface{} {
					return make([]byte, 1024*128)
				}},
			},
			args: args{
				h: nil,
			},
			wantHookCount: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &FastReader{
				chunckSize: tt.fields.chunckSize,
				hooks:      tt.fields.hooks,
				numCPU:     tt.fields.numCPU,
				chunksPool: tt.fields.chunksPool,
			}
			if l := len(r.AppendHook(tt.args.h).hooks); l != tt.wantHookCount {
				t.Errorf("hooks count after calling FastReader.AppendHook() = %v, want %v", l, tt.wantHookCount)
			}
		})
	}
}

func TestFastReader_SetChunckSize(t *testing.T) {
	type fields struct {
		chunckSize int
		hooks      []Hook
		numCPU     int
		chunksPool sync.Pool
	}
	type args struct {
		sz int
	}
	tests := []struct {
		name             string
		fields           fields
		args             args
		wantChunkSize    int
		wantLenChunkPool int
	}{
		{
			name: "name",
			fields: fields{
				chunckSize: 1024 * 128,
				hooks:      []Hook{},
				numCPU:     runtime.NumCPU(),
				chunksPool: sync.Pool{New: func() interface{} {
					return make([]byte, 1024*128)
				}},
			},
			args: args{
				sz: 128,
			},
			wantChunkSize:    128,
			wantLenChunkPool: 128,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &FastReader{
				chunckSize: tt.fields.chunckSize,
				hooks:      tt.fields.hooks,
				numCPU:     tt.fields.numCPU,
				chunksPool: tt.fields.chunksPool,
			}
			got := r.SetChunckSize(tt.args.sz)
			if got.chunckSize != tt.wantChunkSize {
				t.Errorf("chunk size after FastReader.SetChunckSize() = %v, want %v", got.chunckSize, tt.wantChunkSize)
			}
			if item := got.chunksPool.Get().([]byte); len(item) != tt.wantLenChunkPool {
				t.Errorf("size of item of chunk pool after FastReader.SetChunckSize() = %v, want %v", len(item), tt.wantLenChunkPool)
			}
		})
	}
}

func BenchmarkFastReader_Read(b *testing.B) {
	hook, err := NewPostgreSQLHook(dbConnection)
	if err != nil {
		b.Errorf("could not create NewPostgreSQLHook error = %v", err)
	}
	type fields struct {
		chunckSize int
		hooks      []Hook
		numCPU     int
		chunksPool sync.Pool
	}
	type args struct {
		path string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "bench_b1",
			fields: fields{
				chunckSize: 1024 * 128,
				hooks:      []Hook{},
				numCPU:     runtime.NumCPU(),
				chunksPool: sync.Pool{New: func() interface{} {
					return make([]byte, 1024*128)
				}},
			},
			args: args{
				path: b1,
			},
		},
		{
			name: "bench_b1_with_hooks",
			fields: fields{
				chunckSize: 1024 * 128,
				hooks:      []Hook{hook},
				numCPU:     runtime.NumCPU(),
				chunksPool: sync.Pool{New: func() interface{} {
					return make([]byte, 1024*128)
				}},
			},
			args: args{
				path: b1,
			},
		},
	}
	for _, bb := range tests {
		b.Run(bb.name, func(b *testing.B) {
			r := &FastReader{
				chunckSize: bb.fields.chunckSize,
				hooks:      bb.fields.hooks,
				numCPU:     bb.fields.numCPU,
				chunksPool: bb.fields.chunksPool,
			}
			for i := 0; i < b.N; i++ {
				r.Read(bb.args.path)
			}
		})
	}
}
