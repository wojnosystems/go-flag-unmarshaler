package flag_unmarshaler

// Defines a set of objects used with testing

import (
	"github.com/wojnosystems/go-optional"
)

type appConfigMock struct {
	Name        optional.String `flag:"name" flag-short:"n"`
	ThreadCount optional.Int    `flag:"thread-count" flag-short:"c"`
	Databases   []dbConfigMock  `flag:"databases" flag-short:"d"`
	Enabled     bool
	BoolE       bool `flag-short:"e"`
	BoolF       bool `flag-short:"f"`
	BoolG       bool `flag-short:"g"`
	BoolH       bool `flag-short:"h"`
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
		if !database.IsEqual(&o.Databases[i]) {
			return false
		}
	}
	return true
}

type dbConfigMock struct {
	Host        optional.String `flag:"host" flag-short:"h"`
	User        optional.String
	Password    optional.String
	Nested      nestedDbConfigMock
	NestedNamed nestedDbConfigMock `flag:"nestednamed" flag-short:"n"`
}

func (m dbConfigMock) IsEqual(o *dbConfigMock) bool {
	if o == nil {
		return false
	}
	return m.Host.IsEqual(o.Host) &&
		m.User.IsEqual(o.User) &&
		m.Password.IsEqual(o.Password) &&
		m.Nested.IsEqual(&o.Nested) &&
		m.NestedNamed.IsEqual(&o.NestedNamed)
}

type nestedDbConfigMock struct {
	ConnTimeout optional.Duration `flag:"ConnTimeout" flag-short:"c"`
}

func (m nestedDbConfigMock) IsEqual(o *nestedDbConfigMock) bool {
	if o == nil {
		return false
	}
	return m.ConnTimeout.IsEqual(o.ConnTimeout)
}
