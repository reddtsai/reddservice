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

.PHONY: run-swag
run-swag:
	@go run ./api

.PHONY: build-auth
build-auth:
	@go build -ldflags "-X main.VERSION=0.0.1 -X main.CONFIG_PATH=$(CONFIG_PATH)" -o bin/auth ./cmd/auth/
#	go build -gcflags="-m" -o bin/auth ./cmd/auth/

.PHONY: build-gateway
build-auth:
	@go build -ldflags "-X main.VERSION=0.0.1 -X main.CONFIG_PATH=$(CONFIG_PATH)" -o bin/gateway ./cmd/gateway/

.PHONY: test
test:
	@go test -tags unittest ./...

.PHONY: test-cover
test-cover:
	@go test -cover -tags unittest ./...

.PHONY: gen-swagger
gen-swagger:
	swag init -g gateway.go -d ./cmd/gateway -t auth -o api/gateway --pd

.PHONY: gen-mock
gen-mock:
	go generate ./...

.PHONY: build-docker-image
build-docker-image:
	docker build --build-arg APP_VERSION=dev --build-arg CONFIG_PATH=$(CONFIG_PATH) -f deployments/docker/auth/Dockerfile -t reddservice-auth:dev .
#	docker tag reddservice-auth:dev reddtsai/reddservice-auth:dev
	docker push reddtsai/reddservice-auth:dev
	docker build --build-arg APP_VERSION=dev --build-arg CONFIG_PATH=$(CONFIG_PATH) -f deployments/docker/gateway/Dockerfile -t reddservice-gateway:dev .
#	docker tag reddservice-gateway:dev reddtsai/reddservice-gateway:dev
	docker push reddtsai/reddservice-gateway:dev

.PHONY: docker-compose-up
docker-compose-up: build-docker-image
	rm -rf deployments/docker/data/postgres
	docker-compose -p reddservice -f deployments/docker/docker-compose.yml up -d
