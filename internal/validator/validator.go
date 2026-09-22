package validator

import (
	"strconv"
	"strings"
	"time"

	"github.com/Relrige/card-validator-api/internal/domain"
)

func ValidateCard(card domain.CardRequest, now time.Time) *domain.ValidationResponse {
	if card.CardNumber == "" || card.ExpirationMonth == "" || card.ExpirationYear == "" {
		return errorResponse("001", "Missing required fields")
	}

	cleanCardNumber := strings.ReplaceAll(card.CardNumber, " ", "")
	if len(cleanCardNumber) < 12 || len(cleanCardNumber) > 19 || !isValidLuhn(cleanCardNumber) {
		return errorResponse("002", "Invalid card number format")
	}

	month, err := strconv.Atoi(card.ExpirationMonth)
	if err != nil || month < 1 || month > 12 {
		return errorResponse("003", "Invalid expiration month")
	}

	if len(card.ExpirationYear) != 4 {
		return errorResponse("004", "Invalid expiration year length")
	}

	currentYear := now.Year()
	currentMonth := int(now.Month())

	year, err := strconv.Atoi(card.ExpirationYear)
	if err != nil || year < currentYear || (year == currentYear && month < currentMonth) {
		return errorResponse("005", "Card has expired")
	}

	if year > currentYear+20 {
		return errorResponse("006", "Expiration year is too far in the future")
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
