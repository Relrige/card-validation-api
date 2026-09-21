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
