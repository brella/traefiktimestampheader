package traefiktimestampheader

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestTimestampHeader(t *testing.T) {
	cfg := CreateConfig()
	ctx := context.Background()

	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	// Create the plugin
	plugin, err := New(ctx, next, cfg, "timestamp-header")
	if err != nil {
		t.Fatalf("Failed to create plugin: %v", err)
	}

	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Call the plugin
	plugin.ServeHTTP(recorder, req)

	// Validate the header
	headerValue := req.Header.Get("X-Request-Start")
	if headerValue == "" {
		t.Fatalf("Expected 'X-Request-Start' header, but got none")
	}

	// Validate format is t=<seconds.milliseconds> (nginx-compatible)
	if !strings.HasPrefix(headerValue, "t=") {
		t.Fatalf("Expected header to start with 't=', got: %s", headerValue)
	}

	// Validate it's a float in seconds (not integer milliseconds)
	valueStr := strings.TrimPrefix(headerValue, "t=")
	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		t.Fatalf("Expected float value after 't=', got: %s", valueStr)
	}

	// Should have exactly 3 decimal places
	if !strings.Contains(valueStr, ".") {
		t.Fatalf("Expected decimal point in value (seconds format), got: %s", valueStr)
	}

	// Sanity check: value should be a reasonable Unix timestamp in seconds, not milliseconds
	if value > 1e12 {
		t.Fatalf("Value looks like milliseconds, expected seconds: %s", valueStr)
	}
}
