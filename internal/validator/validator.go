package validator

import (
	"github.com/Relrige/card-validator-api/internal/domain"
)

func ValidateCard(card domain.CardRequest) *domain.ValidationResponse {
	if card.CardNumber == "" || card.ExpirationMonth == "" || card.ExpirationYear == "" {
		return errorResponse("001", "Missing required fields")
	}

	if len(card.CardNumber) < 12 || len(card.CardNumber) > 19 || !isValidLuhn(card.CardNumber) {
		return errorResponse("002", "Invalid card number format")
	}

	return &domain.ValidationResponse{
		Valid: true,
	}
}

func errorResponse(code, message string) *domain.ValidationResponse {
	return &domain.ValidationResponse{
		Valid: false,
		Error: &domain.ErrorDetail{
			Code:    code,
			Message: message,
		},
	}
}

func isValidLuhn(cardNumber string) bool {
	sum := 0
	isSecond := false
	for i := len(cardNumber) - 1; i >= 0; i-- {
		if cardNumber[i] < '0' || cardNumber[i] > '9' {
			return false
		}
		digit := int(cardNumber[i] - '0')
		if isSecond {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		isSecond = !isSecond
		sum += digit
	}
	return sum%10 == 0
}
