package readers

import (
	"io/ioutil"
	"strings"

	"github.com/rs/zerolog/log"
)

type DefaultReader struct {
	hooks []Hook
}

func NewDefaultReader() *DefaultReader {
	return &DefaultReader{
		hooks: []Hook{},
	}
}

func (r *DefaultReader) AppendHook(h Hook) *DefaultReader {
	r.hooks = append(r.hooks, h)
	return r
}

func (r *DefaultReader) Read(path string) (map[string]int, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Error().Err(err).Str("file", path).Msg("couldn't open file for reading")
		return nil, err
	}
	tokens := strings.Split(string(content), "\n")
	m := map[string]int{}
	for _, v := range tokens {
		if _, ok := m[v]; !ok {
			for _, h := range r.hooks {
				h.Run(v)
			}
		}
		m[v]++
	}
	return m, nil
}
