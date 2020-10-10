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
	// structPath where the value was set in the structure in go. base.value[index].othervalue = value
	// envName is the environment variable used to look up the value
	// value is what was read from the environment for the envName key
	ReceiveSet(structPath string, envName string, value string)
}
