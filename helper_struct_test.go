package flag_unmarshaler

// Defines a set of objects used with testing

import (
	"github.com/wojnosystems/go-optional"
)

type appConfigMock struct {
	Name        optional.String `flag:"name" flag-short:"n"`
	ThreadCount optional.Int    `flag:"thread-count" flag-short:"c"`
	Databases   []dbConfigMock  `flag:"databases"`
	Enabled     bool
}

func (m appConfigMock) IsEqual(o *appConfigMock) bool {
	if o == nil {
		return false
	}
	if !m.Name.IsEqual(o.Name) || !m.ThreadCount.IsEqual(o.ThreadCount) || m.Enabled != o.Enabled {
		return false
	}
	if len(m.Databases) != len(o.Databases) {
		return false
	}
	for i, database := range m.Databases {
		database.IsEqual(&o.Databases[i])
	}
	return true
}

type dbConfigMock struct {
	Host     optional.String `flag:"host" flag-short:"h"`
	User     optional.String
	Password optional.String
	Nested   nestedDbConfigMock
}

func (m dbConfigMock) IsEqual(o *dbConfigMock) bool {
	if o == nil {
		return false
	}
	return m.Host.IsEqual(o.Host) && m.User.IsEqual(o.User) && m.Password.IsEqual(o.Password) && m.Nested.IsEqual(&o.Nested)
}

type nestedDbConfigMock struct {
	ConnTimeout optional.Duration
}

func (m nestedDbConfigMock) IsEqual(o *nestedDbConfigMock) bool {
	if o == nil {
		return false
	}
	return m.ConnTimeout.IsEqual(o.ConnTimeout)
}
