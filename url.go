package dsign

import (
	"dsign/internal/encoding/percent"
	"net/url"
	"slices"
	"strings"
)

func normalizePath(value string) string {
	if len(value) == 0 || value[0] != '/' {
		return "/" + value
	}
	return value
}

func normalizeQuery(query string) string {
	size := strings.Count(query, "&") + 1
	names := make([]string, 0, size)
	sequences := make(map[string]string, size)

	size = 0

	for sequence := range strings.SplitSeq(query, "&") {
		name, value, _ := strings.Cut(sequence, "=")

		name = percent.Decode(name)
		value = percent.Decode(value)

		if _, ok := sequences[name]; !ok {
			names = append(names, name)
			size += len(name) + 1
		}
		sequences[name] = value
		size += len(value)
	}
	slices.Sort(names)

	var builder strings.Builder
	builder.Grow(size)

	for i, k := range names {
		if i > 0 {
			builder.WriteByte('&')
		}
		builder.WriteString(percent.Encode(k))
		builder.WriteByte('=')
		builder.WriteString(percent.Encode(sequences[k]))
	}
	return builder.String()
}

type CanonicalURL struct {
	value string
}

func (c *CanonicalURL) String() string {
	return c.value
}

func NormalizeURL(u *url.URL) CanonicalURL {
	scheme := u.Scheme
	if scheme == "" {
		scheme = "https"
	}
	host := u.Hostname()
	path := normalizePath(u.EscapedPath())
	query := normalizeQuery(u.RawQuery)

	var builder strings.Builder
	builder.Grow(
		len(scheme) +
			3 +
			len(host) +
			len(path) +
			1 +
			len(query),
	)

	builder.WriteString(scheme)
	builder.WriteString("://")
	builder.WriteString(host)
	builder.WriteString(path)
	builder.WriteByte('?')
	builder.WriteString(query)

	var c CanonicalURL
	c.value = builder.String()
	return c
}
