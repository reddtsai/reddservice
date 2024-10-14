//go:generate mockgen -source=common.go -destination=mock/mock_common.go -package=mock

package main

import (
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/reddtsai/reddservice/api/proto"
	"github.com/reddtsai/reddservice/internal/global"
)

type IGrpcClientConn interface {
	GetAuthClient() pb.AuthServiceClient
}
type GrpcClientConn struct {
	mu      sync.Mutex
	connMap map[string]grpc.ClientConnInterface
}

func NewClientConn() (*GrpcClientConn, error) {
	c := &GrpcClientConn{}
	err := c.clientConn()
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c *GrpcClientConn) clientConn() error {
	cfg := global.GetGrpcClientOptions("auth")
	conn, err := grpc.NewClient(cfg.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	c.mu.Lock()
	c.connMap["auth"] = conn
	c.mu.Unlock()
	return nil
}

func (c *GrpcClientConn) GetAuthClient() pb.AuthServiceClient {
	return pb.NewAuthServiceClient(c.connMap["auth"])
}
