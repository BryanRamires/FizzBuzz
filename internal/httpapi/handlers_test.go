package httpapi

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/BryanRamires/FizzBuzz/internal/config"
	"github.com/BryanRamires/FizzBuzz/internal/stats"
	"github.com/BryanRamires/FizzBuzz/internal/stats/memory"
)

// --- helpers ---

func testLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(io.Discard, nil))
}

func newTestRouter(t *testing.T) (config.Config, http.Handler) {
	t.Helper()

	cfg, err := config.New()
	if err != nil {
		t.Fatal(err)
	}

	cfg.MaxLimit = 50
	cfg.MaxStrLen = 16
	cfg.RedisEnabled = false

	repo := memory.New()
	svc, err := stats.NewService(repo)
	if err != nil {
		t.Fatal(err)
	}

	h := NewHandler(cfg, testLogger(), nil, svc)
	return cfg, NewRouter(cfg, testLogger(), h)
}

func doGET(t *testing.T, router http.Handler, path string) *httptest.ResponseRecorder {
	t.Helper()
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, path, nil)
	router.ServeHTTP(rr, req)
	return rr
}

// --- tests ---

func TestFizzBuzz_OK(t *testing.T) {
	_, router := newTestRouter(t)

	rr := doGET(t, router, "/fizzbuzz?int1=3&int2=5&limit=16&str1=fizz&str2=buzz")
	if rr.Code != http.StatusOK {
		t.Fatalf("status=%d want=%d body=%s", rr.Code, http.StatusOK, rr.Body.String())
	}
	ct := rr.Result().Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "application/json") {
		t.Fatalf("content-type=%q want prefix application/json", ct)
	}

	var got []string
	if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v body=%s", err, rr.Body.String())
	}
	if len(got) != 16 {
		t.Fatalf("len=%d want=16", len(got))
	}
	// critical case: 15 is multiple of both 3 and 5
	if got[14] != "fizzbuzz" {
		t.Fatalf("got[14]=%q want=%q", got[14], "fizzbuzz")
	}
}

func TestFizzBuzz_InvalidInputs_Table(t *testing.T) {
	cfg, router := newTestRouter(t)

	tooLong := strings.Repeat("a", cfg.MaxStrLen+1)
	okLen := strings.Repeat("a", cfg.MaxStrLen)

	tests := []struct {
		name            string
		path            string
		wantStatus      int
		wantErrContains string
	}{
		// missing params
		{"missing_int1", "/fizzbuzz?int2=5&limit=10&str1=a&str2=b", 400, "int1"},
		{"missing_int2", "/fizzbuzz?int1=3&limit=10&str1=a&str2=b", 400, "int2"},
		{"missing_limit", "/fizzbuzz?int1=3&int2=5&str1=a&str2=b", 400, "limit"},
		{"missing_str1", "/fizzbuzz?int1=3&int2=5&limit=10&str2=b", 400, "str1"},
		{"missing_str2", "/fizzbuzz?int1=3&int2=5&limit=10&str1=a", 400, "str2"},

		// ints invalid
		{"int1_not_int", "/fizzbuzz?int1=x&int2=5&limit=10&str1=a&str2=b", 400, "int1"},
		{"int2_not_int", "/fizzbuzz?int1=3&int2=x&limit=10&str1=a&str2=b", 400, "int2"},
		{"limit_not_int", "/fizzbuzz?int1=3&int2=5&limit=x&str1=a&str2=b", 400, "limit"},
		{"int1_zero", "/fizzbuzz?int1=0&int2=5&limit=10&str1=a&str2=b", 400, "int1"},
		{"int2_negative", "/fizzbuzz?int1=3&int2=-5&limit=10&str1=a&str2=b", 400, "int2"},
		{"limit_zero", "/fizzbuzz?int1=3&int2=5&limit=0&str1=a&str2=b", 400, "limit"},

		// limit boundaries
		{"limit_too_large", fmt.Sprintf("/fizzbuzz?int1=3&int2=5&limit=%d&str1=a&str2=b", cfg.MaxLimit+1), 400, "limit"},
		{"limit_max_ok", fmt.Sprintf("/fizzbuzz?int1=3&int2=5&limit=%d&str1=a&str2=b", cfg.MaxLimit), 200, ""},

		// strings
		{"empty_str1", "/fizzbuzz?int1=3&int2=5&limit=10&str1=&str2=b", 400, "non-empty"},
		{"empty_str2", "/fizzbuzz?int1=3&int2=5&limit=10&str1=a&str2=", 400, "non-empty"},
		{"control_chars_str1", "/fizzbuzz?int1=3&int2=5&limit=10&str1=a%0A&str2=b", 400, "control"},
		{"control_chars_str2", "/fizzbuzz?int1=3&int2=5&limit=10&str1=a&str2=b%0D", 400, "control"},
		{"str1_too_long", "/fizzbuzz?int1=3&int2=5&limit=10&str1=" + url.QueryEscape(tooLong) + "&str2=b", 400, "too long"},
		{"str1_max_ok", "/fizzbuzz?int1=3&int2=5&limit=10&str1=" + url.QueryEscape(okLen) + "&str2=b", 200, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := doGET(t, router, tt.path)
			if rr.Code != tt.wantStatus {
				t.Fatalf("status=%d want=%d body=%s", rr.Code, tt.wantStatus, rr.Body.String())
			}

			if tt.wantStatus != 200 {
				var er errorResponse
				if err := json.NewDecoder(rr.Body).Decode(&er); err != nil {
					t.Fatalf("expected json error body, decode err=%v body=%q", err, rr.Body.String())
				}
				if er.Error == "" {
					t.Fatalf("expected non-empty error")
				}
				if tt.wantErrContains != "" && !strings.Contains(strings.ToLower(er.Error), strings.ToLower(tt.wantErrContains)) {
					t.Fatalf("error=%q does not contain %q", er.Error, tt.wantErrContains)
				}
			}
		})
	}
}

func TestStats_NoData_Returns204(t *testing.T) {
	_, router := newTestRouter(t)
	rr := doGET(t, router, "/stats")
	if rr.Code != http.StatusNoContent {
		t.Fatalf("status=%d want=204 body=%s", rr.Code, rr.Body.String())
	}
}

func TestStats_Disabled_Returns204(t *testing.T) {
	cfg, err := config.New()
	if err != nil {
		t.Fatal(err)
	}
	cfg.RedisEnabled = false

	h := NewHandler(cfg, testLogger(), nil, nil)
	router := NewRouter(cfg, testLogger(), h)

	rr := doGET(t, router, "/stats")
	if rr.Code != http.StatusNoContent {
		t.Fatalf("status=%d want=204 body=%s", rr.Code, rr.Body.String())
	}
}

func TestStats_AfterFizzBuzz_ReturnsTop(t *testing.T) {
	_, router := newTestRouter(t)

	// hit twice
	for i := 0; i < 2; i++ {
		rr := doGET(t, router, "/fizzbuzz?int1=3&int2=5&limit=16&str1=fizz&str2=buzz")
		if rr.Code != http.StatusOK {
			t.Fatalf("fizzbuzz hit %d status=%d body=%s", i+1, rr.Code, rr.Body.String())
		}
	}

	rr := doGET(t, router, "/stats")
	if rr.Code != http.StatusOK {
		t.Fatalf("status=%d want=%d body=%s", rr.Code, http.StatusOK, rr.Body.String())
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
	if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v body=%s", err, rr.Body.String())
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

func TestSecurityHeaders_AreSet(t *testing.T) {
	_, router := newTestRouter(t)

	rr := doGET(t, router, "/healthz")
	if rr.Code != http.StatusOK {
		t.Fatalf("status=%d want=200 body=%s", rr.Code, rr.Body.String())
	}

	h := rr.Result().Header
	if h.Get("X-Content-Type-Options") != "nosniff" {
		t.Fatalf("X-Content-Type-Options=%q want nosniff", h.Get("X-Content-Type-Options"))
	}
	if h.Get("X-Frame-Options") != "DENY" {
		t.Fatalf("X-Frame-Options=%q want DENY", h.Get("X-Frame-Options"))
	}
	if h.Get("Referrer-Policy") != "no-referrer" {
		t.Fatalf("Referrer-Policy=%q want no-referrer", h.Get("Referrer-Policy"))
	}
	if h.Get("Content-Security-Policy") != "default-src 'none'" {
		t.Fatalf("Content-Security-Policy=%q want default-src 'none'", h.Get("Content-Security-Policy"))
	}
}

func TestReadyz_RedisDisabled_OK(t *testing.T) {
	cfg, err := config.New()
	if err != nil {
		t.Fatal(err)
	}
	cfg.RedisEnabled = false

	h := NewHandler(cfg, testLogger(), nil, nil)
	router := NewRouter(cfg, testLogger(), h)

	rr := doGET(t, router, "/readyz")
	if rr.Code != http.StatusOK {
		t.Fatalf("status=%d want=200 body=%s", rr.Code, rr.Body.String())
	}
	if strings.TrimSpace(rr.Body.String()) != "ready" {
		t.Fatalf("body=%q want %q", rr.Body.String(), "ready")
	}
}
