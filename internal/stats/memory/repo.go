package memory

import (
	"sync"

	"github.com/BryanRamires/FizzBuzz/internal/stats"
)

var _ stats.Repository = (*Repo)(nil)

type Repo struct {
	mu      sync.RWMutex
	m       map[stats.Key]uint64
	maxKeys int
}

// defaultMaxKeys limits the number of distinct stats entries
// to prevent unbounded memory growth.
const defaultMaxKeys = 100_000

func New() *Repo {
	return &Repo{m: make(map[stats.Key]uint64), maxKeys: defaultMaxKeys}
}

func (r *Repo) Inc(k stats.Key) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.m[k]; !ok && r.maxKeys > 0 && len(r.m) >= r.maxKeys {
		return // best-effort: ignore new keys when full
	}
	r.m[k]++
}

// Linear scan is acceptable: map size is bounded and this is a cold path.
func (r *Repo) Top() (stats.Top, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var (
		bestK   stats.Key
		bestV   uint64
		hasBest bool
	)
	for k, v := range r.m {
		if !hasBest || v > bestV {
			bestK, bestV, hasBest = k, v, true
		}
	}
	if !hasBest {
		return stats.Top{}, false
	}
	return stats.Top{Parameters: bestK, Hits: bestV}, true
}
