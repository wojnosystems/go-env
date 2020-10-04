package env_parser

import "fmt"

type ParseError struct {
	Path        StructEnvPath
	originalErr error
}

func (p *ParseError) Error() string {
	return fmt.Sprintf("environment variable '%s' failed to parse '%s'", p.Path.EnvPath, p.originalErr.Error())
}
