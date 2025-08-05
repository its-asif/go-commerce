package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// @Summary		Ping
// @Description	"Check if server is running fine"
// @Tags			Ping
// @Accept			json
// @Produce		json
// @Failure		200	{string}	string	"All good"
// @Router			/ [GET]
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("🔔 Ding Ding")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Electronic-Commerce server is running...")
}
