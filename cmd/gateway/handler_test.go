package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	mockgrpc "github.com/reddtsai/reddservice/api/proto/mock"
	mockconn "github.com/reddtsai/reddservice/cmd/gateway/mock"
	"github.com/reddtsai/reddservice/internal/global"
)

type TestGateway struct {
	suite.Suite
	ctrl         *gomock.Controller
	mockGrpcConn *mockconn.MockIGrpcClientConn
	mockGrpcAuth *mockgrpc.MockAuthServiceClient

	gateway *Gateway
}

func TestGatewaySuite(t *testing.T) {
	suite.Run(t, new(TestGateway))
}

func (t *TestGateway) SetupSuite() {
	t.ctrl = gomock.NewController(t.T())
	t.mockGrpcConn = mockconn.NewMockIGrpcClientConn(t.ctrl)
	t.mockGrpcAuth = mockgrpc.NewMockAuthServiceClient(t.ctrl)
	// global.Startup("../../conf.d")
	global.Logger, _ = zap.NewDevelopment()
	t.gateway = NewGateway()
	t.gateway.register(NewHandler(t.mockGrpcConn))
}

func (t *TestGateway) TearDownSuite() {
	t.ctrl.Finish()
}

func (t *TestGateway) TestNewGateway() {
	t.NotNil(t.gateway)
}

func (t *TestGateway) TestSignUp_OK() {
	payload := SignUpRequest{
		Account: "test",
		Email:   "test@test.com",
	}
	body, _ := json.Marshal(payload)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/v1/sign-up", bytes.NewReader(body))

	t.mockGrpcConn.EXPECT().GetAuthClient().Return(t.mockGrpcAuth)
	t.mockGrpcAuth.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(nil, nil)

	t.gateway.Handler.ServeHTTP(w, r)

	assert.Equal(t.T(), http.StatusOK, w.Code)
}

func (t *TestGateway) TestSignUp_BadRequest() {
	payload := SignUpRequest{
		Account: "test",
		Email:   "",
	}
	body, _ := json.Marshal(payload)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/v1/sign-up", bytes.NewReader(body))

	t.gateway.Handler.ServeHTTP(w, r)

	assert.Equal(t.T(), http.StatusBadRequest, w.Code)
}
