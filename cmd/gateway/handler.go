package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	pb "github.com/reddtsai/reddservice/api/proto"
	"github.com/reddtsai/reddservice/internal/global"
)

type Handler struct {
	grpcClientConn IGrpcClientConn
}

func NewHandler(conn IGrpcClientConn) *Handler {
	return &Handler{
		grpcClientConn: conn,
	}
}

// @Summary 註冊
// @Description 註冊用戶
// @Tags auth
// @Accept json
// @Produce json
// @Param Request body SignUpRequest true "raw"
// @Success 200 {object} Response{result=SignUpResponse} "ok"
// @Failure 400 {object} Response "bad request"
// @Failure 401 {object} Response "unauthorized"
// @Failure 403 {object} Response "forbidden"
// @Failure 409 {object} Response "conflict"
// @Failure 500 {object} Response "server error"
// @Router /v1/signup [post]
func (h *Handler) SignUp(c *gin.Context) {
	req := SignUpRequest{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authClient := h.grpcClientConn.GetAuthClient()
	_, err = authClient.SignUp(c.Request.Context(), &pb.SignUpRequest{
		Account: req.Account,
		Email:   req.Email,
	})
	if err != nil {
		global.Logger.Error("sign up failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "sign up failed"})
		return
	}
	resp := Response{
		Result: SignUpResponse{},
	}
	c.JSON(http.StatusOK, resp)
}
