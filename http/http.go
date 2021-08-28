package http

import (
	"asscRegsitration/config"
	"asscRegsitration/http/jrpc"
	"asscRegsitration/model"
	"asscRegsitration/service"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type httpService struct {
	logger *zap.Logger
	cfg    *config.Config
	db     model.DB
	r      *mux.Router
	srv    *http.Server
	svc    service.Service
}

func NewService(ctx context.Context, logger *zap.Logger, cfg *config.Config, db model.DB, svc service.Service) *httpService {
	l := logger.With(zap.String("component", "http"))
	addr := fmt.Sprintf("%s:%d", cfg.Http.Host, cfg.Http.Port)

	mr := jrpc.CreateMethodRepository(ctx, logger, svc)
	r := mux.NewRouter()
	r.Handle(fmt.Sprintf("/api/%s", cfg.ApiVersion), mr)

	srv := &http.Server{
		Handler: r,
		Addr:    addr,

		WriteTimeout: cfg.Http.WriteTimeout,
		ReadTimeout:  cfg.Http.ReadTimeout,
	}

	r.Methods("GET", "POST", "OPTIONS", "DELETE", "PUT")

	return &httpService{
		logger: l,
		cfg:    cfg,
		db:     db,
		r:      r,
		srv:    srv,
		svc:    svc,
	}

}

func (srv *httpService) Start() error {
	if srv.cfg.TLS.Enabled {
		srv.r.Schemes("https")
		return srv.srv.ListenAndServeTLS(srv.cfg.TLS.CertPath, srv.cfg.TLS.KeyPath)
	} else {
		return srv.srv.ListenAndServe()
	}
}

func (srv *httpService) Stop() {
	time.Sleep(6 * time.Second)

	err := srv.srv.Shutdown(context.Background())
	if err != nil {
		srv.logger.With(zap.Error(err)).Info("srv stop failed")
	}
}
