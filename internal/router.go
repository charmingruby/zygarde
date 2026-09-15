package internal

import (
	"net/http"

	"github.com/charmingruby/zygard/cell/pkg/httpx"
	"github.com/charmingruby/zygard/cell/pkg/o11y"

	"github.com/go-chi/chi/v5"
)

type PingRequest struct {
	URL  string `json:"url"`
	Path string `json:"path"`
}

type PingResponse struct {
	Message    string `json:"message"`
	ReceiverID string `json:"receiver_id"`
	CallerID   string `json:"called_id"`
}

type PongRequest struct {
	CallerID string `json:"caller_id"`
}

type PongResponse struct {
	Message    string `json:"message"`
	ReceiverID string `json:"receiver_id"`
}

func AddRoutes(r chi.Router, id string) {
	r.Post("/ping", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		log := o11y.LoggerFromContext(ctx)

		req, err := httpx.ParseRequest[PingRequest](w, r)
		if err != nil {
			return
		}

		cl := NewClient(req.URL)

		log.Info("trying to call pong",
			"caller_id", id,
		)

		pong, err := cl.CallPong(ctx, req.Path, id)
		if err != nil {
			log.Error("error from pong",
				"message", err.Error(),
				"url", req.URL,
				"path", req.Path,
			)

			httpx.WriteResponse(w, http.StatusInternalServerError, map[string]string{
				"message": "Internal Server Error",
			})
			return
		}

		log.Info("called pong successfully",
			"caller_id", id,
			"receiver_id", pong.ReceiverID,
		)

		httpx.WriteOKResponse(w, PingResponse{
			Message:    pong.Message,
			ReceiverID: pong.ReceiverID,
			CallerID:   id,
		})
	})

	r.Get("/pong", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		log := o11y.LoggerFromContext(ctx)

		req, err := httpx.ParseRequest[PongRequest](w, r)
		if err != nil {
			return
		}

		log.Info("received pong request",
			"caller_id", req.CallerID,
		)

		httpx.WriteOKResponse(w, PongResponse{
			Message:    "pong",
			ReceiverID: id,
		})
	})
}
