FROM golang:1.21-alpine AS builder
LABEL authors="Nite07"

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init
RUN CGO_ENABLED=0 GOOS=linux go build -o pcgamedb .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/pcgamedb /app/pcgamedb

ENTRYPOINT ["/app/pcgamedb", "server"]