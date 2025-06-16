package main

import (
	"fmt"
	"hezzl/config"
	"hezzl/internal/event"
	"hezzl/internal/repository"
	"hezzl/pkg/broker/nats"
	"hezzl/pkg/db/clickhouse"
	"hezzl/pkg/logger"
	"log"
	"os"
	"os/signal"
	"syscall"

	n "github.com/nats-io/nats.go"
)

func main() {
	config.MustLoad()
	conf := config.GetConfig()

	logger.InitLog(logger.LogConfig{
		Mode:     conf.Env,
		LogPath:  conf.LogOutput,
		LogLevel: conf.LogLevel,
	})

	myLog := logger.GetLogger()
	chErr := make(chan error, 1)

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

	nats, err := nats.New(conf.Nats.Host, conf.Nats.Port, conf.Nats.NameMess)
	if err != nil {
		log.Fatal("failed connect to nats")
	}

	logsRepo := repository.NewLogsRepo(&repository.LogsRepoDeps{
		Logger:       logger.GetLogger(),
		ClickhouseDB: clickhouse,
	})

	loggingEvent := event.NewLogging(&event.LoggingDeps{
		Logger:     logger.GetLogger(),
		NatsBroker: nats,
		ILogsRepo:  logsRepo,
	})

	go func() {
		_, err = loggingEvent.Broker.Conn.Subscribe(nats.NameMess, func(m *n.Msg) {
			loggingEvent.SendLogToDB(m.Data)
		})
		if err != nil {
			chErr <- fmt.Errorf("failed to subscribe to NATS topic: %w", err)
		}

		myLog.Info("subscribe successfully")
	}()

	stop := make(chan os.Signal, 3)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	select {
	case err := <-chErr:
		myLog.Error("func main error", "error", err)
		if err := nats.Close(); err != nil {
			myLog.Error("failed to stop nats", "error", err)
		}

		if err := clickhouse.Close(); err != nil {
			myLog.Error("failed to stop clickhouse", "error", err)
		}

		os.Exit(1)
	case sig := <-stop:
		myLog.Info("received signal, shutting down", "signal", sig)
		var err error

		if err = nats.Close(); err != nil {
			myLog.Error("failed to stop nats", "error", err)
		}

		if err = clickhouse.Close(); err != nil {
			myLog.Error("failed to stop clickhouse", "error", err)
		}

		if err != nil {
			os.Exit(1)
		}

		myLog.Info("shutdown completed successfully")
	}
}
