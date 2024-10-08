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

	ctx context.Context
}

func NewGrpcHandler(ctx context.Context, grpcServer *grpc.Server) {
	pb.RegisterAuthServiceServer(grpcServer, &grpcHandler{
		ctx: ctx,
	})
}

func (h *grpcHandler) SignUp(ctx context.Context, req *pb.SignUpRequest) (*pb.SignUpResponse, error) {
	// TODO: implement SignUp
	return nil, status.Error(codes.Unimplemented, "not implemented")
}
