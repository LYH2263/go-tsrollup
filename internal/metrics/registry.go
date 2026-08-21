package metrics

import "sync/atomic"

// Registry 简单计数器集合。
type Registry struct {
	samples     atomic.Int64
	compactions atomic.Int64
	queries     atomic.Int64
	errors      atomic.Int64
	walBytes    atomic.Int64
}

func NewRegistry() *Registry { return &Registry{} }

func (r *Registry) IncSamples(n int64)     { r.samples.Add(n) }
func (r *Registry) IncCompactions(n int64) { r.compactions.Add(n) }
func (r *Registry) IncQueries(n int64)     { r.queries.Add(n) }
func (r *Registry) IncErrors(n int64)      { r.errors.Add(n) }
func (r *Registry) AddWALBytes(n int64)    { r.walBytes.Add(n) }

func (r *Registry) Snapshot() map[string]int64 {
	return map[string]int64{
		"samples":     r.samples.Load(),
		"compactions": r.compactions.Load(),
		"queries":     r.queries.Load(),
		"errors":      r.errors.Load(),
		"wal_bytes":   r.walBytes.Load(),
	}
}
