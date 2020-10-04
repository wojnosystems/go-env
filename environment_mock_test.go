package env_parser

import "strings"

// Implements the default system environment getters
// this allows us to test this
type envMock struct {
	mock map[string]string
}

func (s *envMock) Get(envNamed string) (out string) {
	if s.mock == nil {
		return
	}
	var ok bool
	out, ok = s.mock[envNamed]
	if !ok {
		return
	}
	return
}

func (s *envMock) Keys(prefix string) (out []string) {
	for key := range s.mock {
		if strings.HasPrefix(key, prefix) {
			out = append(out, key)
		}
	}
	return
}
