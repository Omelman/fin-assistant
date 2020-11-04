package api

import (
	"context"
	"github.com/fin-assistant/internal/config"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/logan/v3"
	"net/http"
)

type service struct {
	cfg     config.Config
	ctx     context.Context

	logger *logan.Entry
}

func NewService(ctx context.Context, cfg config.Config) (*service, error) {
	return &service{
		cfg:    cfg,
		ctx:    ctx,
		logger: cfg.Log(),
	}, nil
}

func (s *service) Run() error {
	s.logger.Info("Starting...")

	r := Router(s.cfg)


	err := http.Serve(s.cfg.Listener(), r)
	if err != nil {
		return errors.Wrap(err, "server stopped with error")
	}
	return nil
}
