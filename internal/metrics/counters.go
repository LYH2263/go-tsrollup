package metrics

import "sync"

// Counter 命名计数器。
type Counter struct {
	mu sync.Mutex
	m  map[string]int64
}

func NewCounter() *Counter {
	return &Counter{m: make(map[string]int64)}
}

func (c *Counter) Add(name string, n int64) {
	c.mu.Lock()
	c.m[name] += n
	c.mu.Unlock()
}

func (c *Counter) Get(name string) int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.m[name]
}

func (c *Counter) All() map[string]int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	out := make(map[string]int64, len(c.m))
	for k, v := range c.m {
		out[k] = v
	}
	return out
}
