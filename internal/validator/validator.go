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

	if resp := validateCardNumber(card.CardNumber); resp != nil {
		return resp
	}

	month, resp := validateMonth(card.ExpirationMonth)
	if resp != nil {
		return resp
	}

	year, resp := validateYear(card.ExpirationYear)
	if resp != nil {
		return resp
	}

	if resp := validateExpiration(month, year, now); resp != nil {
		return resp
	}

	return &domain.ValidationResponse{
		Valid: true,
	}
}

func validateCardNumber(raw string) *domain.ValidationResponse {
	clean := cleanCardNumber(raw)

	if len(clean) < minCardNumberLength || len(clean) > maxCardNumberLength {
		return errorResponse(domain.ErrInvalidCardNumber, "Invalid card number format")
	}
	if !isNumeric(clean) {
		return errorResponse(domain.ErrInvalidCardNumber, "Invalid card number format")
	}
	if !isValidLuhn(clean) {
		return errorResponse(domain.ErrInvalidCardNumber, "Invalid card number format")
	}
	return nil
}

func cleanCardNumber(raw string) string {
	replacer := strings.NewReplacer(" ", "", "-", "")
	return replacer.Replace(raw)
}

func isNumeric(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

func isValidLuhn(cardNumber string) bool {
	sum := 0
	isSecond := false
	for i := len(cardNumber) - 1; i >= 0; i-- {
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

func validateMonth(raw string) (int, *domain.ValidationResponse) {
	month, err := strconv.Atoi(raw)
	if err != nil || month < 1 || month > 12 {
		return 0, errorResponse(domain.ErrInvalidMonth, "Invalid expiration month")
	}
	return month, nil
}

func validateYear(raw string) (int, *domain.ValidationResponse) {
	if len(raw) != 4 {
		return 0, errorResponse(domain.ErrInvalidYearLength, "Invalid expiration year length")
	}
	year, err := strconv.Atoi(raw)
	if err != nil {
		return 0, errorResponse(domain.ErrInvalidYearLength, "Invalid expiration year length")
	}
	return year, nil
}

func validateExpiration(month, year int, now time.Time) *domain.ValidationResponse {
	currentYear := now.Year()
	currentMonth := int(now.Month())

	if year < currentYear || (year == currentYear && month < currentMonth) {
		return errorResponse(domain.ErrCardExpired, "Card has expired")
	}
	if year > currentYear+maxYearsInFuture {
		return errorResponse(domain.ErrYearTooFarInFuture, "Expiration year is too far in the future")
	}
	return nil
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
