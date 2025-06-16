FROM golang:1.24.1-alpine AS builder

WORKDIR /app
COPY go.mod ./
COPY go.sum ./
COPY migrations ./
RUN go mod download

COPY . .

RUN go build -o main ./cmd/main.go
RUN go build -o migrator ./cmd/migrator/migrator.go

FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/migrator .
COPY --from=builder /app/example.env .
COPY --from=builder /app/entrypoint.sh .
RUN chmod +x ./entrypoint.sh
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

CMD ["./entrypoint.sh"]