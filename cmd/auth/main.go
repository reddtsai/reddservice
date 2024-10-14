//go:build !unittest

package main

import (
	"context"
	"database/sql"
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
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	pb "github.com/reddtsai/reddservice/api/proto"
	"github.com/reddtsai/reddservice/db/rdb"
	"github.com/reddtsai/reddservice/internal/auth"
	"github.com/reddtsai/reddservice/internal/global"
)

// here value is set by ldflags
var (
	VERSION     = "dev"
	CONFIG_PATH = "conf.d"
)

type AuthSrv struct {
	wg       sync.WaitGroup
	grpcSrv  *grpc.Server
	grpcPort int
	httpSrv  *http.Server
	httpPort int
	authDB   *sql.DB
}

var (
	_authSrv *AuthSrv
)

const (
	DB_NAME = "auth"
)

func init() {
	_authSrv = new(AuthSrv)
	flag.IntVar(&_authSrv.grpcPort, "grpc-port", 50051, "auth server port")
	flag.IntVar(&_authSrv.httpPort, "http-port", 8081, "metrics server port")
	global.Startup(CONFIG_PATH)
}

func main() {
	flag.Parse()
	defer global.Logger.Sync()
	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, syscall.SIGINT, syscall.SIGTERM)
	_, cancel := context.WithCancel(context.Background())

	global.Logger.Debug("starting auth server")
	srvMetrics := grpcprom.NewServerMetrics()
	grpczap.ReplaceGrpcLoggerV2(global.Logger)
	_authSrv.initDB(DB_NAME)
	authSvc := auth.NewService(_authSrv.authDB)
	authHandler := NewAuthHandler(authSvc)
	_authSrv.grpcServer(
		authHandler,
		srvMetrics.UnaryServerInterceptor(),
		grpczap.UnaryServerInterceptor(global.Logger),
	)
	_authSrv.metricsServer(srvMetrics)
	global.Logger.Debug("auth server started")

	<-shutdownCh

	global.Logger.Debug("shutting down auth server")
	cancel()
	_authSrv.shutdown()
	global.Logger.Debug("auth server stopped")
}

func (srv *AuthSrv) shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := srv.httpSrv.Shutdown(ctx)
	if err != nil {
		global.Logger.Error("metrics server shutdown failed", zap.Error(err))
	}
	srv.grpcSrv.GracefulStop()
	srv.authDB.Close()

	srv.wg.Wait()
}

func (srv *AuthSrv) grpcServer(authHandler *AuthHandler, interceptor ...grpc.UnaryServerInterceptor) {
	srv.grpcSrv = grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptor...,
		),
	)
	pb.RegisterAuthServiceServer(srv.grpcSrv, authHandler)
	srv.wg.Add(1)
	go func() {
		defer srv.wg.Done()
		listen, err := net.Listen("tcp", fmt.Sprintf(":%d", srv.grpcPort))
		if err != nil {
			global.Logger.Fatal("grpc failed to listen", zap.Error(err))
		}
		defer listen.Close()
		err = srv.grpcSrv.Serve(listen)
		if err != nil {
			global.Logger.Fatal("grpc failed to serve", zap.Error(err))
		}
	}()
}

func (srv *AuthSrv) metricsServer(promCollector prometheus.Collector) {
	mux := http.NewServeMux()
	re := prometheus.NewRegistry()
	re.MustRegister(promCollector)
	mux.Handle("/_/metrics", promhttp.HandlerFor(re, promhttp.HandlerOpts{
		EnableOpenMetrics: true,
	}))
	srv.httpSrv = &http.Server{
		Addr:    fmt.Sprintf(":%d", srv.httpPort),
		Handler: mux,
	}
	srv.wg.Add(1)
	go func() {
		defer srv.wg.Done()
		err := srv.httpSrv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			global.Logger.Fatal("metrics failed to serve", zap.Error(err))
		}
	}()
}

func (srv *AuthSrv) initDB(name string) {
	cfg := global.GetPostgresqlConnSetting(name)
	db, err := rdb.ConnPostgresql(cfg.DSN, cfg.MaxOpenConn, cfg.MaxIdleConn, cfg.MaxLifetime)
	if err != nil {
		global.Logger.Fatal("failed to connect to auth db", zap.Error(err))
	}
	srv.authDB = db
}
