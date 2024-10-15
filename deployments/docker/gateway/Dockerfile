FROM golang:1.22.8 AS builder
ENV CGO_ENABLED=0 GOFLAGS=-mod=vendor
ARG APP_VERSION
ARG CONFIG_PATH
WORKDIR /build
COPY . .
RUN go build -ldflags "-X main.VERSION=$APP_VERSION -X main.CONFIG_PATH=$CONFIG_PATH" -o gateway /build/cmd/gateway/

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /build/gateway /app/gateway

ENTRYPOINT ["sh", "-c", "/app/gateway --http-port=$HTTP_PORT"]