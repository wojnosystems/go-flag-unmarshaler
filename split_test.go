package flag_unmarshaller

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSplit(t *testing.T) {
	cases := map[string]struct {
		args     []string
		expected []Group
	}{
		"none": {
			expected: []Group{
				{},
			},
		},
		"bool long": {
			args: []string{"--one"},
			expected: []Group{
				{
					Flags: []KeyValue{
						{
							Key:   "--one",
							Value: "true",
						},
					},
				},
			},
		},
		"bool short": {
			args: []string{"-o"},
			expected: []Group{
				{
					Flags: []KeyValue{
						{
							Key:   "-o",
							Value: "true",
						},
					},
				},
			},
		},
		"one split": {
			args: []string{"--one=two"},
			expected: []Group{
				{
					Flags: []KeyValue{
						{
							Key:   "--one",
							Value: "two",
						},
					},
				},
			},
		},
		"one command split": {
			args: []string{"first", "--one=two"},
			expected: []Group{
				{},
				{
					CommandName: "first",
					Flags: []KeyValue{
						{
							Key:   "--one",
							Value: "two",
						},
					},
				},
			},
		},
		"two command split": {
			args: []string{"first", "--one=two", "second", "-three=3", "--four=five"},
			expected: []Group{
				{},
				{
					CommandName: "first",
					Flags: []KeyValue{
						{
							Key:   "--one",
							Value: "two",
						},
					},
				},
				{
					CommandName: "second",
					Flags: []KeyValue{
						{
							Key:   "-t",
							Value: "3",
						},
						{
							Key:   "-h",
							Value: "3",
						},
						{
							Key:   "-r",
							Value: "3",
						},
						{
							Key:   "-e",
							Value: "3",
						},
						{
							Key:   "-e",
							Value: "3",
						},
						{
							Key:   "--four",
							Value: "five",
						},
					},
				},
			},
		},
		"stop": {
			args: []string{"first", "--one=two", "second", "--", "--four=five"},
			expected: []Group{
				{},
				{
					CommandName: "first",
					Flags: []KeyValue{
						{
							Key:   "--one",
							Value: "two",
						},
					},
				},
				{
					CommandName: "second",
				},
			},
		},
	}
	for caseName, c := range cases {
		t.Run(caseName, func(t *testing.T) {
			actual := Split(c.args)
			assert.Equal(t, c.expected, actual)
		})
	}
}
