package main

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/reddtsai/reddservice/common/api/proto"
)

type grpcHandler struct {
	pb.UnimplementedAuthServiceServer
}

func NewGrpcHandler(grpcServer *grpc.Server) {
	pb.RegisterAuthServiceServer(grpcServer, &grpcHandler{})
}

func (h *grpcHandler) SignUp(ctx context.Context, req *pb.SignUpRequest) (*pb.SignUpResponse, error) {
	// TODO: implement SignUp
	return nil, status.Error(codes.Unimplemented, "not implemented")
}
