package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("🔔 Ding Dong")
	json.NewEncoder(w).Encode("Electronic-Commerce server is running...")
}
