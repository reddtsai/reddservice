package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/reddtsai/reddservice/api/proto"
	"github.com/reddtsai/reddservice/internal/global"
)

const (
	AccessControlMaxAge = 12 * time.Hour
)

type GatewayHandler struct {
	mux http.Handler
}

func NewGatewayHandler() *GatewayHandler {
	h := &GatewayHandler{}
	h.register()
	return h
}

func (h *GatewayHandler) register() {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /v1/sign-up", h.handleSignUp)

	h.mux = mux
	h.mux = SetAccessControl(h.mux)
	h.mux = Middleware(h.mux)
}

func (h *GatewayHandler) handleSignUp(w http.ResponseWriter, r *http.Request) {
	authClient, err := connAuthClient()
	if err != nil {
		global.Logger.Error("conn auth client failed", zap.Error(err))
		// TODO
		return
	}
	_, err = authClient.SignUp(r.Context(), &pb.SignUpRequest{
		Account: r.FormValue("account"),
		Email:   r.FormValue("email"),
	})
	if err != nil {
		global.Logger.Error("sign up failed", zap.Error(err))
		// TODO
		return
	}
	// TODO
	JSON(w, http.StatusOK, map[string]string{"message": "sign up"})
}

func JSON(w http.ResponseWriter, code int, obj any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		global.Logger.Error("json marshal failed", zap.Error(err))
		return
	}
	_, err = w.Write(jsonBytes)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		global.Logger.Error("response write failed", zap.Error(err))
	}
}

func SetAccessControl(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,HEAD,OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin,Content-Length,Content-Type,Authorization")
		w.Header().Set("Access-Control-Expose-Headers", "*")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		maxAge := strconv.FormatInt(int64(AccessControlMaxAge/time.Second), 10)
		w.Header().Set("Access-Control-Max-Age", maxAge)

		next.ServeHTTP(w, r)
	})
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		resp := &RespW{w, http.StatusOK}
		next.ServeHTTP(resp, r)
		duration := time.Since(start)
		ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr))
		if err != nil {
			ip = "unknown"
		}
		msg := fmt.Sprintf("[%s] %15s |%3d| %13s | %-7s %s", ServiceName, ip, resp.status, duration.String(), r.Method, r.URL.Path)
		global.Logger.Info(msg, zap.String("method", r.Method), zap.String("path", r.URL.Path), zap.Duration("duration", duration))
	})
}

type RespW struct {
	http.ResponseWriter
	status int
}

func (rw *RespW) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func connAuthClient() (pb.AuthServiceClient, error) {
	addr := "localhost:8080"
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return pb.NewAuthServiceClient(conn), nil
}
