package stats

import "testing"

// petit fake repo in-memory pour tester le service
type fakeRepo struct {
	m map[Key]uint64
}

func newFakeRepo() *fakeRepo { return &fakeRepo{m: map[Key]uint64{}} }

func (f *fakeRepo) Inc(k Key) { f.m[k]++ }

func (f *fakeRepo) Top() (Top, bool) {
	var (
		bestK   Key
		bestV   uint64
		hasBest bool
	)
	for k, v := range f.m {
		if !hasBest || v > bestV {
			bestK, bestV, hasBest = k, v, true
		}
	}
	if !hasBest {
		return Top{}, false
	}
	return Top{Parameters: bestK, Hits: bestV}, true
}

func TestService_RecordAndMostFrequent(t *testing.T) {
	repo := newFakeRepo()
	svc, _ := NewService(repo)

	k := Key{Int1: 3, Int2: 5, Limit: 16, Str1: "fizz", Str2: "buzz"}

	svc.Record(k)
	svc.Record(k)

	top, ok := svc.MostFrequent()
	if !ok {
		t.Fatal("expected ok")
	}
	if top.Hits != 2 {
		t.Fatalf("hits=%d want=2", top.Hits)
	}
	if top.Parameters != k {
		t.Fatalf("params=%+v want=%+v", top.Parameters, k)
	}
}
