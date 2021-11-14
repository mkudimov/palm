package readers

import "sync"

type cmap struct {
	s map[string]int
	m sync.RWMutex
}

func newCmap() *cmap {
	return &cmap{
		s: map[string]int{},
	}
}

func (m *cmap) increase(k string) {
	m.m.Lock()
	defer m.m.Unlock()
	m.s[k]++
}

func (m *cmap) get(k string) (int, bool) {
	m.m.RLock()
	defer m.m.RUnlock()
	v, ok := m.s[k]
	return v, ok
}
