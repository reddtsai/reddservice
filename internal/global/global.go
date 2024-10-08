package global

import (
	"net/http"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var (
	once        sync.Once
	Logger      *zap.Logger
	SugarLogger *zap.SugaredLogger
)

func init() {
	newZapLogger()
}

func newZapLogger() {
	once.Do(func() {
		Logger, _ = zap.NewProduction()
		SugarLogger = Logger.Sugar()
	})
}

func MetricsServer(addr string, promCollector prometheus.Collector) *http.Server {
	mux := http.NewServeMux()
	re := prometheus.NewRegistry()
	re.MustRegister(promCollector)
	mux.Handle("/metrics", promhttp.HandlerFor(re, promhttp.HandlerOpts{
		EnableOpenMetrics: true,
	}))

	return &http.Server{
		Addr:    addr,
		Handler: mux,
	}
}

func GrpcServer(interceptor ...grpc.UnaryServerInterceptor) *grpc.Server {
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptor...,
		),
	)

	return grpcServer
}
