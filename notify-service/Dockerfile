FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o notify-service cmd/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/notify-service .
COPY --from=builder /app/.env.example .env

EXPOSE 8081

CMD ["./notify-service"]
