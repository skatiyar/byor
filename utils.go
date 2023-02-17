package byor

func appendVarint(buf []byte, n int32) []byte {
	temp := make([]byte, 4)
	for i := 0; i < len(temp); i += 1 {
		temp[i] = byte((n >> (i * 8) & 0xFF))
	}
	return append(buf, temp...)
}

func getVarint(buf []byte) int32 {
	return int32((buf[3] & 0xFF) | (buf[2] & 0xFF) | (buf[1] & 0xFF) | (buf[0] & 0xFF))
}
