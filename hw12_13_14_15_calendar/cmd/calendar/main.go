package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/FedoseevAlex/OtusGolangHomeWork/hw12_13_14_15_calendar/internal/app"
	"github.com/FedoseevAlex/OtusGolangHomeWork/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/FedoseevAlex/OtusGolangHomeWork/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/FedoseevAlex/OtusGolangHomeWork/hw12_13_14_15_calendar/internal/storage/memory"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config, err := NewConfig(configFile)
	if err != nil {
		log.Fatal(err)
	}
	appLogger := logger.New(config.Logger.Level, config.Logger.File)

	storage := memorystorage.New()
	calendar := app.New(appLogger, storage)

	server := internalhttp.NewServer(calendar, config.Server.Host, config.Server.Port)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGINT, syscall.SIGHUP)

		select {
		case <-ctx.Done():
			return
		case <-signals:
		}

		signal.Stop(signals)
		cancel()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			appLogger.Error("failed to stop http server: ", logger.LogArgs{"error": err})
		}
	}()

	appLogger.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		appLogger.Error("failed to start http server", logger.LogArgs{"error": err})
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
