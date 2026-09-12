package internal

import (
	"net/http"
	"zygarde/lib/httpx"
	"zygarde/lib/o11y"

	"github.com/go-chi/chi/v5"
)

type PingRequest struct {
	URL  string `json:"url"`
	Path string `json:"path"`
}

type PingResponse struct {
	Message string `json:"message"`
}

type PongResponse = PingResponse

func AddRoutes(r chi.Router) {
	r.Post("/ping", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		log := o11y.LoggerFromContext(ctx)

		req, err := httpx.ParseRequest[PingRequest](w, r)
		if err != nil {
			return
		}

		cl := NewCellClient(req.URL)

		pong, err := cl.Ping(ctx, req.Path)
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

		httpx.WriteOKResponse(w, PingResponse{
			Message: pong.Message,
		})
	})

	r.Get("/pong", func(w http.ResponseWriter, r *http.Request) {
		httpx.WriteOKResponse(w, PongResponse{
			Message: "pong",
		})
	})
}
