GO_MODULE=github.com/reddtsai/reddservice
CONFIG_PATH=conf.d

.PHONY: gen-grpc
gen-grpc:
	@protoc -I=api/proto --go_out=api/proto --go_opt=paths=source_relative --go-grpc_out=api/proto --go-grpc_opt=paths=source_relative auth.proto common.proto

.PHONY: run-auth
run-auth:
	@go run ./cmd/auth/ --grpc-port=50051 --http-port=8080

.PHONY: run-gateway
run-gateway:
	@go run ./cmd/gateway/ --http-port=80

.PHONY: build-auth
build-auth:
	@go build -ldflags "-X main.VERSION=0.0.1 -X $(GO_MODULE)/internal/global.CONFIG_PATH=$(CONFIG_PATH)" -o bin/auth ./cmd/auth/
#	go build -gcflags="-m" -o bin/auth ./cmd/auth/

.PHONY: docker-compose-up
docker-compose-up:
	rm -rf deployments/docker/data/postgres
	docker-compose -p reddservice -f deployments/docker/docker-compose.yml up -d