package hex

import (
	"encoding/binary"
)

var UpperTable [256]uint16
var LowerTable [256]uint16

func init() {
	const upperAlphabet = "0123456789ABCDEF"
	for i := range 256 {
		UpperTable[i] = uint16(upperAlphabet[i&0x0F])<<8 | uint16(upperAlphabet[i>>4])
	}

	const lowerAlphabet = "0123456789abcdef"
	for i := range 256 {
		LowerTable[i] = uint16(lowerAlphabet[i&0x0F])<<8 | uint16(lowerAlphabet[i>>4])
	}
}

type Case byte

const (
	Upper Case = iota
	Lower
)

func Encode(data []byte, alphabet Case) []byte {
	if len(data) == 0 {
		return data
	}

	var table [256]uint16
	switch alphabet {
	case Upper:
		table = UpperTable
	case Lower:
		table = LowerTable
	default:
		table = UpperTable
	}

	src := data
	out := make([]byte, len(src)*2)

	dst := out

	for len(src) >= 8 {
		p := (*[16]byte)(dst)

		w := uint64(table[src[0]]) |
			uint64(table[src[1]])<<16 |
			uint64(table[src[2]])<<32 |
			uint64(table[src[3]])<<48

		binary.LittleEndian.PutUint64(p[:8], w)

		w = uint64(table[src[4]]) |
			uint64(table[src[5]])<<16 |
			uint64(table[src[6]])<<32 |
			uint64(table[src[7]])<<48

		binary.LittleEndian.PutUint64(p[8:], w)

		src = src[8:]
		dst = dst[16:]
	}

	if len(src) >= 4 {
		p := (*[8]byte)(dst)

		w := uint64(table[src[0]]) |
			uint64(table[src[1]])<<16 |
			uint64(table[src[2]])<<32 |
			uint64(table[src[3]])<<48

		binary.LittleEndian.PutUint64(p[:8], w)

		src = src[4:]
		dst = dst[8:]
	}

	switch len(src) {
	case 3:
		p := (*[6]byte)(dst)

		w := uint32(table[src[0]]) |
			uint32(table[src[1]])<<16

		binary.LittleEndian.PutUint32(p[:4], w)
		binary.LittleEndian.PutUint16(p[4:], table[src[2]])
	case 2:
		p := (*[4]byte)(dst)

		w := uint32(table[src[0]]) |
			uint32(table[src[1]])<<16

		binary.LittleEndian.PutUint32(p[:4], w)
	case 1:
		p := (*[2]byte)(dst)

		binary.LittleEndian.PutUint16(p[:2], table[src[0]])
	}

	return out
}
