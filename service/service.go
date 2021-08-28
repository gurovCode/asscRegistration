package service

import (
	api "asscRegsitration"
	"asscRegsitration/config"
	"asscRegsitration/model"
	"context"

	"go.uber.org/zap"
)

type Service interface {
	TeamRegister(ctx context.Context, team api.Team) (int, error)
}

type service struct {
	cfg    *config.Config
	logger *zap.Logger
	db     model.DB
}

func NewService(logger *zap.Logger, cfg *config.Config, db model.DB) Service {
	l := logger.With(zap.String("component", "service"))

	svc := &service{
		logger: l,
		cfg:    cfg,
		db:     db,
	}

	return svc
}
