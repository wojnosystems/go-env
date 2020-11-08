package env_parser

import "fmt"

type ParseError struct {
	Path        StructEnvPath
	originalErr error
}

func wrapWithParseError(err error, structPath string, envPath string) *ParseError {
	return &ParseError{
		Path: StructEnvPath{
			StructPath: structPath,
			EnvPath:    envPath,
		},
		originalErr: err,
	}
}

func (p *ParseError) Error() string {
	return fmt.Sprintf("environment variable '%s' failed to parse '%s'", p.Path.EnvPath, p.originalErr.Error())
}
