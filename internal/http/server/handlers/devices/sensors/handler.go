package sensors

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

// Handle Получение информации по сенсорам.
func (h *Handler) Handle() func(http.ResponseWriter, *http.Request) (interface{}, error) {
	return func(w http.ResponseWriter, r *http.Request) (interface{}, error) {
		res, err := h.devicesClient.Sensors()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error().Err(err).Msg("devices sensors fail")

			return nil, errors.Wrap(err, "devices sensors fail")
		}

		return map[string]interface{}{
			"response": res,
		}, nil
	}
}
