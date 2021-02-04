package flag_unmarshaler

// Defines a set of objects used with testing

import (
	"github.com/wojnosystems/go-optional/v2"
	"time"
)

func stringEqual(a, b optional.String) bool {
	equal := false
	a.IfSetElse(func(aVal string) {
		b.IfSet(func(bVal string) {
			equal = aVal == bVal
		})
	}, func() {
		equal = !b.IsSet()
	})
	return equal
}

func intEqual(a, b optional.Int) bool {
	equal := false
	a.IfSetElse(func(aVal int) {
		b.IfSet(func(bVal int) {
			equal = aVal == bVal
		})
	}, func() {
		equal = !b.IsSet()
	})
	return equal
}

func durationEqual(a, b optional.Duration) bool {
	equal := false
	a.IfSetElse(func(aVal time.Duration) {
		b.IfSet(func(bVal time.Duration) {
			equal = aVal == bVal
		})
	}, func() {
		equal = !b.IsSet()
	})
	return equal
}

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
	if !stringEqual(m.Name, o.Name) || !intEqual(m.ThreadCount, o.ThreadCount) || m.Enabled != o.Enabled {
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
	return stringEqual(m.Host, o.Host) &&
		stringEqual(m.User, o.User) &&
		stringEqual(m.Password, o.Password) &&
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
	return durationEqual(m.ConnTimeout, o.ConnTimeout)
}
