package encoding

import (
	"encoding/binary"
	"math"
)

func AppendUvarint(b []byte, x uint64) []byte {
	var tmp [binary.MaxVarintLen64]byte
	n := binary.PutUvarint(tmp[:], x)
	return append(b, tmp[:n]...)
}

func AppendFloat64(b []byte, v float64) []byte {
	var tmp [8]byte
	binary.LittleEndian.PutUint64(tmp[:], math.Float64bits(v))
	return append(b, tmp[:]...)
}

func ReadFloat64(b []byte) (float64, []byte, bool) {
	if len(b) < 8 {
		return 0, b, false
	}
	u := binary.LittleEndian.Uint64(b[:8])
	return math.Float64frombits(u), b[8:], true
}

func AppendInt64(b []byte, v int64) []byte {
	var tmp [8]byte
	binary.LittleEndian.PutUint64(tmp[:], uint64(v))
	return append(b, tmp[:]...)
}

func ReadInt64(b []byte) (int64, []byte, bool) {
	if len(b) < 8 {
		return 0, b, false
	}
	u := binary.LittleEndian.Uint64(b[:8])
	return int64(u), b[8:], true
}
