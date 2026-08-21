package segment

import (
	"encoding/binary"
	"hash/crc32"
	"time"
)

const (
	magic   = 0x5453524c // TSRL
	version = 1
	hdrSize = 16
)

// Row 段内一行窗口结果。
type Row struct {
	Series string
	Start  time.Time
	End    time.Time
	Value  float64
	Count  int64
	Agg    string
}

func putHeader(buf []byte, nrows uint32) {
	binary.LittleEndian.PutUint32(buf[0:4], magic)
	binary.LittleEndian.PutUint32(buf[4:8], version)
	binary.LittleEndian.PutUint32(buf[8:12], nrows)
	binary.LittleEndian.PutUint32(buf[12:16], 0) // reserved / crc placeholder
}

func checkMagic(buf []byte) bool {
	if len(buf) < 4 {
		return false
	}
	return binary.LittleEndian.Uint32(buf[0:4]) == magic
}

func checksum(data []byte) uint32 {
	return crc32.ChecksumIEEE(data)
}
