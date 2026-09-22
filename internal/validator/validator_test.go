package validator

import (
	"testing"
	"time"

	"github.com/Relrige/card-validator-api/internal/domain"
)

func TestValidateCard(t *testing.T) {
	now := time.Date(2026, time.September, 22, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name      string
		req       domain.CardRequest
		wantValid bool
		wantCode  string
	}{
		{
			name: "valid visa test card",
			req: domain.CardRequest{
				CardNumber:      "4111111111111111",
				ExpirationMonth: "12",
				ExpirationYear:  "2028",
			},
			wantValid: true,
		},
		{
			name: "spaces and dashes are stripped before validation",
			req: domain.CardRequest{
				CardNumber:      "4111-1111 1111-1111",
				ExpirationMonth: "12",
				ExpirationYear:  "2028",
			},
			wantValid: true,
		},
		{
			name: "valid: expires this exact month",
			req: domain.CardRequest{
				CardNumber:      "4111111111111111",
				ExpirationMonth: "09",
				ExpirationYear:  "2026",
			},
			wantValid: true,
		},
		{
			name:      "missing all fields",
			req:       domain.CardRequest{},
			wantValid: false,
			wantCode:  domain.ErrMissingFields,
		},
		{
			name: "missing card number only",
			req: domain.CardRequest{
				ExpirationMonth: "12",
				ExpirationYear:  "2028",
			},
			wantValid: false,
			wantCode:  domain.ErrMissingFields,
		},
		{
			name: "fails luhn check (spec example)",
			req: domain.CardRequest{
				CardNumber:      "1111111111111",
				ExpirationMonth: "10",
				ExpirationYear:  "2028",
			},
			wantValid: false,
			wantCode:  domain.ErrInvalidCardNumber,
		},
		{
			name: "card number too short",
			req: domain.CardRequest{
				CardNumber:      "4111",
				ExpirationMonth: "10",
				ExpirationYear:  "2028",
			},
			wantValid: false,
			wantCode:  domain.ErrInvalidCardNumber,
		},
		{
			name: "card number too long",
			req: domain.CardRequest{
				CardNumber:      "41111111111111111111",
				ExpirationMonth: "10",
				ExpirationYear:  "2028",
			},
			wantValid: false,
			wantCode:  domain.ErrInvalidCardNumber,
		},
		{
			name: "card number contains letters",
			req: domain.CardRequest{
				CardNumber:      "411111111111abcd",
				ExpirationMonth: "10",
				ExpirationYear:  "2028",
			},
			wantValid: false,
			wantCode:  domain.ErrInvalidCardNumber,
		},
		{
			name: "invalid month zero",
			req: domain.CardRequest{
				CardNumber:      "4111111111111111",
				ExpirationMonth: "0",
				ExpirationYear:  "2028",
			},
			wantValid: false,
			wantCode:  domain.ErrInvalidMonth,
		},
		{
			name: "invalid month 13",
			req: domain.CardRequest{
				CardNumber:      "4111111111111111",
				ExpirationMonth: "13",
				ExpirationYear:  "2028",
			},
			wantValid: false,
			wantCode:  domain.ErrInvalidMonth,
		},
		{
			name: "non-numeric month",
			req: domain.CardRequest{
				CardNumber:      "4111111111111111",
				ExpirationMonth: "ab",
				ExpirationYear:  "2028",
			},
			wantValid: false,
			wantCode:  domain.ErrInvalidMonth,
		},
		{
			name: "year not 4 digits",
			req: domain.CardRequest{
				CardNumber:      "4111111111111111",
				ExpirationMonth: "12",
				ExpirationYear:  "28",
			},
			wantValid: false,
			wantCode:  domain.ErrInvalidYearLength,
		},
		{
			name: "expired: past year (spec example)",
			req: domain.CardRequest{
				CardNumber:      "4111111111111111",
				ExpirationMonth: "01",
				ExpirationYear:  "2021",
			},
			wantValid: false,
			wantCode:  domain.ErrCardExpired,
		},
		{
			name: "expired: current year, earlier month",
			req: domain.CardRequest{
				CardNumber:      "4111111111111111",
				ExpirationMonth: "08",
				ExpirationYear:  "2026",
			},
			wantValid: false,
			wantCode:  domain.ErrCardExpired,
		},
		{
			name: "year too far in the future",
			req: domain.CardRequest{
				CardNumber:      "4111111111111111",
				ExpirationMonth: "12",
				ExpirationYear:  "2099",
			},
			wantValid: false,
			wantCode:  domain.ErrYearTooFarInFuture,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateCard(tt.req, now)

			if got.Valid != tt.wantValid {
				t.Fatalf("Valid = %v, want %v (error: %+v)", got.Valid, tt.wantValid, got.Error)
			}
			if !tt.wantValid {
				if got.Error == nil {
					t.Fatalf("expected an Error, got nil")
				}
				if got.Error.Code != tt.wantCode {
					t.Fatalf("Error.Code = %q, want %q", got.Error.Code, tt.wantCode)
				}
			}
		})
	}
}

func TestPassesLuhnCheck(t *testing.T) {
	cases := []struct {
		number string
		want   bool
	}{
		{"4111111111111111", true},
		{"1111111111111", false},
		{"79927398713", true},
	}

	for _, c := range cases {
		if got := isValidLuhn(c.number); got != c.want {
			t.Errorf("isValidLuhn(%q) = %v, want %v", c.number, got, c.want)
		}
	}
}
