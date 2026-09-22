package validator

import (
	"strconv"
	"strings"
	"time"

	"github.com/Relrige/card-validator-api/internal/domain"
)

const (
	minCardNumberLength = 12
	maxCardNumberLength = 19
	maxYearsInFuture    = 20
)

func ValidateCard(card domain.CardRequest, now time.Time) *domain.ValidationResponse {
	if card.CardNumber == "" || card.ExpirationMonth == "" || card.ExpirationYear == "" {
		return errorResponse(domain.ErrMissingFields, "Missing required fields")
	}

	cleanCardNumber := strings.ReplaceAll(card.CardNumber, " ", "")
	if len(cleanCardNumber) < minCardNumberLength || len(cleanCardNumber) > maxCardNumberLength || !isValidLuhn(cleanCardNumber) {
		return errorResponse(domain.ErrInvalidCardNumber, "Invalid card number format")
	}

	month, err := strconv.Atoi(card.ExpirationMonth)
	if err != nil || month < 1 || month > 12 {
		return errorResponse(domain.ErrInvalidMonth, "Invalid expiration month")
	}

	if len(card.ExpirationYear) != 4 {
		return errorResponse(domain.ErrInvalidYearLength, "Invalid expiration year length")
	}

	currentYear := now.Year()
	currentMonth := int(now.Month())

	year, err := strconv.Atoi(card.ExpirationYear)
	if err != nil || year < currentYear || (year == currentYear && month < currentMonth) {
		return errorResponse(domain.ErrCardExpired, "Card has expired")
	}

	if year > currentYear+maxYearsInFuture {
		return errorResponse(domain.ErrYearTooFarInFuture, "Expiration year is too far in the future")
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
