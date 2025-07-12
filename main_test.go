package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestHandler_WithoutBackend(t *testing.T) {
	os.Setenv("DIGEST_USER", "")
	os.Setenv("DIGEST_PASS", "dummy")
	os.Setenv("BACKEND_URL", "http://127.0.0.1:65535")

	req := httptest.NewRequest("GET", "http://localhost/test", nil)
	rw := httptest.NewRecorder()

	handler(rw, req)

	resp := rw.Result()
	if resp.StatusCode != http.StatusBadGateway {
		t.Errorf("Expected status 502, got %d", resp.StatusCode)
	}
}

func TestHandler_BadRequest(t *testing.T) {
	os.Setenv("DIGEST_USER", "")
	os.Setenv("DIGEST_PASS", "dummy")
	os.Setenv("BACKEND_URL", "http://127.0.0.1:65535")

	req := httptest.NewRequest("GET", "/", nil)
	rw := httptest.NewRecorder()

	handler(rw, req)

	resp := rw.Result()
	if resp.StatusCode != http.StatusBadGateway {
		t.Errorf("Expected status 502, got %d", resp.StatusCode)
	}
}

