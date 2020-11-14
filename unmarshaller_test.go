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
		"nested named wrong name": {
			input: Group{
				"do-thing",
				[]KeyValue{
					{
						Key:   "--databases[0].nestedwrong.ConnTimeout",
						Value: "30s",
					},
				},
			},
			expected: appConfigMock{
				Databases: []dbConfigMock{
					{},
				},
			},
		},
		"nested named correct name": {
			input: Group{
				"do-thing",
				[]KeyValue{
					{
						Key:   "--databases[0].nestednamed.ConnTimeout",
						Value: "10s",
					},
					{
						Key:   "--databases[0].h",
						Value: "test.example.com",
					},
					{
						Key:   "--databases[1].n.ConnTimeout",
						Value: "20s",
					},
					{
						Key:   "--databases[2].n.c",
						Value: "30s",
					},
					{
						Key:   "-d[3].n.c",
						Value: "40s",
					},
					{
						Key:   "-d[4].nestednamed.c",
						Value: "50s",
					},
				},
			},
			expected: appConfigMock{
				Databases: []dbConfigMock{
					{
						NestedNamed: nestedDbConfigMock{
							ConnTimeout: optional.DurationFrom(10 * time.Second),
						},
						Host: optional.StringFrom("test.example.com"),
					},
					{
						NestedNamed: nestedDbConfigMock{
							ConnTimeout: optional.DurationFrom(20 * time.Second),
						},
					},
					{
						NestedNamed: nestedDbConfigMock{
							ConnTimeout: optional.DurationFrom(30 * time.Second),
						},
					},
					{
						NestedNamed: nestedDbConfigMock{
							ConnTimeout: optional.DurationFrom(40 * time.Second),
						},
					},
					{
						NestedNamed: nestedDbConfigMock{
							ConnTimeout: optional.DurationFrom(50 * time.Second),
						},
					},
				},
			},
		},
		"multi-bool": {
			input: Group{
				"do-thing",
				[]KeyValue{
					{
						Key: "-efh",
					},
				},
			},
			expected: appConfigMock{
				BoolE: true,
				BoolF: true,
				BoolG: false,
				BoolH: true,
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
