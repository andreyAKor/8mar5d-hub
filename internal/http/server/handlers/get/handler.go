package get

import (
	"net/http"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	"github.com/andreyAKor/8mar5d-hub/internal/http/clients/nut"
)

type Handler struct {
	nutClient *nut.Client
}

func New(nutClient *nut.Client) *Handler {
	return &Handler{
		nutClient: nutClient,
	}
}

func (h *Handler) Handle() func(http.ResponseWriter, *http.Request) (interface{}, error) {
	return func(w http.ResponseWriter, r *http.Request) (interface{}, error) {
		list, err := h.nutClient.GetList(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error().Err(err).Msg("get list fail")

			return nil, errors.Wrap(err, "get list fail")
		}

		return list, nil
	}
}
