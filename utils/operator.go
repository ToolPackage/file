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

func ConvertUint16ToByte(v uint16, buf []byte, offset int) {
	buf[offset] = byte(v >> 8 & 0xff)
	buf[offset+1] = byte(v & 0xff)
}

func ConvertUint32ToByte(v uint32, buf []byte, offset int) {
	buf[offset] = byte(v >> 24 & 0xff)
	buf[offset+1] = byte(v >> 16 & 0xff)
	buf[offset+2] = byte(v >> 8 & 0xff)
	buf[offset+3] = byte(v & 0xff)
}

func ConvertInt64ToByte(v int64, buf []byte, offset int) {
	buf[offset] = byte(v >> 56 & 0xff)
	buf[offset+1] = byte(v >> 48 & 0xff)
	buf[offset+2] = byte(v >> 40 & 0xff)
	buf[offset+3] = byte(v >> 32 & 0xff)
	buf[offset+4] = byte(v >> 24 & 0xff)
	buf[offset+5] = byte(v >> 16 & 0xff)
	buf[offset+6] = byte(v >> 8 & 0xff)
	buf[offset+7] = byte(v & 0xff)
}

func ConvertByteToUint16(buf []byte, offset int) uint16 {
	return uint16(buf[offset])<<8 + uint16(buf[offset+1])
}

func ConvertByteToUint32(buf []byte, offset int) uint32 {
	var v uint32
	v += uint32(buf[offset]) << 24
	v += uint32(buf[offset+1]) << 16
	v += uint32(buf[offset+2]) << 8
	v += uint32(buf[offset+3])
	return v
}

func ConvertByteToInt64(buf []byte, offset int) int64 {
	var v int64
	v += int64(buf[offset]) << 56
	v += int64(buf[offset+1]) << 48
	v += int64(buf[offset+2]) << 40
	v += int64(buf[offset+3]) << 32
	v += int64(buf[offset+4]) << 24
	v += int64(buf[offset+5]) << 16
	v += int64(buf[offset+6]) << 8
	v += int64(buf[offset+7])
	return v
}
