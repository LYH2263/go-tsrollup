package segment

import "encoding/binary"

// Header 段文件头解析结果。
type Header struct {
	Magic   uint32
	Version uint32
	Rows    uint32
	CRC     uint32
}

func ParseHeader(buf []byte) (Header, error) {
	var h Header
	if len(buf) < hdrSize {
		return h, errShort
	}
	h.Magic = binary.LittleEndian.Uint32(buf[0:4])
	h.Version = binary.LittleEndian.Uint32(buf[4:8])
	h.Rows = binary.LittleEndian.Uint32(buf[8:12])
	h.CRC = binary.LittleEndian.Uint32(buf[12:16])
	return h, nil
}

var errShort = errString("segment: short header")

type errString string

func (e errString) Error() string { return string(e) }
