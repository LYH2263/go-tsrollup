package buffer

// FloatRing 定长浮点环形缓冲。
type FloatRing struct {
	buf []float64
	pos int
	n   int
}

func NewFloatRing(cap int) *FloatRing {
	if cap <= 0 {
		cap = 64
	}
	return &FloatRing{buf: make([]float64, cap)}
}

func (r *FloatRing) Push(v float64) {
	r.buf[r.pos] = v
	r.pos = (r.pos + 1) % len(r.buf)
	if r.n < len(r.buf) {
		r.n++
	}
}

func (r *FloatRing) Len() int { return r.n }

func (r *FloatRing) Snapshot() []float64 {
	out := make([]float64, r.n)
	start := (r.pos - r.n + len(r.buf)) % len(r.buf)
	for i := 0; i < r.n; i++ {
		out[i] = r.buf[(start+i)%len(r.buf)]
	}
	return out
}

func (r *FloatRing) Sum() float64 {
	var s float64
	for i := 0; i < r.n; i++ {
		idx := (r.pos - r.n + i + len(r.buf)) % len(r.buf)
		s += r.buf[idx]
	}
	return s
}
