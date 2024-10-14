package main

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/reddtsai/reddservice/api/proto"
	"github.com/reddtsai/reddservice/internal/auth"
	"github.com/reddtsai/reddservice/internal/auth/mock"
)

func TestNewAuthHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := mock.NewMockIAuthService(ctrl)
	handler := NewAuthHandler(mockAuthService)

	assert.NotNil(t, handler)
}

func TestSignUp(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := mock.NewMockIAuthService(ctrl)
	handler := NewAuthHandler(mockAuthService)
	var id int64

	t.Run("successful sign up", func(t *testing.T) {
		req := &pb.SignUpRequest{
			Account: "test_account",
			Email:   "test_account@example.com",
		}

		mockAuthService.EXPECT().CreateUser(auth.CreateUserInput{
			Account: req.Account,
			Email:   req.Email,
		}).Return(id, nil)

		resp, err := handler.SignUp(context.Background(), req)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, int32(0), resp.Status.Status)
	})

	t.Run("failed sign up due to internal error", func(t *testing.T) {
		req := &pb.SignUpRequest{
			Account: "test_account",
			Email:   "test_account@example.com",
		}

		mockAuthService.EXPECT().CreateUser(auth.CreateUserInput{
			Account: req.Account,
			Email:   req.Email,
		}).Return(id, errors.New("internal error"))

		resp, err := handler.SignUp(context.Background(), req)
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, codes.Internal, status.Code(err))
	})
}
