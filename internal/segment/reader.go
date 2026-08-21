package segment

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"

	"github.com/LYH2263/go-tsrollup/internal/errs"
)

// View 只读段视图。
type View struct {
	path string
	rows []Row
}

func OpenView(path string) (*View, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	rows, err := decode(data)
	if err != nil {
		return nil, err
	}
	return &View{path: path, rows: rows}, nil
}

func decode(data []byte) ([]Row, error) {
	if len(data) < hdrSize {
		return nil, errs.Corrupt("short header")
	}
	if !checkMagic(data) {
		return nil, errs.Corrupt("bad magic")
	}
	ver := binary.LittleEndian.Uint32(data[4:8])
	if ver != version {
		return nil, errs.Corrupt(fmt.Sprintf("bad version %d", ver))
	}
	wantCRC := binary.LittleEndian.Uint32(data[12:16])
	payload := data[hdrSize:]
	if checksum(payload) != wantCRC {
		return nil, errs.Corrupt("checksum mismatch")
	}
	var rows []Row
	if err := json.Unmarshal(payload, &rows); err != nil {
		return nil, errs.Corrupt("json: " + err.Error())
	}
	return rows, nil
}

func (v *View) RowsFor(series string) []Row {
	out := make([]Row, 0)
	for _, r := range v.rows {
		if r.Series == series {
			out = append(out, r)
		}
	}
	return out
}

func (v *View) All() []Row {
	out := make([]Row, len(v.rows))
	copy(out, v.rows)
	return out
}

func (v *View) Path() string { return v.path }

func (v *View) Close() error { return nil }

func (v *View) Len() int { return len(v.rows) }
