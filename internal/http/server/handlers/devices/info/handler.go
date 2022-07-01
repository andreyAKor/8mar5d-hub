package info

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

// Handle Получение информации о контроллере.
func (h *Handler) Handle() func(http.ResponseWriter, *http.Request) (interface{}, error) {
	return func(w http.ResponseWriter, r *http.Request) (interface{}, error) {
		res, err := h.devicesClient.Info()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error().Err(err).Msg("devices info fail")

			return nil, errors.Wrap(err, "devices info fail")
		}

		return map[string]interface{}{
			"response": res,
		}, nil
	}
}
