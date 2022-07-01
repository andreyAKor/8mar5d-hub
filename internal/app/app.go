package app

import (
	"context"
	"io"

	"github.com/rs/zerolog/log"

	"github.com/andreyAKor/8mar5d-hub/internal/http/server"
	metricsSensors "github.com/andreyAKor/8mar5d-hub/internal/metrics/sensors"
)

var _ io.Closer = (*App)(nil)

type App struct {
	srv            *server.Server
	sensorsMetrics *metricsSensors.Metric
}

func New(srv *server.Server, sensorsMetrics *metricsSensors.Metric) (*App, error) {
	return &App{
		srv:            srv,
		sensorsMetrics: sensorsMetrics,
	}, nil
}

// Run application.
func (a *App) Run(ctx context.Context) error {
	go func() {
		if err := a.srv.Run(ctx); err != nil {
			log.Fatal().Err(err).Msg("http-server listen fail")
		}
	}()
	go func() {
		if err := a.sensorsMetrics.Run(ctx); err != nil {
			log.Fatal().Err(err).Msg("nut metrics running fail")
		}
	}()

	return nil
}

// Close application.
func (a *App) Close() error {
	return nil
}
