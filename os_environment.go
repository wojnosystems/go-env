package env_parser

import (
	"os"
	"strings"
)

// osEnv Implements the default system environment getters
type osEnv struct {
}

func (s *osEnv) Get(envNamed string) string {
	return os.Getenv(envNamed)
}

func (s *osEnv) Keys(prefix string) (out []string) {
	return SelectKeysWithPrefix(os.Environ(), prefix)
}

// SelectKeysWithPrefix filters the keys to only include those that contain the prefix
// Offered here as a generic way to filter keys in any implementing interfaces
func SelectKeysWithPrefix(keys []string, prefix string) (out []string) {
	for _, key := range keys {
		if strings.HasPrefix(key, prefix) {
			out = append(out, key)
		}
	}
	if out == nil {
		out = []string{}
	}
	return
}
