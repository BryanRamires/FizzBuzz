package memory

import (
	"context"
	"encoding/json"
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

func (r *Repo) Inc(ctx context.Context, k stats.Key) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.m[k]; !ok && r.maxKeys > 0 && len(r.m) >= r.maxKeys {
		return nil // best-effort: ignore new keys when full
	}
	r.m[k]++
	return nil
}

// Linear scan is acceptable: map size is bounded and this is a cold path.
func (r *Repo) Top(ctx context.Context) (stats.Top, bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var (
		bestK   stats.Key
		bestV   uint64
		hasBest bool
		bestMem string
	)

	for k, v := range r.m {
		if !hasBest || v > bestV {
			bestK, bestV, hasBest = k, v, true
			bestMem = keyMember(k)
			continue
		}
		// Match Redis ZREVRANGE tie-breaking to unify back-end:
		// pick the lexicographically highest member.
		if v == bestV {
			mem := keyMember(k)
			if mem > bestMem {
				bestK = k
				bestMem = mem
			}
		}
	}

	if !hasBest {
		return stats.Top{}, false, nil
	}
	return stats.Top{Parameters: bestK, Hits: bestV}, true, nil
}

func keyMember(k stats.Key) string {
	// Top() is a cold path, so the allocation cost is acceptable.
	b, _ := json.Marshal(k)
	return string(b)
}
