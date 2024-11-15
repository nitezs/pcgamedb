FROM golang:1.21-alpine AS builder
LABEL authors="Nite07"

RUN apk add --no-cache git

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init
ARG version=dev
RUN if [ "$version" = "dev" ]; then \
    version=$(git describe --tags --always); \
    fi && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w -X github.com/nitezs/pcgamedb/constant.Version=${version}" -o pcgamedb .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/pcgamedb /app/pcgamedb

ENTRYPOINT ["/app/pcgamedb", "server"]