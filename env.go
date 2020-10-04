package env_parser

import (
	"github.com/wojnosystems/go-parse-register"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// Env creates an environment parser given the provided registry
type Env struct {
	// envReader is the source of environment variables.
	// If you leave it blank, it will default to using the operating system environment variables with no prefixes.
	envReader envReader
	// ParseRegistry maps go-default and custom types to members of the provided structure. If left blank, defaults to just Go's primitives being mapped
	ParseRegistry *parse_register.Registry
}

// Unmarshall reads the environment variables and writes them to into.
// into should be a reference to a struct
// This method will do some basic checks on the into value, but to help developers pass in the correct values
func (e *Env) UnmarshallWithEmitter(into interface{}, emitter SetReceiver) (err error) {
	rootV := reflect.ValueOf(into)
	err = e.validateDestination(rootV, rootV.Type())
	if err != nil {
		return
	}
	err = e.unmarshall("", rootV.Elem(), rootV.Elem().Type(), emitter)
	return
}

// Unmarshall reads the environment variables and writes them to into.
// into should be a reference to a struct
// This method will do some basic checks on the into value, but to help developers pass in the correct values
func (e *Env) Unmarshall(into interface{}) (err error) {
	return e.UnmarshallWithEmitter(into, nil)
}

// validateDestination does some basic checks to help users of this class avoid common pitfalls with more helpful messages
func (e *Env) validateDestination(rootV reflect.Value, rootT reflect.Type) (err error) {
	if rootV.IsNil() {
		return NewErrProgramming("'into' argument must be not be nil")
	}
	if rootT.Kind() != reflect.Ptr {
		return NewErrProgramming("'into' argument must be a reference")
	}
	if rootV.Elem().Kind() != reflect.Struct {
		return NewErrProgramming("'into' argument must be a struct")
	}
	return nil
}

// unmarshall is the internal method, which can be called recursively. This performs the heavy-lifting
func (e *Env) unmarshall(structParentPath string, structRefV reflect.Value, structRefT reflect.Type, emitter SetReceiver) (err error) {
	for i := 0; i < structRefV.NumField(); i++ {
		fieldV := structRefV.Field(i)
		if fieldV.CanSet() {
			fieldT := structRefT.Field(i)
			fieldName := fieldNameOrDefault(fieldT)
			structFullPath := appendStructPath(structParentPath, fieldName)

			if fieldT.Type.Kind() == reflect.Slice {
				err = e.unmarshallSlice(structFullPath, fieldV, emitter)
				if err != nil {
					return
				}
			} else {
				envValue := e.envs().Get(structToEnvPath(structFullPath))
				if "" != envValue {
					var wasCalled bool
					wasCalled, err = e.parseRegistry().SetValue(fieldV.Addr().Interface(), envValue)
					if err != nil {
						err = &ParseError{
							Path: StructEnvPath{
								StructPath: structFullPath,
								EnvPath:    structToEnvPath(structFullPath),
							},
							originalErr: err,
						}
						return
					}
					if wasCalled {
						if emitter != nil {
							emitter.ReceiveSet(structFullPath, envValue)
						}
					} else {
						// fall back
						if fieldT.Type.Kind() == reflect.Struct {
							err = e.unmarshall(structFullPath, fieldV, fieldV.Type(), emitter)
							if err != nil {
								return
							}
						}
					}
				} else {
					// fall back
					if fieldT.Type.Kind() == reflect.Struct {
						err = e.unmarshall(structFullPath, fieldV, fieldT.Type, emitter)
						if err != nil {
							return
						}
					}
				}
			}
		}
	}
	return nil
}

var defaultEnvReader = &osEnv{}

// envs obtains the current envReader, or uses the osEnvs by default
func (e *Env) envs() envReader {
	if e.envReader == nil {
		e.envReader = defaultEnvReader
	}
	return e.envReader
}

var defaultParseRegister = parse_register.RegisterGoPrimitives(&parse_register.Registry{})

// parseRegistry obtains a copy of the current registry, or uses the default go primitives, for convenience
func (e *Env) parseRegistry() *parse_register.Registry {
	if e.ParseRegistry == nil {
		e.ParseRegistry = defaultParseRegister
	}
	return e.ParseRegistry
}

// fieldNameOrDefault attempts to read the tags to obtain an alternate name, if no tag found, defaults back to
// using the name provided to the field when the member was defined in Go
func fieldNameOrDefault(fieldT reflect.StructField) (fieldName string) {
	fieldName = fieldT.Tag.Get("env")
	if "" == fieldName {
		fieldName = fieldT.Name
	}
	return
}

// appendStructPath concatenates the parent path name with the current field's name
func appendStructPath(parent string, name string) string {
	if parent != "" {
		return parent + "." + name
	}
	return name
}

// unmarshallSlice operates on a slice of objects. It will initialize the slice, then populate all of its members
// from the environment variables
func (e *Env) unmarshallSlice(sliceFieldPath string, sliceValue reflect.Value, emitter SetReceiver) (err error) {
	var length int
	length, err = elementsInSliceWithAddressPrefix(e.envs(), structToEnvPath(sliceFieldPath)+"_")
	if err != nil {
		return
	}
	if length > 0 {
		newSlice := reflect.MakeSlice(sliceValue.Type(), length, length)
		sliceValue.Set(newSlice)
		for i := 0; i < length; i++ {
			sliceElement := newSlice.Index(i)
			err = e.unmarshall(sliceFieldPath+"["+strconv.FormatInt(int64(i), 10)+"]", sliceElement, sliceElement.Type(), emitter)
			if err != nil {
				return
			}
		}
	}
	return
}

// elementsInSliceWithAddressPrefix returns how big a slice should be to hold all of the variables defined in the environment
func elementsInSliceWithAddressPrefix(env envReader, pathPrefix string) (length int, err error) {
	maxIndex := int64(-1)
	for _, key := range env.Keys(pathPrefix) {
		possibleNumber := envIndexRegexp.FindString(key[len(pathPrefix):])
		if "" != possibleNumber {
			var index int64
			index, err = strconv.ParseInt(possibleNumber, 10, 0)
			if err != nil {
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

var envIndexRegexp = regexp.MustCompile(`^(\d+)`)

const replaceWith = "_"

func structToEnvPath(structPath string) (envPath string) {
	envPath = strings.ReplaceAll(structPath, "].", replaceWith)
	envPath = strings.ReplaceAll(envPath, ".", replaceWith)
	envPath = strings.ReplaceAll(envPath, "]", replaceWith)
	envPath = strings.ReplaceAll(envPath, "[", replaceWith)
	return
}
