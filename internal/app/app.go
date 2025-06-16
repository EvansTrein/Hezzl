package app

import (
	"context"
	"fmt"
	"hezzl/config"
	"hezzl/internal/controller"
	"hezzl/internal/event"
	"hezzl/internal/repository"
	"hezzl/internal/service"
	"hezzl/pkg/broker/nats"
	"hezzl/pkg/db/clickhouse"
	"hezzl/pkg/db/postgres"
	"hezzl/pkg/db/redis"
	"hezzl/pkg/logger"
	"log"
	"log/slog"
	"net/http"
	"time"
)

const (
	gracefulShutdownTimer = time.Second * 10
)

type App struct {
	logger     *slog.Logger
	http       *http.Server
	postgres   *postgres.PostgresDB
	clickhouse *clickhouse.ClickhouseDB
	redis      *redis.RedisDB
	nats       *nats.NatsBroker
}

func New() *App {
	conf := config.GetConfig()

	postgres, err := postgres.New(conf.StoragePath)
	if err != nil {
		log.Fatal("failed connect to postgres")
	}

	clickhouse, err := clickhouse.New(
		conf.Clickhouse.Host,
		conf.Clickhouse.Port,
		conf.Clickhouse.DB,
		conf.Clickhouse.User,
		conf.Clickhouse.Password,
	)
	if err != nil {
		log.Fatal("failed connect to clickhouse")
	}

	redis, err := redis.New(
		conf.Redis.Host,
		conf.Redis.Port,
		conf.Redis.Password,
		conf.Redis.TTLKeys,
		conf.Redis.NumberDB,
	)
	if err != nil {
		log.Fatal("failed connect to redis")
	}

	nats, err := nats.New(conf.Nats.Host, conf.Nats.Port, conf.Nats.NameMess)
	if err != nil {
		log.Fatal("failed connect to nats")
	}

	// Init repository
	goodRepo := repository.NewGoodsRepo(&repository.GoodsRepoDeps{
		Logger:     logger.GetLogger(),
		PostgresDB: postgres,
	})

	cacheRepo := repository.NewCacheRepo(&repository.CacheRepoDeps{
		Logger:  logger.GetLogger(),
		RedisDB: redis,
	})

	logsRepo := repository.NewLogsRepo(&repository.LogsRepoDeps{
		Logger:       logger.GetLogger(),
		ClickhouseDB: clickhouse,
	})

	// Init event
	eventLog := event.NewLogging(&event.LoggingDeps{
		Logger:     logger.GetLogger(),
		NatsBroker: nats,
		ILogsRepo:  logsRepo,
	})

	// Init service
	goodService := service.NewGoods(&service.GoodsDeps{
		Logger:        logger.GetLogger(),
		IGoodsRepo:    goodRepo,
		ICacheRepo:    cacheRepo,
		IEventManager: eventLog,
	})

	// Init controllers
	baseController := controller.NewBaseController(&controller.BaseControllerDeps{
		Logger: logger.GetLogger(),
	})

	goodsController := controller.NewGoods(&controller.GoodsDeps{
		BaseController: baseController,
		IGoodsService:  goodService,
	})

	handler := NewActiveHandlers(&activeHandlersDeps{
		Goods: goodsController,
	})

	// Init server
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", conf.HttpServer.Host, conf.HttpServer.Port),
		Handler: handler.InitRouters(),
	}

	return &App{
		logger:     logger.GetLogger(),
		http:       server,
		postgres:   postgres,
		clickhouse: clickhouse,
		redis:      redis,
		nats:       nats,
	}
}

func (a *App) Start() error {
	a.logger.Info("app: successfully started", "port", config.GetConfig().HttpServer.Port)
	if err := a.http.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func (a *App) Stop() error {
	a.logger.Info("app: stop started")

	ctx, cancel := context.WithTimeout(context.Background(), gracefulShutdownTimer)
	defer cancel()

	if err := a.http.Shutdown(ctx); err != nil {
		a.logger.Error("failed to stop http server", "error", err)
		return err
	}

	if err := a.nats.Close(); err != nil {
		a.logger.Error("failed to stop nats", "error", err)
		return err
	}

	if err := a.postgres.Close(); err != nil {
		a.logger.Error("failed to stop postgres", "error", err)
		return err
	}

	if err := a.clickhouse.Close(); err != nil {
		a.logger.Error("failed to stop clickhouse", "error", err)
		return err
	}

	if err := a.redis.Close(); err != nil {
		a.logger.Error("failed to stop redis", "error", err)
		return err
	}

	a.logger.Info("app: stop successful")
	return nil
}
