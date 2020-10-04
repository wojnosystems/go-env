package env_parser

import (
	"os"
	"strings"
)

// Implements the default system environment getters
// this allows us to test this
type osEnv struct {
}

func (s *osEnv) Get(envNamed string) string {
	return os.Getenv(envNamed)
}

func (s *osEnv) Keys(prefix string) (out []string) {
	for _, key := range os.Environ() {
		if strings.HasPrefix(key, prefix) {
			out = append(out, key)
		}
	}
	return
}
