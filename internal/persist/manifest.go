package persist

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Manifest 简单清单。
type Manifest struct {
	Segments int `json:"segments"`
}

func (d *Dir) WriteManifest(n int) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	path := filepath.Join(d.root, "manifest.json")
	tmp := path + ".tmp"
	b, err := json.Marshal(Manifest{Segments: n})
	if err != nil {
		return err
	}
	if err := os.WriteFile(tmp, b, 0o644); err != nil {
		return err
	}
	f, err := os.OpenFile(tmp, os.O_RDWR, 0)
	if err != nil {
		return err
	}
	if err := f.Sync(); err != nil {
		_ = f.Close()
		return err
	}
	_ = f.Close()
	return os.Rename(tmp, path)
}

func (d *Dir) ReadManifest() (Manifest, error) {
	var m Manifest
	b, err := os.ReadFile(filepath.Join(d.root, "manifest.json"))
	if err != nil {
		return m, err
	}
	err = json.Unmarshal(b, &m)
	return m, err
}
