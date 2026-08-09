package dsign_test

import (
	"dsign"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

var tests = []struct {
	name   string
	input  string
	output string
}{
	{
		name:   "does not change the URL unnecessarily",
		input:  "http://example.dynata.com/./hello/../world//test?a=%3F&b=%2A&c=%2B",
		output: "http://example.dynata.com/./hello/../world//test?a=%3F&b=%2A&c=%2B",
	},
	{
		name:   "adds missing scheme as https",
		input:  "//example.dynata.com/hello?world=test",
		output: "https://example.dynata.com/hello?world=test",
	},
	{
		name:   "orders query string parameters",
		input:  "http://example.dynata.com/?z=1&y=2&x=3",
		output: "http://example.dynata.com/?x=3&y=2&z=1",
	},
	{
		name:   "adds path prefix slash",
		input:  "http://example.dynata.com?hello=world",
		output: "http://example.dynata.com/?hello=world",
	},
	{
		name:   "drops port number",
		input:  "http://example.dynata.com:8080/?hello=world",
		output: "http://example.dynata.com/?hello=world",
	},
	{
		name:   "keeps only last duplicate query sequence",
		input:  "http://example.dynata.com/?hello=world&hello=dlrow&hello=empty",
		output: "http://example.dynata.com/?hello=empty",
	},
}

func TestNormalizeURL(t *testing.T) {
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			u, err := url.Parse(test.input)
			require.NoError(t, err)
			output := dsign.NormalizeURL(u)
			require.Equal(t, test.output, output.String())
		})
	}
}
