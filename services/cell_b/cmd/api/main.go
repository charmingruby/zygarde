package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
	"zygarde/lib/httpx"
	"zygarde/lib/o11y"
	"zygarde/lib/validator"

	"github.com/charmingruby/zygard/cell_b/config"
	"github.com/charmingruby/zygard/cell_b/internal"
)

const shutdownTimeout = 30 * time.Second

func main() {
	if err := run(); err != nil {
		os.Exit(1)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer stop()

	log := o11y.InitLogger()

	cfg, err := config.Load()
	if err != nil {
		log.Error("config: error loading config", "error", err)
		return err
	}

	val := validator.New()
	srv, router := httpx.NewServer(cfg.Port, val)

	internal.AddRoutes(router)

	shutdownErrCh := make(chan error, 1)
	go shutdown(ctx, log, shutdownErrCh, srv)

	log.Info("server: running...", "port", cfg.Port)
	if err := srv.Start(); err != nil {
		log.Error("server: start error", "error", err)

		stop()
		<-shutdownErrCh
		return err
	}

	if err := <-shutdownErrCh; err != nil {
		log.Error("shutdown: shutdown error", "error", err)
		return err
	}

	log.Info("shutdown: gracefully shutdown")
	return nil
}

func shutdown(
	ctx context.Context,
	log *slog.Logger,
	errCh chan error,
	srv *httpx.Server,
) {
	<-ctx.Done()
	log.Info("shutdown: signal received, starting graceful shutdown...")

	ctxTimeout, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	var err error

	if closeErr := srv.Close(ctxTimeout); closeErr != nil {
		if errors.Is(closeErr, context.DeadlineExceeded) {
			closeErr = errors.New("http server: deadline exceeded, forcing shutdown")
		}

		err = errors.Join(err, closeErr)
	}

	errCh <- err
}
