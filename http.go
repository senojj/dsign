package dsign

import (
	"bytes"
	"crypto/sha256"
	"dsign/internal/encoding/hex"
	"io"
	"net/http"
)

func EncodeRequest(req *http.Request) ([]byte, error) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	req.Body = io.NopCloser(bytes.NewReader(body))
	canonicalURL := EncodeURL(req.URL)

	var buffer bytes.Buffer
	buffer.Grow(
		len(req.Method) +
			len(canonicalURL) +
			len(body),
	)
	buffer.WriteString(req.Method)
	buffer.Write(canonicalURL)
	buffer.Write(body)

	digest := sha256.New()
	digest.Write(buffer.Bytes())
	return hex.Encode(digest.Sum(nil), hex.Lower), nil
}

func SignRequest(signer *Signer, req *http.Request) error {
	signingString, err := EncodeRequest(req)
	if err != nil {
		return err
	}

	signature := signer.Sign(signingString)
	expiration := signature.Timestamp().String()
	accessKey := signature.AccessKey()

	req.Header.Set("dynata-expiration", expiration)
	req.Header.Set("dynata-access-key", accessKey)
	req.Header.Set("dynata-signature", signature.String())
	return nil
}

type SigningTransport struct {
	*http.Transport

	Signer *Signer
}

func (t *SigningTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	err := SignRequest(t.Signer, req)
	if err != nil {
		return nil, err
	}
	return t.Transport.RoundTrip(req)
}
