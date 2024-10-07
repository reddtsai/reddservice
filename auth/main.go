package main

import (
	"flag"
	"fmt"
	"net"

	"google.golang.org/grpc"

	"github.com/reddtsai/reddservice/common/logger"
)

var addr = 50051

func init() {
	flag.IntVar(&addr, "addr", 50051, "auth server port")
}

func main() {
	logger.New()
	defer logger.Sync()

	grpcServer := grpc.NewServer()
	NewGrpcHandler(grpcServer)
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", addr))
	if err != nil {
		logger.Fatal("grpc failed to listen", err)
	}
	defer listen.Close()
	err = grpcServer.Serve(listen)
	if err != nil {
		logger.Fatal("grpc failed to serve", err)
	}
}
