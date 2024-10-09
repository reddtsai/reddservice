package main

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/reddtsai/reddservice/api/proto"
)

type grpcHandler struct {
	pb.UnimplementedAuthServiceServer
}

func NewGrpcHandler(ctx context.Context, grpcServer *grpc.Server) {
	pb.RegisterAuthServiceServer(grpcServer, &grpcHandler{})
}

func (h *grpcHandler) SignUp(ctx context.Context, req *pb.SignUpRequest) (*pb.SignUpResponse, error) {
	// TODO: implement SignUp
	err := status.Error(codes.Unimplemented, "not implemented")
	return nil, err
}
