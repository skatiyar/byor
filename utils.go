package byor

func appendVarint(buf []byte, n int32) []byte {
	temp := make([]byte, 4)
	temp[0] = byte((n >> 24) & 0xFF)
	temp[1] = byte((n >> 16) & 0xFF)
	temp[2] = byte((n >> 8) & 0xFF)
	temp[3] = byte((n >> 0) & 0xFF)
	return append(buf, temp...)
}

func getVarint(buf []byte) int32 {
	return int32((buf[3] & 0xFF) | (buf[2] & 0xFF) | (buf[1] & 0xFF) | (buf[0] & 0xFF))
}
