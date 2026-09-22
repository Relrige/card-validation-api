FROM golang:1.26-alpine AS builder

WORKDIR /app
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /api-server ./cmd/api

FROM alpine:3.19

WORKDIR /root/
COPY --from=builder /api-server .

EXPOSE 8080

CMD ["./api-server"]