protoc:
	protoc -I . -I scrapper/api/third_party --go_out=. --go-grpc_out=. --grpc-gateway_out=. scrapper/api/proto/scrapper.proto

protoc-swagger:
	protoc -I . -I scrapper/api/third_party --go_out=. --go-grpc_out=. --grpc-gateway_out=. --openapiv2_out=. scrapper/api/proto/scrapper.proto


protoc-bot:
	protoc -I . -I bot/api/third_party --go_out=. --go-grpc_out=. --grpc-gateway_out=. --openapiv2_out=. bot/api/proto/server.proto

all: bot scrapper

bot:
	go run bot/cmd/bot/main.go

scrapper:
	go run scrapper/cmd/main/main.go

.PHONY: all bot scrapper protoc protoc-swagger


## mockgen -source=pkg/github/client.go -destination=internal/mocks/github_mock.go -package=mocks
