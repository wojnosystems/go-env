package env_parser

import "fmt"

type ParseError struct {
	Path        StructEnvPath
	originalErr error
}

func newParseError(structPath string, envPath string, original error) *ParseError {
	return &ParseError{
		Path: StructEnvPath{
			StructPath: structPath,
			EnvPath:    envPath,
		},
		originalErr: original,
	}
}

func (p *ParseError) Error() string {
	return fmt.Sprintf("environment variable '%s' failed to parse because %s", p.Path.EnvPath, p.originalErr.Error())
}
