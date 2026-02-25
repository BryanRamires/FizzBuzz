package main

import (
	"context"
	"net"
	"net/http"
	"strings"
	"time"
)

func runHealthcheck(addr string) int {
	target := addr
	if strings.HasPrefix(target, ":") {
		target = "127.0.0.1" + target
	} else if host, port, err := net.SplitHostPort(target); err == nil {
		_ = host
		target = net.JoinHostPort("127.0.0.1", port)
	} else {
		return 1
	}

	url := "http://" + target + "/healthz"

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 1
	}

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return 1
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 1
	}
	return 0
}
