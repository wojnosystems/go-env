package env_parser

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSelectKeysWithPrefix(t *testing.T) {
	cases := map[string]struct {
		keys     []string
		prefix   string
		expected []string
	}{
		"empty": {
			expected: []string{},
		},
		"matches subset": {
			keys:     []string{"abcde", "abcdf", "abcdg", "bcd", "azy"},
			prefix:   "abcd",
			expected: []string{"abcde", "abcdf", "abcdg"},
		},
		"matches none": {
			keys:     []string{"abcde", "abcdf", "abcdg", "bcd", "azy"},
			prefix:   "x",
			expected: []string{},
		},
	}

	for caseName, c := range cases {
		t.Run(caseName, func(t *testing.T) {
			actual := SelectKeysWithPrefix(c.keys, c.prefix)
			assert.Equal(t, c.expected, actual)
		})
	}
}
