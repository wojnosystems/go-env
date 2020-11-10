package env_parser

import (
	parse_register "github.com/wojnosystems/go-parse-register"
	"reflect"
	"strconv"
	"strings"
)

// envInternal hides the methods that implement the intoStruct parser
type envInternal struct {
	// envReader is the source of environment variables.
	// If you leave it blank, it will default to using the operating system environment variables with no prefixes.
	envReader envReader
	// ParseRegistry maps go-default and custom types to members of the provided structure. If left blank, defaults to just Go's primitives being mapped
	parseRegistry parse_register.ValueSetter
	emitter       SetReceiver
}

// SetValue
func (e *envInternal) SetValue(structFullPath string, fieldV reflect.Value) (handled bool, err error) {
	envPath := structToEnvPath(structFullPath)
	envValue := e.envReader.Get(envPath)
	if "" != envValue {
		// Some environment value was set, use it
		valueDst := fieldV.Addr().Interface()
		handled, err = e.parseRegistry.SetValue(valueDst, envValue)
		if err != nil {
			err = wrapWithParseError(err, structFullPath, envPath)
			return
		}
		if handled {
			e.emitter.ReceiveSet(structFullPath, envPath, envValue)
			return
		}
	}
	return
}

func (e *envInternal) SliceLen(structFullPath string) (length int, err error) {
	envPath := structToEnvPath(structFullPath)
	pathPrefix := envPath + "_"
	maxIndex := int64(-1)
	for _, key := range e.envReader.Keys(pathPrefix) {
		possibleNumber := envIndexRegexp.FindString(key[len(pathPrefix):])
		if "" != possibleNumber {
			var index int64
			index, err = strconv.ParseInt(possibleNumber, 10, 0)
			if err != nil {
				err = wrapWithParseError(err, structFullPath, envPath)
				return
			}
			if index > maxIndex {
				maxIndex = index
			}
		}
	}
	length = int(maxIndex + 1)
	return
}

const replaceWith = "_"

func structToEnvPath(structPath string) (envPath string) {
	envPath = strings.ReplaceAll(structPath, "].", replaceWith)
	envPath = strings.ReplaceAll(envPath, ".", replaceWith)
	envPath = strings.ReplaceAll(envPath, "]", replaceWith)
	envPath = strings.ReplaceAll(envPath, "[", replaceWith)
	return
}
