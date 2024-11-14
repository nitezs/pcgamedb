go install github.com/swaggo/swag/cmd/swag@latest
swag init
CGO_ENABLED=0
go build -o gamedb .
