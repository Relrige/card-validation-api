package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Relrige/card-validator-api/internal/domain"
	"github.com/Relrige/card-validator-api/internal/validator"
)

func ValidateCardHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, domain.ValidationResponse{
			Valid: false,
			Error: &domain.ErrorDetail{Code: "405", Message: "Method not allowed"},
		})
		return
	}

	var req domain.CardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, domain.ValidationResponse{
			Valid: false,
			Error: &domain.ErrorDetail{Code: "400", Message: "Invalid JSON format"},
		})
		return
	}

	response := validator.ValidateCard(req, time.Now())

	writeJSON(w, http.StatusOK, response)
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}
