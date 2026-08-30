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

	a.di.DB()

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

	a.httpServer = &http.Server{
		Addr:         cfg.Address(),
		Handler:      a.di.Router(),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
}



package app

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/goNiki/subservice/internal/infrastructure/config"
	"github.com/goNiki/subservice/internal/infrastructure/database"
	"github.com/goNiki/subservice/internal/infrastructure/database/txmanager"
	"github.com/goNiki/subservice/internal/infrastructure/logger"
	"github.com/goNiki/subservice/internal/infrastructure/migrator"
	"github.com/goNiki/subservice/internal/repository"
	subRepo "github.com/goNiki/subservice/internal/repository/subscription"
	txRepo "github.com/goNiki/subservice/internal/repository/transaction"
	"github.com/goNiki/subservice/internal/service"
	"github.com/goNiki/subservice/internal/service/subscription"
	"github.com/goNiki/subservice/internal/transport/handlers"
	"github.com/goNiki/subservice/internal/transport/router"
)

const configPath = "./.env"
const MigDir = "./migrations"

type diContainer struct {
	//инфраструктура
	loggerCfg   config.Logger
	postgresCfg config.Postgres
	serverCfg   config.Server
	logger      *logger.Logger
	db          *database.DB
	txManager   *txmanager.TxManager
	migrator    *migrator.Migrator

	//router
	router http.Handler

	//handler
	subHandler *handlers.Handler

	//service
	subService service.SubscriptionService

	//repository
	subRepository repository.SubscriptionRepository
	txRepository  repository.TransactionRepository
}

func (d *diContainer) initConfig() {
	if d.loggerCfg != nil && d.postgresCfg != nil && d.serverCfg != nil {
		return
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		slog.Error("failed to get config: ", "error", err)
		os.Exit(1)
	}

	d.loggerCfg = cfg.Logger
	d.postgresCfg = cfg.Postgres
	d.serverCfg = cfg.Server
}

func (d *diContainer) LoggerCfg() config.Logger {
	d.initConfig()

	return d.loggerCfg
}

func (d *diContainer) PostgresCfg() config.Postgres {
	d.initConfig()

	return d.postgresCfg
}

func (d *diContainer) ServerCfg() config.Server {
	d.initConfig()

	return d.serverCfg
}

func (d *diContainer) Logger() *logger.Logger {
	if d.logger == nil {
		logger, err := logger.InitLogger(d.LoggerCfg())
		if err != nil {
			slog.Error("Failed to create logger: " + err.Error())
			os.Exit(1)
		}

		d.logger = logger
	}

	return d.logger
}

func (d *diContainer) DB() *database.DB {
	if d.db == nil {
		db, err := database.InitDatabase(d.PostgresCfg())
		if err != nil {
			slog.Error("failed to connection to DB: " + err.Error())
			os.Exit(1)
		}

		d.db = db
	}

	return d.db
}

func (d *diContainer) TxManager() *txmanager.TxManager {
	if d.txManager == nil {
		d.txManager = txmanager.NewTxManager(d.DB())
	}

	return d.txManager
}

func (d *diContainer) Migrator() *migrator.Migrator {
	if d.migrator == nil {
		migrator, err := migrator.NewMigrator(d.DB().Pool, MigDir)
		if err != nil {
			slog.Error("failed to connect migration: " + err.Error())
			os.Exit(1)
		}
		d.migrator = migrator
	}

	return d.migrator
}

func (d *diContainer) SubRepository() repository.SubscriptionRepository {
	if d.subRepository == nil {
		d.subRepository = subRepo.NewSubscriptionRepository()
	}

	return d.subRepository

}

func (d *diContainer) TxRepository() repository.TransactionRepository {
	if d.txRepository == nil {
		d.txRepository = txRepo.NewTransactionSubRepository()
	}

	return d.txRepository
}

func (d *diContainer) SubService() service.SubscriptionService {
	if d.subService == nil {
		d.subService = subscription.NewSubscriptionService(
			d.TxManager(),
			d.DB(),
			d.SubRepository(),
			d.TxRepository(),
		)
	}

	return d.subService
}

func (d *diContainer) SubHandler() *handlers.Handler {
	if d.subHandler == nil {
		d.subHandler = handlers.NewHandler(
			d.SubService(),
			d.Logger(),
		)
	}

	return d.subHandler
}

func (d *diContainer) Router() http.Handler {
	if d.router == nil {
		d.router = router.NewRouter(d.SubHandler())
	}
	return d.router
}

func (d *diContainer) Close() {
	if d.db != nil {
		d.db.Close()
	}
	if d.logger != nil {
		_ = d.logger.Close()
	}
}
