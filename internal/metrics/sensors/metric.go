package sensors

import (
	"context"
	"time"

	"github.com/andreyAKor/8mar5d-hub/internal/clients/devices"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/rs/zerolog/log"
)

var metrics = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Namespace: "mar8d5_hub",
	Name:      "sensors",
	Help:      "Variables of UPS list",
}, []string{"device", "sensor"})

type Metric struct {
	interval      time.Duration
	devicesClient *devices.Client
}

func New(interval string, devicesClient *devices.Client) (*Metric, error) {
	intervalDur, err := time.ParseDuration(interval)
	if err != nil {
		return nil, errors.Wrapf(err, "interval parsing fail (%s)", interval)
	}

	return &Metric{
		interval:      intervalDur,
		devicesClient: devicesClient,
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
		sensors, err := m.devicesClient.Sensors()
		if err != nil {
			log.Warn().Err(err).Msg("sensors fail")
			continue
		}
		if sensors.Data == nil {
			log.Warn().Err(err).Msg("sensors data is empty")
			continue
		}

		for _, sensor := range *sensors.Data {
			if sensor.Type == devices.SensorTypeDHT11 {
				data, err := m.devicesClient.SensorDHT11()
				if err != nil {
					log.Warn().Err(err).Msg("sensor DHT-11 fail")
					continue
				}
				if data.Data == nil {
					log.Warn().Err(err).Msg("sensors data is empty on DHT-11")
					continue
				}

				metrics.WithLabelValues(m.devicesClient.Host(), "temperature").Set(data.Data.Temperature.Value)
				metrics.WithLabelValues(m.devicesClient.Host(), "humidity").Set(data.Data.Humidity.Value)
			}
			if sensor.Type == devices.SensorTypeMQ2 {
				data, err := m.devicesClient.SensorMQ2()
				if err != nil {
					log.Warn().Err(err).Msg("sensor MQ-2 fail")
					continue
				}
				if data.Data == nil {
					log.Warn().Err(err).Msg("sensors data is empty on MQ-2")
					continue
				}

				metrics.WithLabelValues(m.devicesClient.Host(), "ratio").Set(float64(data.Data.Ratio.Value))
				metrics.WithLabelValues(m.devicesClient.Host(), "lpg").Set(float64(data.Data.Lpg.Value))
				metrics.WithLabelValues(m.devicesClient.Host(), "methane").Set(float64(data.Data.Methane.Value))
				metrics.WithLabelValues(m.devicesClient.Host(), "smoke").Set(float64(data.Data.Smoke.Value))
				metrics.WithLabelValues(m.devicesClient.Host(), "hydrogen").Set(float64(data.Data.Hydrogen.Value))
			}
		}
	}

	return nil
}
