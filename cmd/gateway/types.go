package main

type SignUpRequest struct {
	Account string `json:"account" binding:"required"`
	Email   string `json:"email" binding:"required"`
}

type Response struct {
	Code   int    `json:"code"`
	Error  string `json:"error"`
	Result any    `json:"result"`
}

type SignUpResponse struct {
}
