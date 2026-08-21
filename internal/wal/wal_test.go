package wal_test

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/LYH2263/go-tsrollup/internal/wal"
)

func TestAppendReplay(t *testing.T) {
	dir := t.TempDir()
	l, err := wal.Open(filepath.Join(dir, "wal"))
	if err != nil {
		t.Fatal(err)
	}
	rec := wal.Record{
		Series: "a",
		Name:   "a",
		Ts:     time.Unix(1, 0).UTC(),
		Value:  1.5,
	}
	if err := l.Append(context.Background(), rec); err != nil {
		t.Fatal(err)
	}
	if err := l.Flush(); err != nil {
		t.Fatal(err)
	}
	recs, err := l.Replay()
	if err != nil {
		t.Fatal(err)
	}
	if len(recs) != 1 || recs[0].Value != 1.5 {
		t.Fatalf("%+v", recs)
	}
	_ = l.Close()
}

func TestAppendHonorsCancel(t *testing.T) {
	l, err := wal.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err = l.Append(ctx, wal.Record{Series: "x", Ts: time.Now().UTC()})
	if err == nil {
		t.Fatal("expected cancel")
	}
}
