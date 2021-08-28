package jrpc

import (
	api "asscRegsitration"
	"asscRegsitration/service"
	"context"

	"github.com/intel-go/fastjson"
	"github.com/osamingo/jsonrpc"
	"go.uber.org/zap"
)

type teamRegister struct {
	ctx    context.Context
	logger *zap.Logger
	svc    service.Service
}

func TeamRegisterHandler(ctx context.Context, logger *zap.Logger, svc service.Service) jsonrpc.Handler {
	l := logger.With(zap.String("method", "customers_by_ids"))

	res := &teamRegister{
		ctx:    ctx,
		logger: l,
		svc:    svc,
	}

	return res
}

type TeamRegisterRequest struct {
	Team api.Team `json:"team"`
}

type TeamRegisterResponse struct {
	Success bool `json:"success"`
	Id      int  `json:"id,omitempty"`
}

func (h *teamRegister) ServeJSONRPC(ctx context.Context, p *fastjson.RawMessage) (interface{}, *jsonrpc.Error) {
	var req TeamRegisterRequest
	if err := jsonrpc.Unmarshal(p, &req); err != nil {
		return nil, &jsonrpc.Error{Code: InvalidParamsErrorCode, Message: err.Message}
	}

	id, err := h.svc.TeamRegister(ctx, req.Team)
	if err != nil {
		return &TeamRegisterResponse{Success: false}, nil
	}

	return &TeamRegisterResponse{Success: true, Id: id}, nil
}
