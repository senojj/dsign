package dsign

import (
	"runtime"
)

type Secret struct {
	data   []byte
	pinner runtime.Pinner
}

// NewSecret constructs a new Secret instance by copying the bytes from data.
// All bytes in data are replaced with zero before the Secret is returned. If
// the bytes in data originated from a string value, that value may remain in
// the application's global heap until the GC eventually cleans it up. Avoid
// storing the secret value as a string if this is a concern.
func NewSecret(data []byte) *Secret {
	s := new(Secret)
	if data == nil {
		return s
	}
	s.data = make([]byte, len(data))
	copy(s.data, data)
	s.pinner.Pin(s.data)

	for i := range len(data) {
		data[i] = 0
	}
	runtime.KeepAlive(data)
	return s
}

func (s *Secret) Clone() *Secret {
	clone := new(Secret)
	if s == nil || s.data == nil {
		return clone
	}
	clone.data = make([]byte, len(s.data))
	copy(clone.data, s.data)
	clone.pinner.Pin(clone.data)
	return clone
}

func (s *Secret) Expose() []byte {
	if s == nil {
		return nil
	}
	return s.data
}

func (s *Secret) Destroy() {
	if s == nil || s.data == nil {
		return
	}

	for i := range len(s.data) {
		s.data[i] = 0
	}
	runtime.KeepAlive(s.data)
	s.pinner.Unpin()
	s.data = nil
}

func (s *Secret) String() string {
	return "[SECRET]"
}
