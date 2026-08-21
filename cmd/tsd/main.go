package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/LYH2263/go-tsrollup"
)

func main() {
	addr := flag.String("addr", ":8113", "listen address")
	root := flag.String("root", "./data", "engine root")
	web := flag.String("web", "web", "static web dir")
	window := flag.Duration("window", time.Minute, "tumbling window size")
	aggKind := flag.String("agg", "sum", "sum|avg|max")
	flag.Parse()

	eng, err := tsrollup.Open(tsrollup.Options{
		Root:    *root,
		Window:  tsrollup.WindowSpec{Size: *window},
		AggKind: *aggKind,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer eng.Close()

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(*web)))
	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, map[string]any{"ok": true, "time": time.Now().UTC()})
	})
	mux.HandleFunc("/api/stats", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, eng.Stats())
	})
	mux.HandleFunc("/api/series", func(w http.ResponseWriter, r *http.Request) {
		list, err := eng.ListSeries()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		writeJSON(w, list)
	})
	mux.HandleFunc("/api/append", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "POST only", 405)
			return
		}
		var body struct {
			Series string            `json:"series"`
			Labels map[string]string `json:"labels"`
			Ts     string            `json:"ts"`
			Value  float64           `json:"value"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		ts := time.Now().UTC()
		if body.Ts != "" {
			parsed, err := time.Parse(time.RFC3339Nano, body.Ts)
			if err != nil {
				http.Error(w, err.Error(), 400)
				return
			}
			ts = parsed
		}
		if err := eng.AppendContext(r.Context(), tsrollup.Sample{
			Series: body.Series,
			Labels: body.Labels,
			Ts:     ts,
			Value:  body.Value,
		}); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		writeJSON(w, map[string]string{"status": "ok"})
	})
	mux.HandleFunc("/api/query", func(w http.ResponseWriter, r *http.Request) {
		series := r.URL.Query().Get("series")
		from, to := parseRange(r)
		res, err := eng.Query(series, nil, from, to)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		writeJSON(w, res)
	})
	mux.HandleFunc("/api/compact", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "POST only", 405)
			return
		}
		if err := eng.Compact(r.Context()); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		writeJSON(w, map[string]string{"status": "ok"})
	})
	mux.HandleFunc("/api/snapshot", func(w http.ResponseWriter, r *http.Request) {
		snap, err := eng.Snapshot()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		writeJSON(w, snap)
	})

	log.Printf("tsd listening on %s root=%s", *addr, *root)
	if err := http.ListenAndServe(*addr, mux); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func parseRange(r *http.Request) (time.Time, time.Time) {
	var from, to time.Time
	if s := r.URL.Query().Get("from"); s != "" {
		if n, err := strconv.ParseInt(s, 10, 64); err == nil {
			from = time.Unix(0, n).UTC()
		} else if t, err := time.Parse(time.RFC3339Nano, s); err == nil {
			from = t
		}
	}
	if s := r.URL.Query().Get("to"); s != "" {
		if n, err := strconv.ParseInt(s, 10, 64); err == nil {
			to = time.Unix(0, n).UTC()
		} else if t, err := time.Parse(time.RFC3339Nano, s); err == nil {
			to = t
		}
	}
	return from, to
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}

var _ = context.Background
