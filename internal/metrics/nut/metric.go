package nut

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/rs/zerolog/log"

	"github.com/andreyAKor/8mar5d-hub/internal/http/clients/nut"
)

var metrics = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Namespace: "8mar5d-hub",
	Name:      "ups_variables",
	Help:      "Variables of UPS list",
}, []string{"ups", "variable"})

type Metric struct {
	interval  time.Duration
	nutClient *nut.Client
}

func New(interval string, nutClient *nut.Client) (*Metric, error) {
	intervalDur, err := time.ParseDuration(interval)
	if err != nil {
		return nil, errors.Wrapf(err, "interval parsing fail (%s)", interval)
	}

	return &Metric{
		interval:  intervalDur,
		nutClient: nutClient,
	}, nil
}

func (m *Metric) Run(ctx context.Context) error {
	ticker := time.NewTicker(m.interval)

	go func() {
		select {
		case <-ctx.Done():
			ticker.Stop()
		}
	}()

	for _ = range ticker.C {
		list, err := m.nutClient.GetList(ctx)
		if err != nil {
			log.Warn().Err(err).Msg("get list fail")
			continue
		}

		for _ = range list {
		}
	}

	return nil
}
