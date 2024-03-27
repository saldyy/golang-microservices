package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Repeat(character string) string {
	var repeated string
	for i := 0; i < 5; i++ {
		repeated += character
	}
	return repeated
}

func BenchmarkHealthCheckHandler(b *testing.B) {
	// Create a new HTTP request to the handler
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		b.Fatalf("Failed to create request: %v", err)
	}

	// Create a new ResponseRecorder (which implements http.ResponseWriter) to record the response
	rr := httptest.NewRecorder()

	// Run the benchmark b.N times
	for i := 0; i < b.N; i++ {
		// Call the handler function directly with the ResponseRecorder and Request
		HealthCheckHandler(rr, req)
	}
}
