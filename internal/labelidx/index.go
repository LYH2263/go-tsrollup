package labelidx

import (
	"sync"

	"github.com/LYH2263/go-tsrollup/internal/clone"
)

// Index 标签倒排索引。
type Index struct {
	mu      sync.RWMutex
	byID    map[string]map[string]string
	byLabel map[string]map[string]map[string]struct{} // key -> val -> series set
	byName  map[string]map[string]struct{}
}

func New() *Index {
	return &Index{
		byID:    make(map[string]map[string]string),
		byLabel: make(map[string]map[string]map[string]struct{}),
		byName:  make(map[string]map[string]struct{}),
	}
}

func (idx *Index) Add(id string, labels map[string]string) {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	// 深拷贝：倒排索引不得与调用方 map 共享底层，否则外部事后改动会穿透进 byID。
	labs := clone.Labels(labels)
	idx.byID[id] = labs
	name := labs["__name__"]
	if name == "" {
		// id 前缀到 { 之前视作 name
		for i := 0; i < len(id); i++ {
			if id[i] == '{' {
				name = id[:i]
				break
			}
		}
		if name == "" {
			name = id
		}
	}
	if idx.byName[name] == nil {
		idx.byName[name] = make(map[string]struct{})
	}
	idx.byName[name][id] = struct{}{}
	for k, v := range labs {
		if idx.byLabel[k] == nil {
			idx.byLabel[k] = make(map[string]map[string]struct{})
		}
		if idx.byLabel[k][v] == nil {
			idx.byLabel[k][v] = make(map[string]struct{})
		}
		idx.byLabel[k][v][id] = struct{}{}
	}
}

func (idx *Index) LabelsOf(id string) map[string]string {
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	return clone.Labels(idx.byID[id])
}

func (idx *Index) Len() int {
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	return len(idx.byID)
}
