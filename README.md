# Card Validation API

A small HTTP service that validates a bank card's number and expiration date
and returns a structured `valid: true/false` verdict.

## How it works

The service validates three fields on every request, in order, and stops at
the first failure:

1. **Presence** — `card_number`, `expiration_month`, `expiration_year` must
   all be non-empty.
2. **Card number** — digits only (spaces/dashes are stripped first),
   12–19 digits long, and must pass the [Luhn checksum](https://en.wikipedia.org/wiki/Luhn_algorithm).
3. **Expiration month** — an integer from `1` to `12`.
4. **Expiration year** — exactly 4 digits, not in the past (relative to the
   card's month), and not more than 20 years in the future.


# API

### POST /validate

**Request**

```json
{
  "card_number": "4111111111111111",
  "expiration_month": "12",
  "expiration_year": "2028"
}
```

**Response — valid card** (`200 OK`)

```json
{
  "valid": true
}
```

**Response — invalid card** (`200 OK`)

```json
{
  "valid": false,
  "error": {
    "code": "005",
    "message": "Card has expired"
  }
}
```

| Status | When |
|---|---|
| `200 OK` | Request parsed successfully; body says whether the card is valid |
| `400 Bad Request` | Request body isn't valid JSON |
| `405 Method Not Allowed` | Any method other than `POST` |

### Error codes

| Code | Meaning |
|---|---|
| `001` | One or more required fields are missing |
| `002` | Card number isn't 12–19 digits or fails the Luhn check |
| `003` | Expiration month isn't an integer between 1 and 12 |
| `004` | Expiration year isn't exactly 4 digits |
| `005` | Card is expired |
| `006` | Expiration year is more than 20 years in the future |

## Running

### Docker

```bash
docker build -t card-validator-api .
docker run -p 8080:8080 card-validator-api
```


## Project structure

```
cmd/api/            entry point, HTTP server wiring
internal/domain/     request/response types, error codes
internal/handler/    HTTP layer: parsing, status codes, JSON encoding
internal/validator/  pure domain logic: card number, month, year, expiry
```