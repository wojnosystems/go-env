package env_parser

import (
	"fmt"
	into_struct "github.com/wojnosystems/go-into-struct"
	parse_register "github.com/wojnosystems/go-parse-register"
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
func (e *envInternal) SetValue(structFullPath into_struct.Path) (handled bool, err error) {
	field := structFullPath.Top()
	if field == nil {
		return
	}
	envPath := structToEnvPath(structFullPath)
	envValue := e.envReader.Get(envPath)
	if "" != envValue {
		// Some environment value was set, use it
		valueDst := field.Value().Addr().Interface()
		handled, err = e.parseRegistry.SetValue(valueDst, envValue)
		if err != nil {
			err = newParseError(structFullPath.String(), envPath, err)
			return
		}
		if handled {
			e.emitter.ReceiveSet(structFullPath, envPath, envValue)
			return
		}
	}
	return
}

func (e *envInternal) SliceLen(structFullPath into_struct.Path) (length int, err error) {
	envPath := structToEnvPath(structFullPath)
	pathPrefix := envPath + "_"
	maxIndex := int64(-1)
	for _, key := range e.envReader.Keys(pathPrefix) {
		possibleNumber := envIndexRegexp.FindString(key[len(pathPrefix):])
		if "" != possibleNumber {
			var index int64
			index, err = strconv.ParseInt(possibleNumber, 10, 0)
			if err != nil {
				err = newParseError(structFullPath.String(), envPath, err)
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

const envFieldSeparator = "_"

func structToEnvPath(structPath into_struct.Path) string {
	envParts := make([]string, 0, len(structPath.Parts()))
	for _, pathPart := range structPath.Parts() {
		fieldEnvName := pathPart.StructField().Tag.Get("env")
		if fieldEnvName == "" {
			fieldEnvName = pathPart.StructField().Name
		}
		switch t := pathPart.(type) {
		case into_struct.PathSliceParter:
			envParts = append(envParts, fmt.Sprintf("%s%s%d%s", fieldEnvName, envFieldSeparator, t.Index(), envFieldSeparator))
		default:
			envParts = append(envParts, fieldEnvName)
		}
	}
	sb := strings.Builder{}
	for i, part := range envParts {
		sb.WriteString(part)
		if i < len(envParts)-1 && !strings.HasSuffix(part, envFieldSeparator) {
			sb.WriteString(envFieldSeparator)
		}
	}
	return sb.String()
}
