package app

import (
	"context"
	"errors"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

type App struct {
	di         *diContainer
	httpServer *http.Server
}

func New() *App {
	a := &App{}

	a.initDI()
	a.initHttpServer()

	a.di.DB() // Initialize DB connection

	return a
}

func (a *App) Run() error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)

	go func() {
		a.di.Logger().Info("http server started", "addr", a.httpServer.Addr)
		errCh <- a.httpServer.ListenAndServe()
	}()

	select {
	case err := <-errCh:
		if !errors.Is(err, http.ErrServerClosed) {
			a.di.Close()
			return err
		}
	case <-ctx.Done():
		a.di.Logger().Info("shutdown signal received")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.httpServer.Shutdown(shutdownCtx); err != nil {
		a.di.Close()
		return err
	}

	a.di.Close()

	return nil
}

func (a *App) initDI() {
	a.di = &diContainer{}
}

func (a *App) initHttpServer() {
	cfg := a.di.ServerCfg()

	readHeaderTimeout := cfg.ReadHeaderTimeout()
	if readHeaderTimeout <= 0 {
		readHeaderTimeout = 5 * time.Second
	}

	a.httpServer = &http.Server{
		Addr:              ":" + cfg.Port(),
		Handler:           a.di.Router(),
		ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout:       cfg.TimeOut(),
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       cfg.Idletimeout(),
	}
}
