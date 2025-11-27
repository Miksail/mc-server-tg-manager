FROM golang:1.25 AS builder
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o server ./cmd/mc-server-tg-manager/main.go

# Финальный минимальный образ
FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=builder /app/server .
CMD ["/app/server"]