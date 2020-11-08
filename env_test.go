package env_parser

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/wojnosystems/go-optional"
	"testing"
	"time"
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
		"db nested": {
			env: &envMock{
				mock: map[string]string{
					"Databases_0_Nested_ConnTimeout": "30s",
				},
			},
			expected: appConfigMock{
				Databases: []dbConfigMock{
					{
						Nested: nestedDbConfigMock{ConnTimeout: optional.DurationFrom(30 * time.Second)},
					},
				},
			},
		},
	}

	for caseName, c := range cases {
		t.Run(caseName, func(t *testing.T) {
			actual := &appConfigMock{}
			e := NewWithEnvReader(c.env)
			err := e.Unmarshall(actual)
			assert.NoError(t, err)
			assert.True(t, c.expected.IsEqual(actual))
		})
	}
}

func TestEnv_UnmarshallErrors(t *testing.T) {
	cases := map[string]struct {
		env      *envMock
		expected error
	}{
		"field does not exist": {
			env: &envMock{
				mock: map[string]string{
					"IDontExist": "nothing to set",
				},
			},
		},
		"time field does not parse": {
			env: &envMock{
				mock: map[string]string{
					"Databases_0_Nested_ConnTimeout": "P30s",
				},
			},
			expected: &ParseError{
				Path: StructEnvPath{
					StructPath: "Databases[0].Nested.ConnTimeout",
					EnvPath:    "Databases_0_Nested_ConnTimeout",
				},
				originalErr: fmt.Errorf("time: invalid duration P30s"),
			},
		},
		"int field does not parse": {
			env: &envMock{
				mock: map[string]string{
					"ThreadCount": "P",
				},
			},
			expected: &ParseError{
				Path: StructEnvPath{
					StructPath: "ThreadCount",
					EnvPath:    "ThreadCount",
				},
				originalErr: fmt.Errorf(`strconv.ParseInt: parsing "P": invalid syntax`),
			},
		},
		"slice length field does not parse": {
			env: &envMock{
				mock: map[string]string{
					"Databases_9999999999999999999999999999999999999999999999999999999999999999_.Host": "P",
				},
			},
			expected: &ParseError{
				Path: StructEnvPath{
					StructPath: "Databases",
					EnvPath:    "Databases",
				},
				originalErr: fmt.Errorf(`strconv.ParseInt: parsing "9999999999999999999999999999999999999999999999999999999999999999": value out of range`),
			},
		},
	}
	for caseName, c := range cases {
		t.Run(caseName, func(t *testing.T) {
			actual := &appConfigMock{}
			e := NewWithEnvReader(c.env)
			err := e.Unmarshall(actual)
			if c.expected == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, c.expected.Error())
			}
		})
	}
}
