.PHONY: gen-grpc
gen-grpc:
	@protoc -I=common/api/proto --go_out=common/api/proto --go_opt=paths=source_relative --go-grpc_out=common/api/proto --go-grpc_opt=paths=source_relative auth.proto common.proto

.PHONY: go-mod-tidy
go-mod-tidy:
	@for dir in ./auth ./common ./gateway; do \
		cd $$dir; \
		go mod tidy; \
		cd ..; \
	done