package env_parser

import (
	"github.com/wojnosystems/go-parse-register"
	"reflect"
	"regexp"
	"strconv"
)

// Env creates an environment parser given the provided registry
type Env struct {
	envReader envReader
	// ParseRegistry maps go-default and custom types to members of the provided structure
	ParseRegistry *parse_register.Registry
}

// Unmarshall reads the environment variables and writes them to into.
// into should be a reference to a struct
// This method will do some basic checks on the into value, but to help developers pass in the correct values
func (e *Env) Unmarshall(into interface{}) (err error) {
	rootV := reflect.ValueOf(into)
	err = e.validateDestination(rootV, rootV.Type())
	if err != nil {
		return
	}
	err = e.unmarshall("", rootV.Elem(), rootV.Elem().Type())
	if err != nil {
		return
	}
	return nil
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
func (e *Env) unmarshall(parentName string, structRefV reflect.Value, structRefT reflect.Type) (err error) {
	for i := 0; i < structRefV.NumField(); i++ {
		fieldV := structRefV.Field(i)
		if fieldV.CanSet() {
			fieldT := structRefT.Field(i)
			fieldName := fieldNameOrDefault(fieldT)
			fullPath := nameFromParent(fieldName, parentName)

			if fieldT.Type.Kind() == reflect.Slice {
				err = e.unmarshallSlice(fullPath, fieldV)
				if err != nil {
					return
				}
			} else {
				envValue := e.envs().Get(fullPath)
				if "" != envValue {
					var wasCalled bool
					wasCalled, err = e.parseRegistry().SetValue(fieldV.Addr().Interface(), envValue)
					if err != nil {
						return
					}
					if !wasCalled {
						// fall back
						if fieldT.Type.Kind() == reflect.Struct {
							err = e.unmarshall(fullPath, fieldV, fieldV.Type())
							if err != nil {
								return
							}
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

// nameFromParent concatenates the parent path name with the current field's name
func nameFromParent(name string, parent string) string {
	if parent != "" {
		return parent + "." + name
	}
	return name
}

// unmarshallSlice operates on a slice of objects. It will initialize the slice, then populate all of its members
// from the environment variables
func (e *Env) unmarshallSlice(path string, sliceValue reflect.Value) (err error) {
	var length int
	length, err = elementsInSliceWithAddressPrefix(e.envs(), path+"_")
	if err != nil {
		return
	}
	if length > 0 {
		newSlice := reflect.MakeSlice(sliceValue.Type(), length, length)
		sliceValue.Set(newSlice)
		for i := 0; i < length; i++ {
			sliceElement := newSlice.Index(i)
			err = e.unmarshall(path+"_"+strconv.FormatInt(int64(i), 10)+"_", sliceElement, sliceElement.Type())
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
