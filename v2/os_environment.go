package v2

import (
	"os"
	"strings"
)

// OsEnv Implements the default system environment getters
type OsEnv struct {
}

func (s *OsEnv) Get(envNamed string) string {
	return os.Getenv(envNamed)
}

func (s *OsEnv) Keys(prefix string) (out []string) {
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
