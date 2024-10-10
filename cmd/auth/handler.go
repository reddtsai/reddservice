package main

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/reddtsai/reddservice/api/proto"
	"github.com/reddtsai/reddservice/internal/auth"
)

type AuthHandler struct {
	pb.UnimplementedAuthServiceServer

	authSvc auth.IAuthService
}

func NewAuthHandler(authSvc auth.IAuthService) *AuthHandler {
	return &AuthHandler{
		authSvc: authSvc,
	}
}

func (h *AuthHandler) SignUp(ctx context.Context, req *pb.SignUpRequest) (*pb.SignUpResponse, error) {
	_, err := h.authSvc.CreateUser(auth.ServiceInput[auth.CreateUser]{
		Data: auth.CreateUser{
			Name:  req.Username,
			Email: req.Email,
		},
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.SignUpResponse{}, nil
}
