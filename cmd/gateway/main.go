//go:build !unittest

package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/reddtsai/reddservice/internal/global"
)

// here value is set by ldflags
var (
	VERSION     = "dev"
	CONFIG_PATH = "conf.d"
)

type GatewaySrv struct {
	wg             sync.WaitGroup
	httpSrv        *http.Server
	httpPort       int
	grpcClientConn IGrpcClientConn
}

var (
	_gatewaySrv *GatewaySrv
)

func init() {
	_gatewaySrv = new(GatewaySrv)
	flag.IntVar(&_gatewaySrv.httpPort, "http-port", 8081, "gateway server port")
	global.Startup(CONFIG_PATH)
}

func main() {
	flag.Parse()
	defer global.Logger.Sync()

	global.Logger.Debug("starting gateway server")
	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, syscall.SIGINT, syscall.SIGTERM)
	_gatewaySrv.connGrpc()
	_gatewaySrv.httpServer()
	global.Logger.Debug("gateway server started")

	<-shutdownCh
	global.Logger.Debug("shutting down gateway server")
	_gatewaySrv.shutdown()
	global.Logger.Debug("gateway server stopped")
}

func (srv *GatewaySrv) shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := srv.httpSrv.Shutdown(ctx)
	if err != nil {
		global.Logger.Error("gateway server shutdown failed", zap.Error(err))
	}
	srv.wg.Wait()
}

func (srv *GatewaySrv) httpServer() {
	g := NewGateway()
	h := NewHandler(srv.grpcClientConn)
	g.register(h)
	srv.httpSrv = &http.Server{
		Addr:    fmt.Sprintf(":%d", srv.httpPort),
		Handler: g.Handler,
	}
	srv.wg.Add(1)
	go func() {
		defer srv.wg.Done()
		err := srv.httpSrv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			global.Logger.Fatal("gateway failed to serve", zap.Error(err))
		}
	}()
}

func (srv *GatewaySrv) connGrpc() {
	conn, err := NewClientConn()
	if err != nil {
		global.Logger.Fatal("conn grpc client failed", zap.Error(err))
	}
	srv.grpcClientConn = conn
}
