// Package traefiktimestampheader is a plugin
package traefiktimestampheader

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// Config holds the plugin configuration.
type Config struct{}

// CreateConfig creates and initializes the plugin configuration.
func CreateConfig() *Config {
	return &Config{}
}

// RequestTimestamp is the struct implementing the Traefik plugin interface.
type RequestTimestamp struct {
	next http.Handler
}

// New creates a new RequestTimestamp plugin instance.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	return &RequestTimestamp{
		next: next,
	}, nil
}

// ServeHTTP adds the "X-Request-Start" header and forwards the request.
func (r *RequestTimestamp) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	msec := time.Now().UnixMilli()
	req.Header.Set("X-Request-Start", fmt.Sprintf("t=%d", msec))
	r.next.ServeHTTP(rw, req)
}
