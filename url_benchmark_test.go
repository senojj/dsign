package dsign_test

import (
	"dsign"
	"net/url"
	"testing"
)

var benchmarks = []struct {
	name  string
	input string
}{
	{
		name:  "simple",
		input: "https://example.dynata.com/hello?hello=world&hello=dlrow",
	},
	{
		name:  "complex",
		input: "https://example.dynata.com/hello/./world/..?tes%3Bt=tha%2Ft&what=the&after=%3F%3F%2F%2F%2F%2B%3F%3F%20%20%20%20%3B%3B&next=twas&then",
	},
}

func BenchmarkMyCanonicalURL(b *testing.B) {
	for _, test := range benchmarks {
		b.Run(test.name, func(b *testing.B) {
			for b.Loop() {
				u, err := url.Parse(test.input)
				if err != nil {
					b.Fatal(err)
				}
				_ = dsign.EncodeURL(u)
			}
		})
	}
}
