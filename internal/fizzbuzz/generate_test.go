package fizzbuzz

import (
	"testing"
)

// Generate assumes validated inputs.
// Parameter validation is performed at the API layer to keep this function focused on pure business logic.
func TestGenerate_ClassicFizzBuzz(t *testing.T) {
	got := Generate(3, 5, 16, "fizz", "buzz")

	want := []string{
		"1", "2", "fizz", "4", "buzz", "fizz", "7", "8", "fizz", "buzz",
		"11", "fizz", "13", "14", "fizzbuzz", "16",
	}

	if len(got) != len(want) {
		t.Fatalf("len=%d want=%d", len(got), len(want))
	}

	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("index %d: got=%s want=%s", i, got[i], want[i])
		}
	}
}

func TestGenerate_SameMultiples(t *testing.T) {
	got := Generate(3, 3, 6, "a", "b")

	want := []string{"1", "2", "ab", "4", "5", "ab"}

	if len(got) != len(want) {
		t.Fatalf("len=%d want=%d", len(got), len(want))
	}

	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("index %d: got=%s want=%s", i, got[i], want[i])
		}
	}
}

func TestGenerate_LimitOne(t *testing.T) {
	got := Generate(2, 3, 1, "a", "b")

	want := []string{"1"}

	if len(got) != len(want) || got[0] != want[0] {
		t.Fatalf("got=%v want=%v", got, want)
	}
}

func TestGenerate_NoMultiples(t *testing.T) {
	got := Generate(50, 100, 10, "x", "y")

	want := []string{
		"1", "2", "3", "4", "5", "6", "7", "8", "9", "10",
	}

	if len(got) != len(want) {
		t.Fatalf("len=%d want=%d", len(got), len(want))
	}

	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("index %d: got=%s want=%s", i, got[i], want[i])
		}
	}
}

func TestGenerate_SpecialStrings(t *testing.T) {
	got := Generate(2, 3, 6, "📦", "🚚")

	want := []string{
		"1",
		"📦",
		"🚚",
		"📦",
		"5",
		"📦🚚",
	}

	if len(got) != len(want) {
		t.Fatalf("len=%d want=%d", len(got), len(want))
	}

	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("index %d: got=%s want=%s", i, got[i], want[i])
		}
	}
}
