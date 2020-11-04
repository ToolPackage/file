package utils

func OrString(str ...string) string {
	for _, s := range str {
		if s != "" {
			return s
		}
	}
	return ""
}

func Min(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

func ConvertInt16ToByte(v uint16, buf []byte, offset int) {
	buf[offset] = byte(v >> 8)
	buf[offset+1] = byte(v & 0xffff)
}

func ConvertByteToInt16(buf []byte, offset int) uint16 {
	return uint16(buf[offset])<<8 + uint16(buf[offset+1])
}
