package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hilltracer/otus-go/hw12_13_14_15_calendar/internal/app"
	"github.com/hilltracer/otus-go/hw12_13_14_15_calendar/internal/config"
	"github.com/hilltracer/otus-go/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/hilltracer/otus-go/hw12_13_14_15_calendar/internal/server/http"
	"github.com/hilltracer/otus-go/hw12_13_14_15_calendar/internal/server/internalgrpc"
	"github.com/hilltracer/otus-go/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/hilltracer/otus-go/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/hilltracer/otus-go/hw12_13_14_15_calendar/internal/storage/sql"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.yaml", "Path to configuration file")
}

func main() {
	os.Exit(run())
}

func run() int {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return 0
	}

	cfg, err := config.NewConfig(configFile)
	if err != nil {
		fmt.Printf("failed to read config: %v\n", err)
		return 1
	}

	logg := logger.New(cfg.Logger.Level)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	var storage storage.Repository
	switch cfg.Storage.Type {
	case "sql":
		pwd := os.Getenv(cfg.Storage.PG.PasswordEnv)
		dsn := fmt.Sprintf(
			"postgresql://%s:%s@%s:%d/%s?sslmode=%s",
			cfg.Storage.PG.User, pwd, cfg.Storage.PG.Host, cfg.Storage.PG.Port, cfg.Storage.PG.DBName, cfg.Storage.PG.SSLMode,
		)
		pgStore, err := sqlstorage.Connect(ctx, dsn)
		if err != nil {
			logg.Error("db connect: " + err.Error())
			return 1
		}
		storage = pgStore
	default:
		storage = memorystorage.New()
	}

	calendar := app.New(logg, storage)
	addr := net.JoinHostPort(cfg.HTTP.Host, cfg.HTTP.Port)
	server := internalhttp.NewServer(logg, calendar, addr)

	go func() {
		logg.Info("HTTP server starting on " + addr)
		if err := server.Start(ctx); err != nil {
			logg.Error("failed to start http server: " + err.Error())
		}
	}()

	grpcAddr := net.JoinHostPort(cfg.GRPC.Host, cfg.GRPC.Port)
	gsrv := internalgrpc.New(calendar, logg)

	go func() {
		logg.Info("gRPC server starting on " + grpcAddr)
		if err := gsrv.Start(grpcAddr); err != nil {
			logg.Error("failed to start grpc server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	<-ctx.Done()

	// graceful shutdown
	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	_ = server.Stop(ctxTimeout)
	gsrv.Stop()

	return 0
}
