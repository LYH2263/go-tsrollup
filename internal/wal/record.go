package wal

import "time"

// Record WAL 记录。
type Record struct {
	Series string            `json:"series"`
	Name   string            `json:"name"`
	Labels map[string]string `json:"labels"`
	Ts     time.Time         `json:"ts"`
	Value  float64           `json:"value"`
}
