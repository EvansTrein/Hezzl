FROM golang:1.24.1-alpine AS builder

WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod tidy

COPY . .

RUN go build -o events ./cmd/events/events.go

FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/events .
COPY --from=builder /app/example.env .

CMD ["sh", "-c", "sleep 3 && ./events -config ./example.env"]