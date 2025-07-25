package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("🔔 Ding Ding")
	json.NewEncoder(w).Encode("Electronic-Commerce server is running...")
}
