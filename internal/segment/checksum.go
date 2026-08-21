package segment

// VerifyPayload 校验 payload 与期望 crc。
func VerifyPayload(payload []byte, want uint32) bool {
	return checksum(payload) == want
}

// CRCOf 计算 payload CRC。
func CRCOf(payload []byte) uint32 {
	return checksum(payload)
}
