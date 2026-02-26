package main

import (
	"context"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func runHealthcheck(addr string) int {
	port, ok := extractPort(addr)
	if !ok {
		return 1
	}

	u := &url.URL{
		Scheme: "http",
		Host:   net.JoinHostPort("127.0.0.1", port),
		Path:   "/healthz",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return 1
	}

	client := &http.Client{}
	// #nosec G704 -- healthcheck makes a local call to 127.0.0.1 only (host forced, port validated)
	resp, err := client.Do(req)
	if err != nil {
		return 1
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return 1
	}
	return 0
}

func extractPort(addr string) (string, bool) {
	if strings.HasPrefix(addr, ":") {
		p := strings.TrimPrefix(addr, ":")
		return validatePort(p)
	}

	_, p, err := net.SplitHostPort(addr)
	if err == nil {
		return validatePort(p)
	}

	return "", false
}

func validatePort(p string) (string, bool) {
	n, err := strconv.Atoi(p)
	if err != nil || n < 1 || n > 65535 {
		return "", false
	}
	return strconv.Itoa(n), true
}
