FROM golang:1.22.8 AS builder
ENV CGO_ENABLED=0 GOFLAGS=-mod=vendor
ARG APP_VERSION
ARG CONFIG_PATH
WORKDIR /build
COPY . .
RUN go build -ldflags "-X main.VERSION=$APP_VERSION -X main.CONFIG_PATH=$CONFIG_PATH" -o auth /build/cmd/auth/

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /build/auth /app/auth

ENTRYPOINT ["sh", "-c", "/app/auth --grpc-port=$GRPC_PORT --http-port=$HTTP_PORT"]