package encoding

var HexTable [256]uint16

func init() {
	const alphabet = "0123456789ABCDEF"
	for i := range 256 {
		HexTable[i] = uint16(alphabet[i&0x0F])<<8 | uint16(alphabet[i>>4])
	}
}
