package env_parser

// envReader reads environment variables
type envReader interface {
	// Get the value of a single environment with the name envNamed
	Get(envNamed string) string
	// Keys get a list of keys that begin with the prefix. If "" is passed, matches all and returns all keys
	Keys(prefix string) []string
}

type SetReceiver interface {
	// Receive the notice that a value was parsed and set at the fullPath in the destination structure
	// This will allow the flick library to know which values were updated from which source.
	ReceiveSet(fullPath string, value string)
}
