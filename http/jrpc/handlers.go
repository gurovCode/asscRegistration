package jrpc

import (
	"asscRegsitration/service"
	"context"
	"net/http"

	"github.com/osamingo/jsonrpc"
	"go.uber.org/zap"
)

func CreateMethodRepository(ctx context.Context, logger *zap.Logger, svc service.Service) http.Handler {
	l := logger.With(zap.String("component", "api"))
	mr := jsonrpc.NewMethodRepository()

	err := mr.RegisterMethod("team_register",
		TeamRegisterHandler(ctx, l, svc),
		TeamRegisterRequest{},
		TeamRegisterResponse{})
	if err != nil {
		logger.With(zap.Error(err)).Error("failed register team_register")
	}

	return mr
}
