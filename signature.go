package dsign

import (
	"crypto/hmac"
	"crypto/sha256"
	"dsign/internal/encoding/hex"
	"time"
	"unsafe"
)

const timestampFormat string = "2006-01-02T15:04:05.000Z"

type Timestamp struct {
	time.Time
}

func (e Timestamp) String() string {
	return e.Format(timestampFormat)
}

type Signature struct {
	data []byte
	key  string
	ts   Timestamp
}

func (s *Signature) AccessKey() string {
	return s.key
}

func (s *Signature) Timestamp() Timestamp {
	return s.ts
}

func (s *Signature) String() string {
	return unsafe.String(unsafe.SliceData(s.data), len(s.data))
}

func digest(key, value []byte) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write(value)
	return hex.Encode(mac.Sum(nil), hex.Lower)
}

type Signer struct {
	AccessKey string
	SecretKey *Secret
	TimeFunc  func() time.Time
}

func (s *Signer) Close() error {
	s.SecretKey.Destroy()
	return nil
}

func (s *Signer) Sign(payload []byte) Signature {
	var ts Timestamp
	if s.TimeFunc != nil {
		ts = Timestamp{s.TimeFunc()}
	} else {
		ts = Timestamp{time.Now()}
	}

	first := digest([]byte(ts.String()), payload)
	second := digest(unsafe.Slice(unsafe.StringData(s.AccessKey), len(s.AccessKey)), first)
	third := digest(s.SecretKey.Expose(), second)

	signature := Signature{data: third, key: s.AccessKey, ts: ts}

	return signature
}
