package handlers

import (
	"encoding/json"
	"net/http"
)

// HealthHandler returns a basic health status for readiness/liveness probes.
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
