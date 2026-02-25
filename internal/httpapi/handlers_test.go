package httpapi

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/BryanRamires/FizzBuzz/internal/config"
	"github.com/BryanRamires/FizzBuzz/internal/stats"
	"github.com/BryanRamires/FizzBuzz/internal/stats/memory"
)

func TestFizzBuzz_OK(t *testing.T) {
	req := httptest.NewRequest(
		http.MethodGet,
		"/fizzbuzz?int1=3&int2=5&limit=16&str1=fizz&str2=buzz",
		nil,
	)

	cfg, _ := config.New()
	repo := memory.New()
	svc, _ := stats.NewService(repo)
	h := NewHandler(cfg, nil, svc)

	rr := httptest.NewRecorder()
	NewRouter(cfg, testLogger(), h).ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status=%d want=%d", rr.Code, http.StatusOK)
	}

	var got []string
	if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
		t.Fatal(err)
	}

	if len(got) != 16 {
		t.Fatalf("len=%d want=16", len(got))
	}

	// we only assert the most critical case: 15 is a multiple of both 3 and 5
	// verifying the full sequence would be redundant since it is already covered by unit tests
	if got[14] != "fizzbuzz" {
		t.Fatalf("got[14]=%q want=%q", got[14], "fizzbuzz")
	}
}

func TestFizzBuzz_InvalidInputs(t *testing.T) {
	tests := []struct {
		name string
		url  string
	}{
		{"missing int1", "/fizzbuzz?int2=5&limit=10&str1=a&str2=b"},
		{"zero int1", "/fizzbuzz?int1=0&int2=5&limit=10&str1=a&str2=b"},
		{"limit too large", "/fizzbuzz?int1=3&int2=5&limit=999999&str1=a&str2=b"},
		{"empty str1", "/fizzbuzz?int1=3&int2=5&limit=10&str1=&str2=b"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, _ := config.New()
			repo := memory.New()
			svc, _ := stats.NewService(repo)
			h := NewHandler(cfg, nil, svc)

			rr := httptest.NewRecorder()
			req := httptest.NewRequest("GET", tt.url, nil)

			NewRouter(cfg, testLogger(), h).ServeHTTP(rr, req)

			if rr.Code != http.StatusBadRequest {
				t.Fatalf("status=%d want=400", rr.Code)
			}
		})
	}
}

func TestStats_AfterFizzBuzz_ReturnsTop(t *testing.T) {
	cfg, _ := config.New()
	repo := memory.New()
	svc, _ := stats.NewService(repo)
	h := NewHandler(cfg, nil, svc)

	srv := httptest.NewServer(NewRouter(cfg, testLogger(), h))
	defer srv.Close()

	// hit twice
	_, _ = http.Get(srv.URL + "/fizzbuzz?int1=3&int2=5&limit=16&str1=fizz&str2=buzz")
	_, _ = http.Get(srv.URL + "/fizzbuzz?int1=3&int2=5&limit=16&str1=fizz&str2=buzz")

	resp, err := http.Get(srv.URL + "/stats")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			t.Errorf("failed to close response body: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status=%d want=%d", resp.StatusCode, http.StatusOK)
	}

	var got struct {
		Parameters struct {
			Int1  int    `json:"int1"`
			Int2  int    `json:"int2"`
			Limit int    `json:"limit"`
			Str1  string `json:"str1"`
			Str2  string `json:"str2"`
		} `json:"parameters"`
		Hits uint64 `json:"hits"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Fatal(err)
	}

	if got.Hits != 2 {
		t.Fatalf("hits=%d want=2", got.Hits)
	}
	if got.Parameters.Int1 != 3 || got.Parameters.Int2 != 5 || got.Parameters.Limit != 16 {
		t.Fatalf("bad params: %+v", got.Parameters)
	}
	if got.Parameters.Str1 != "fizz" || got.Parameters.Str2 != "buzz" {
		t.Fatalf("bad strings: %+v", got.Parameters)
	}
}

func testLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(io.Discard, nil))
}
