package readers

import (
	"bufio"
	"io"
	"sync"
)

type bufioReader struct {
	m sync.Mutex
	r *bufio.Reader
}

func newBufioReader(rd io.Reader) *bufioReader {
	return &bufioReader{
		r: bufio.NewReader(rd),
	}
}

func (r *bufioReader) Read(p []byte) (n int, err error) {
	r.m.Lock()
	n, err = r.r.Read(p)
	r.m.Unlock()
	return
}

func (r *bufioReader) ReadBytes(delim byte) ([]byte, error) {
	r.m.Lock()
	bytes, err := r.r.ReadBytes(delim)
	r.m.Unlock()
	return bytes, err
}
