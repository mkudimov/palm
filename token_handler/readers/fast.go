package readers

import (
	"io"
	"os"
	"runtime"
	"strings"
	"sync"

	"github.com/rs/zerolog/log"
)

type (
	FastReader struct {
		chunckSize int
		hooks      []Hook
		numCPU     int
		chunksPool sync.Pool
	}
)

func NewFastReader() *FastReader {
	chunkSize := 1024 * 1024
	return &FastReader{
		chunckSize: chunkSize,
		hooks:      []Hook{},
		numCPU:     runtime.NumCPU(),
		chunksPool: sync.Pool{New: func() interface{} {
			return make([]byte, chunkSize)
		}},
	}
}

func (r *FastReader) SetChunckSize(sz int) *FastReader {
	r.chunckSize = sz
	r.chunksPool = sync.Pool{New: func() interface{} {
		return make([]byte, sz)
	}}
	return r
}

func (r *FastReader) AppendHook(h Hook) *FastReader {
	r.hooks = append(r.hooks, h)
	return r
}

func (r *FastReader) Read(path string) (map[string]int, error) {
	file, err := os.Open(path)
	if err != nil {
		log.Error().Err(err).Str("file", path).Msg("couldn't open file for reading")
		return nil, err
	}
	defer file.Close()
	m := newCmap()
	errChan := make(chan error, r.numCPU)
	rd := newBufioReader(file)
	var wg sync.WaitGroup
	wg.Add(r.numCPU)
	for i := 0; i < r.numCPU; i++ {
		go r.readAndProcess(rd, m, &wg, errChan)
	}
	wg.Wait()
	if err := <-errChan; err != nil {
		log.Error().Err(err).Str("file", path).Msg("error happened during file handling")
		return nil, err
	}
	close(errChan)
	return m.s, err
}

func (r *FastReader) readAndProcess(br *bufioReader, store *cmap, wg *sync.WaitGroup, errChan chan<- error) {
	for {
		buffer, err := r.read(br)
		if err != nil {
			if err != io.EOF {
				log.Error().Err(err).Msg("error when reading file")
				errChan <- err
			}
			break
		}
		r.process(buffer, store)
	}
	errChan <- nil
	wg.Done()
}

func (r *FastReader) read(br *bufioReader) ([]byte, error) {
	buffer := r.chunksPool.Get().([]byte)
	readBytes, err := br.Read(buffer)
	if err != nil {
		return nil, err
	}
	buffer = buffer[:readBytes]
	untilNewLine, err := br.ReadBytes('\n')
	if err != io.EOF {
		buffer = append(buffer, untilNewLine...)
	}
	return buffer, nil
}

func (r *FastReader) process(buffer []byte, store *cmap) {
	tokens := strings.Split(string(buffer), "\n")
	if tokens[len(tokens)-1] == "" {
		tokens = tokens[:len(tokens)-1]
	}
	for _, vv := range tokens {
		if _, ok := store.get(vv); !ok {
			for _, h := range r.hooks {
				h.Run(vv)
			}
		}
		store.increase(vv)
	}
}
