go install github.com/swaggo/swag/cmd/swag@latest
swag init
CGO_ENABLED=0
go build -o pcgamedb -ldflags "-s -w -X github.com/nitezs/pcgamedb/constant.Version=$(git describe --tags --always)" .
