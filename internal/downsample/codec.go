package downsample

import (
	"encoding/json"
	"time"
)

type wireBucket struct {
	Start int64   `json:"start"`
	End   int64   `json:"end"`
	Sum   float64 `json:"sum"`
	Max   float64 `json:"max"`
	Count int64   `json:"count"`
}

func EncodeBuckets(in []Bucket) ([]byte, error) {
	wire := make([]wireBucket, len(in))
	for i, b := range in {
		wire[i] = wireBucket{
			Start: b.Start.UnixNano(),
			End:   b.End.UnixNano(),
			Sum:   b.Sum,
			Max:   b.Max,
			Count: b.Count,
		}
	}
	return json.Marshal(wire)
}

func DecodeBuckets(data []byte) ([]Bucket, error) {
	var wire []wireBucket
	if err := json.Unmarshal(data, &wire); err != nil {
		return nil, err
	}
	out := make([]Bucket, len(wire))
	for i, w := range wire {
		out[i] = Bucket{
			Start: time.Unix(0, w.Start).UTC(),
			End:   time.Unix(0, w.End).UTC(),
			Sum:   w.Sum,
			Max:   w.Max,
			Count: w.Count,
			set:   w.Count > 0,
		}
	}
	return out, nil
}
