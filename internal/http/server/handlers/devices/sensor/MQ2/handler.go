package MQ2

import (
	"net/http"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	"github.com/andreyAKor/8mar5d-hub/internal/clients/devices"
)

type Handler struct {
	devicesClient *devices.Client
}

func New(devicesClient *devices.Client) *Handler {
	return &Handler{
		devicesClient: devicesClient,
	}
}

// Handle Получение метрик сенсора MQ-2.
func (h *Handler) Handle() func(http.ResponseWriter, *http.Request) (interface{}, error) {
	return func(w http.ResponseWriter, r *http.Request) (interface{}, error) {
		res, err := h.devicesClient.SensorMQ2()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error().Err(err).Msg("devices sensor MQ-2 fail")

			return nil, errors.Wrap(err, "devices sensor MQ-2 fail")
		}

		return map[string]interface{}{
			"response": res,
		}, nil
	}
}
