package flag_unmarshaler

import (
	"github.com/stretchr/testify/assert"
	"github.com/wojnosystems/go-optional"
	"testing"
	"time"
)

func TestType_Unmarshall(t *testing.T) {
	cases := map[string]struct {
		input    Group
		expected appConfigMock
	}{
		"empty": {
			input: Group{
				"",
				[]KeyValue{},
			},
			expected: appConfigMock{},
		},
		"name": {
			input: Group{
				"do-thing",
				[]KeyValue{
					{
						Key:   "--name",
						Value: "chris",
					},
				},
			},
			expected: appConfigMock{
				Name: optional.StringFrom("chris"),
			},
		},
		"thread count": {
			input: Group{
				"do-thing",
				[]KeyValue{
					{
						Key:   "-c",
						Value: "5",
					},
				},
			},
			expected: appConfigMock{
				ThreadCount: optional.IntFrom(5),
			},
		},
		"slice": {
			input: Group{
				"do-thing",
				[]KeyValue{
					{
						Key:   "--databases[1].h",
						Value: "example.org",
					},
					{
						Key:   "--databases[0].host",
						Value: "example.com",
					},
				},
			},
			expected: appConfigMock{
				Databases: []dbConfigMock{
					{
						Host: optional.StringFrom("example.com"),
					},
					{
						Host: optional.StringFrom("example.org"),
					},
				},
			},
		},
		"nested": {
			input: Group{
				"do-thing",
				[]KeyValue{
					{
						Key:   "--databases[0].Nested.ConnTimeout",
						Value: "30s",
					},
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
		"bool": {
			input: Group{
				"do-thing",
				[]KeyValue{
					{
						Key:   "--Enabled",
						Value: "true",
					},
				},
			},
			expected: appConfigMock{
				Enabled: true,
			},
		},
	}
	for caseName, c := range cases {
		t.Run(caseName, func(t *testing.T) {
			var actual appConfigMock
			underTest := New(&c.input)
			err := underTest.Unmarshal(&actual)
			assert.NoError(t, err)
			assert.True(t, c.expected.IsEqual(&actual))
		})
	}
}
