package admin

import (
	"time"
)

// Dashboard 管理页摘要。
type Dashboard struct {
	Series    int       `json:"series"`
	Samples   int64     `json:"samples"`
	Segments  int       `json:"segments"`
	OpenWins  int       `json:"open_windows"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewDashboard(series, segments, openWins int, samples int64) Dashboard {
	return Dashboard{
		Series:    series,
		Samples:   samples,
		Segments:  segments,
		OpenWins:  openWins,
		UpdatedAt: time.Now().UTC(),
	}
}
