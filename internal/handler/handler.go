package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Relrige/card-validator-api/internal/domain"
	"github.com/Relrige/card-validator-api/internal/validator"
)

func ValidateCardHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(domain.ValidationResponse{
			Valid: false,
			Error: &domain.ErrorDetail{Code: "405", Message: "Method not allowed"},
		})
		return
	}

	var req domain.CardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(domain.ValidationResponse{
			Valid: false,
			Error: &domain.ErrorDetail{Code: "400", Message: "Invalid JSON format"},
		})
		return
	}

	response := validator.ValidateCard(req)

	if !response.Valid {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	json.NewEncoder(w).Encode(response)
}
