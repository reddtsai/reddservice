package main

import (
	"net/http"
)

type GatewayHandler struct {
}

func NewGatewayHandler() *GatewayHandler {
	return &GatewayHandler{}
}

func (h *GatewayHandler) registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/v1/sign-up", h.handleSignUp)
}

func (h *GatewayHandler) handleSignUp(w http.ResponseWriter, r *http.Request) {
	// TODO: implement handleSignUp
}
