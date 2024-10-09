package global

import (
	"net/http"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
)

const (
	KEY_LOG_LEVEL = "LOG_LEVEL"
)

var (
	once        sync.Once
	Logger      *zap.Logger
	SugarLogger *zap.SugaredLogger
)

func init() {
	// TODO: add configuration
	viper.AutomaticEnv()
	viper.SetDefault(KEY_LOG_LEVEL, -1) // -1 ~ 5
	newZapLogger()
}

func newZapLogger() {
	level := -1
	if isSet := viper.IsSet(KEY_LOG_LEVEL); isSet {
		level = viper.GetInt(KEY_LOG_LEVEL)
	}
	once.Do(func() {
		config := zap.NewProductionConfig()
		config.Level = zap.NewAtomicLevelAt(zapcore.Level(level))
		logger, err := config.Build()
		if err != nil {
			panic(err)
		}
		Logger = logger
		zap.ReplaceGlobals(Logger)
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
