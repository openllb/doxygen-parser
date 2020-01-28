package doxygen

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type testCase struct {
	name     string
	input    string
	expected *Group
}

func cleanup(value string) string {
	return strings.ReplaceAll(value, strings.Repeat("\t", 3), "")
}

func TestParse(t *testing.T) {
	for _, tc := range []testCase{
		{
			"empty",
			``,
			&Group{},
		},
		{
			"no commands",
			`Hello world
			`,
			&Group{
				Doc: "Hello world\n",
			},
		},
		{
			"returns",
			`Hello world
			@return an object
			`,
			&Group{
				Doc: "Hello world\n",
				Return: Return{
					Description: "an object",
				},
			},
		},
		{
			"docs after commands",
			`Hello world
			@return an object

			More docs
			`,
			&Group{
				Doc: "Hello world\n\nMore docs\n",
				Return: Return{
					Description: "an object",
				},
			},
		},
		{
			"params",
			`Hello world
			@param arg1 an param
			@return an object
			`,
			&Group{
				Doc: "Hello world\n",
				Params: []Param{
					{
						Name: "arg1",
						Description: "an param",
					},
				},
				Return: Return{
					Description: "an object",
				},
			},
		},
		{
			"multi-line params",
			`Hello world
			@param arg1 multi
			line params
			@return an object
			`,
			&Group{
				Doc: "Hello world\n",
				Params: []Param{
					{
						Name: "arg1",
						Description: "multi line params",
					},
				},
				Return: Return{
					Description: "an object",
				},
			},
		},
		{
			"multiple params",
			`Hello world
			@param arg1 multi
			line params
			@param arg2 another param
			@return an object
			`,
			&Group{
				Doc: "Hello world\n",
				Params: []Param{
					{
						Name: "arg1",
						Description: "multi line params",
					},
					{
						Name: "arg2",
						Description: "another param",
					},
				},
				Return: Return{
					Description: "an object",
				},
			},
		},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			group, err := Parse(strings.NewReader(cleanup(tc.input)))
			require.NoError(t, err)
			require.Equal(t, *tc.expected, *group)
		})
	}
}
