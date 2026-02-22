package memory

import (
	"testing"

	"github.com/BryanRamires/FizzBuzz/internal/stats"
)

func TestRepoTop(t *testing.T) {
	r := New()

	k1 := stats.Key{Int1: 3, Int2: 5, Limit: 10, Str1: "fizz", Str2: "buzz"}
	k2 := stats.Key{Int1: 2, Int2: 7, Limit: 10, Str1: "a", Str2: "b"}

	r.Inc(k1)
	r.Inc(k1)
	r.Inc(k2)

	top, ok := r.Top()
	if !ok {
		t.Fatal("expected ok")
	}

	if top.Hits != 2 {
		t.Fatalf("hits=%d want=2", top.Hits)
	}
	if top.Parameters != k1 {
		t.Fatalf("params=%+v want=%+v", top.Parameters, k1)
	}
}
