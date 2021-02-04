package v2

// Defines a set of objects used with testing

import (
	"github.com/wojnosystems/go-optional/v2"
)

type appConfigMock struct {
	Name        optional.String
	ThreadCount optional.Int
	Databases   []dbConfigMock
}

func (m appConfigMock) IsEqual(o *appConfigMock) bool {
	if o == nil {
		return false
	}
	if !m.Name.IsEqual(o.Name) || !m.ThreadCount.IsEqual(o.ThreadCount) {
		return false
	}
	if len(m.Databases) != len(o.Databases) {
		return false
	}
	for i, database := range m.Databases {
		if !database.IsEqual(&o.Databases[i]) {
			return false
		}
	}
	return true
}

type dbConfigMock struct {
	Host     optional.String
	User     optional.String
	Password optional.String
	Nested   nestedDbConfigMock `env:"NEST"`
}

func (m dbConfigMock) IsEqual(o *dbConfigMock) bool {
	if o == nil {
		return false
	}
	return m.Host.IsEqual(o.Host) && m.User.IsEqual(o.User) && m.Password.IsEqual(o.Password) && m.Nested.IsEqual(&o.Nested)
}

type nestedDbConfigMock struct {
	ConnTimeout optional.Duration `env:"TIMEOUT"`
}

func (m nestedDbConfigMock) IsEqual(o *nestedDbConfigMock) bool {
	if o == nil {
		return false
	}
	return m.ConnTimeout.IsEqual(o.ConnTimeout)
}
