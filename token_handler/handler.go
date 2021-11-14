package tokenhandler

import (
	"encoding/json"
	"io/ioutil"

	"github.com/mkudimov/palm/token_handler/readers"
)

type Handler struct {
	reader readers.Reader
	dump   func(m map[string]int, out string) error
}

func NewHandler(reader readers.Reader) *Handler {
	return &Handler{
		reader: reader,
		dump:   dump,
	}
}
func (h *Handler) Handle(in, out string) error {
	m, err := h.reader.Read(in)
	if err != nil {
		return err
	}
	return h.dump(m, out)
}

func dump(m map[string]int, out string) error {
	mm := map[string]int{}
	for k, v := range m {
		if v > 1 {
			mm[k] = v
		}
	}
	bytes, err := json.Marshal(mm)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(out, bytes, 0644)
}
