package main

import (
	"context"
	"errors"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"os-simulator-plan/internal/transport/realtime"

	"go.uber.org/zap"
)

func main() {
	addr := flag.String("addr", ":8080", "HTTP listen address")
	flag.Parse()
	logger, err := zap.NewProduction()
	if err != nil {
		_, _ = os.Stderr.WriteString("failed to initialize logger\n")
		os.Exit(1)
	}
	defer func() {
		_ = logger.Sync()
	}()

	manager := realtime.NewSessionManager()
	transport := realtime.NewServer(manager)

	httpServer := &http.Server{
		Addr:              *addr,
		Handler:           transport.Handler(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	logger.Info("server listening", zap.String("addr", *addr))
	errCh := make(chan error, 1)
	go func() {
		errCh <- httpServer.ListenAndServe()
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errCh:
		if err == nil || errors.Is(err, http.ErrServerClosed) {
			return
		}
		logger.Error("server failed", zap.Error(err))
		os.Exit(1)
	case sig := <-sigCh:
		logger.Info("received shutdown signal", zap.String("signal", sig.String()))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		logger.Error("server shutdown failed", zap.Error(err))
		os.Exit(1)
	}
}
