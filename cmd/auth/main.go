package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	grpczap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/reddtsai/reddservice/internal/global"
)

var (
	wg       sync.WaitGroup
	grpcSrv  *grpc.Server
	grpcPort = 50051
	httpSrv  *http.Server
	httpPort = 8081
)

func init() {
	flag.IntVar(&grpcPort, "grpc-port", 50051, "auth server port")
	flag.IntVar(&httpPort, "http-port", 8081, "metrics server port")
}

func main() {
	flag.Parse()
	defer global.Logger.Sync()
	global.Logger.Info("starting auth server")
	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())
	srvMetrics := grpcprom.NewServerMetrics()
	grpczap.ReplaceGrpcLoggerV2(global.Logger)
	grpcSrv = global.GrpcServer(
		srvMetrics.UnaryServerInterceptor(),
		grpczap.UnaryServerInterceptor(global.Logger),
	)
	NewGrpcHandler(ctx, grpcSrv)
	httpSrv = global.MetricsServer(fmt.Sprintf(":%d", httpPort), srvMetrics)
	wg.Add(1)
	go grpcListen()
	wg.Add(1)
	go httpListen()
	global.Logger.Info("auth server started")

	<-shutdownCh
	cancel()
	shutdown()
	wg.Wait()
	global.Logger.Info("auth server stopped")
}

func shutdown() {
	global.Logger.Info("shutting down auth server")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := httpSrv.Shutdown(ctx)
	if err != nil {
		global.Logger.Error("metrics server shutdown failed", zap.Error(err))
	}
	grpcSrv.GracefulStop()
}

func grpcListen() {
	defer wg.Done()
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		global.Logger.Fatal("grpc failed to listen", zap.Error(err))
	}
	defer listen.Close()
	err = grpcSrv.Serve(listen)
	if err != nil {
		global.Logger.Fatal("grpc failed to serve", zap.Error(err))
	}
}

func httpListen() {
	defer wg.Done()
	err := httpSrv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		global.Logger.Fatal("metrics failed to serve", zap.Error(err))
	}
}
