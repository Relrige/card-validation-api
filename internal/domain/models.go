package domain

type CardRequest struct {
	CardNumber      string `json:"card_number"`
	ExpirationMonth string `json:"expiration_month"`
	ExpirationYear  string `json:"expiration_year"`
}

type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ValidationResponse struct {
	Valid bool         `json:"valid"`
	Error *ErrorDetail `json:"error,omitempty"`
}

const (
	ErrMissingFields      = "001"
	ErrInvalidCardNumber  = "002"
	ErrInvalidMonth       = "003"
	ErrInvalidYearLength  = "004"
	ErrCardExpired        = "005"
	ErrYearTooFarInFuture = "006"
)
