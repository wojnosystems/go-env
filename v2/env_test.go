package v2

import (
	"github.com/stretchr/testify/assert"
	"github.com/wojnosystems/go-optional/v2"
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
					"Databases_0_NEST_TIMEOUT": "30s",
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
		expected string
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
					"Databases_0_NEST_TIMEOUT": "P30s",
				},
			},
			expected: `environment variable 'Databases_0_NEST_TIMEOUT' failed to parse because time: invalid duration P30s`,
		},
		"int field does not parse": {
			env: &envMock{
				mock: map[string]string{
					"ThreadCount": "P",
				},
			},
			expected: `environment variable 'ThreadCount' failed to parse because strconv.ParseInt: parsing "P": invalid syntax`,
		},
		"slice length field does not parse": {
			env: &envMock{
				mock: map[string]string{
					"Databases_9999999999999999999999999999999999999999999999999999999999999999_.Host": "P",
				},
			},
			expected: `environment variable 'Databases' failed to parse because strconv.ParseInt: parsing "9999999999999999999999999999999999999999999999999999999999999999": value out of range`,
		},
	}
	for caseName, c := range cases {
		t.Run(caseName, func(t *testing.T) {
			actual := &appConfigMock{}
			e := NewWithEnvReader(c.env)
			err := e.Unmarshall(actual)
			if c.expected == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, c.expected)
			}
		})
	}
}
