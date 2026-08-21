package segment

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Writer 原子写段：写临时文件 → Sync → Rename。
type Writer struct {
	final string
	tmp   string
	f     *os.File
	rows  []Row
}

func Create(path string) (*Writer, error) {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	tmp := path + ".tmp"
	f, err := os.OpenFile(tmp, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, err
	}
	return &Writer{final: path, tmp: tmp, f: f}, nil
}

func (w *Writer) WriteRow(r Row) error {
	w.rows = append(w.rows, r)
	return nil
}

func (w *Writer) Abort() error {
	if w.f != nil {
		_ = w.f.Close()
		w.f = nil
	}
	_ = os.Remove(w.tmp)
	return nil
}

func (w *Writer) Close() error {
	if w.f == nil {
		return fmt.Errorf("segment: writer closed")
	}
	payload, err := json.Marshal(w.rows)
	if err != nil {
		_ = w.Abort()
		return err
	}
	hdr := make([]byte, hdrSize)
	putHeader(hdr, uint32(len(w.rows)))
	crc := checksum(payload)
	binary.LittleEndian.PutUint32(hdr[12:16], crc)
	if _, err := w.f.Write(hdr); err != nil {
		_ = w.Abort()
		return err
	}
	if _, err := w.f.Write(payload); err != nil {
		_ = w.Abort()
		return err
	}
	// 必须 Sync 后再 Rename（bug09 正确行为）。
	if err := w.f.Sync(); err != nil {
		_ = w.Abort()
		return err
	}
	if err := w.f.Close(); err != nil {
		_ = os.Remove(w.tmp)
		return err
	}
	w.f = nil
	if err := os.Rename(w.tmp, w.final); err != nil {
		_ = os.Remove(w.tmp)
		return err
	}
	return nil
}
