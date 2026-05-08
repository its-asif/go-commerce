package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHomeHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	HomeHandler(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", rr.Code)
	}
	var msg string
	if err := json.NewDecoder(rr.Body).Decode(&msg); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if msg == "" {
		t.Fatalf("expected non-empty message")
	}
}
