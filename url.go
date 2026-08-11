package dsign

import (
	"bytes"
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

func EncodeQuery(query string) []byte {
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

	var buffer bytes.Buffer
	buffer.Grow(size)

	for i, k := range names {
		if i > 0 {
			buffer.WriteByte('&')
		}
		buffer.WriteString(percent.Encode(k))
		buffer.WriteByte('=')
		buffer.WriteString(percent.Encode(sequences[k]))
	}
	return buffer.Bytes()
}

func EncodeURL(u *url.URL) []byte {
	scheme := u.Scheme
	if scheme == "" {
		scheme = "https"
	}
	host := u.Hostname()
	path := normalizePath(u.EscapedPath())
	query := EncodeQuery(u.RawQuery)

	var buffer bytes.Buffer
	buffer.Grow(
		len(scheme) +
			3 +
			len(host) +
			len(path) +
			1 +
			len(query),
	)

	buffer.WriteString(scheme)
	buffer.WriteString("://")
	buffer.WriteString(host)
	buffer.WriteString(path)
	buffer.WriteByte('?')
	buffer.Write(query)

	return buffer.Bytes()
}

func SignURL(signer Signer, u *url.URL) {
	signingString := EncodeURL(u)
	signature := signer.Sign(signingString)
	expiration := signature.Timestamp().String()
	accessKey := signature.AccessKey()

	query := u.Query()
	query.Set("dynata-expiration", expiration)
	query.Set("dynata-access-key", accessKey)
	query.Set("dynata-signature", signature.String())
	u.RawQuery = query.Encode()
}
