package traefiktimestampheader

import (
	"context"
	"net/http"
	"net/http/httptest"
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

	// Validate format is t=<milliseconds>
	if !strings.HasPrefix(headerValue, "t=") {
		t.Fatalf("Expected header to start with 't=', got: %s", headerValue)
	}
}
