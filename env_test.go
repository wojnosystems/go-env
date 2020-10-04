package env_parser

import (
	"github.com/stretchr/testify/assert"
	"github.com/wojnosystems/go-optional"
	"github.com/wojnosystems/go-optional-parse-registry"
	"github.com/wojnosystems/go-parse-register"
	"testing"
)

func TestEnv_Unmarshall(t *testing.T) {
	cases := map[string]struct {
		env      *envMock
		expected appConfigMock
	}{
		"nothing": {
			env: &envMock{},
		},
		"name": {
			env: &envMock{
				mock: map[string]string{
					"Name": "SuperServer",
				},
			},
			expected: appConfigMock{
				Name: optional.StringFrom("SuperServer"),
			},
		},
		"db[1].Host": {
			env: &envMock{
				mock: map[string]string{
					"Databases_1_Host": "example.com",
				},
			},
			expected: appConfigMock{
				Databases: []dbConfigMock{
					{},
					{
						Host: optional.StringFrom("example.com"),
					},
				},
			},
		},
	}

	for caseName, c := range cases {
		t.Run(caseName, func(t *testing.T) {
			actual := &appConfigMock{}
			e := Env{
				envReader:     c.env,
				ParseRegistry: optional_parse_registry.Register(parse_register.RegisterGoPrimitives(&parse_register.Registry{})),
			}
			err := e.Unmarshall(actual)
			assert.NoError(t, err)
			assert.True(t, c.expected.IsEqual(actual))
		})
	}
}
