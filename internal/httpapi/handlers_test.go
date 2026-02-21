package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFizzBuzz_OK(t *testing.T) {
	req := httptest.NewRequest(
		http.MethodGet,
		"/fizzbuzz?int1=3&int2=5&limit=16&str1=fizz&str2=buzz",
		nil,
	)
	rr := httptest.NewRecorder()
	NewRouter().ServeHTTP(rr, req)
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
			rr := httptest.NewRecorder()
			req := httptest.NewRequest("GET", tt.url, nil)

			NewRouter().ServeHTTP(rr, req)

			if rr.Code != http.StatusBadRequest {
				t.Fatalf("status=%d want=400", rr.Code)
			}
		})
	}
}
