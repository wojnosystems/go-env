package v2

import (
	into_struct "github.com/wojnosystems/go-into-struct"
	"github.com/wojnosystems/go-optional-parse-registry/v2"
	"github.com/wojnosystems/go-parse-register"
	"regexp"
)

// Env creates an environment parser given the provided registry
type Env struct {
	config envInternal
}

func New() *Env {
	return NewWithParseRegistryEmitterEnvReader(defaultParseRegister, defaultNoOpSetReceiver, defaultEnvReader)
}

func NewWithParseRegistryEmitterEnvReader(parseRegistry parse_register.ValueSetter, emitter SetReceiver, reader envReader) *Env {
	return &Env{
		config: envInternal{
			envReader:     reader,
			parseRegistry: parseRegistry,
			emitter:       emitter,
		},
	}
}

func NewWithParseRegistry(registry parse_register.ValueSetter) *Env {
	return NewWithParseRegistryEmitterEnvReader(registry, defaultNoOpSetReceiver, defaultEnvReader)
}

func NewWithEnvReader(reader envReader) *Env {
	return NewWithParseRegistryEmitterEnvReader(defaultParseRegister, defaultNoOpSetReceiver, reader)
}

func NewWithEmitter(emitter SetReceiver) *Env {
	return NewWithParseRegistryEmitterEnvReader(defaultParseRegister, emitter, defaultEnvReader)
}

func NewWithParseRegistryWithEmitter(registry parse_register.ValueSetter, emitter SetReceiver) *Env {
	return NewWithParseRegistryEmitterEnvReader(registry, emitter, defaultEnvReader)
}

// Unmarshall reads the environment variables and writes them to into.
// into should be a reference to a struct
// This method will do some basic checks on the into value, but to help developers pass in the correct values
func (e *Env) Unmarshall(into interface{}) (err error) {
	return into_struct.Unmarshall(into, &e.config)
}

var (
	defaultEnvReader       = &OsEnv{}
	defaultNoOpSetReceiver = &SetReceiverNoOp{}
	defaultParseRegister   = optional_parse_registry.NewWithGoPrimitives()
	envIndexRegexp         = regexp.MustCompile(`^(\d+)`)
)
