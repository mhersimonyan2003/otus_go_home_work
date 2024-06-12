package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	"github.com/mhersimonyan2003/otus_go_home_work/hw12_13_14_15_calendar/internal/app"
	internalConfig "github.com/mhersimonyan2003/otus_go_home_work/hw12_13_14_15_calendar/internal/config"
	"github.com/mhersimonyan2003/otus_go_home_work/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/mhersimonyan2003/otus_go_home_work/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/mhersimonyan2003/otus_go_home_work/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/mhersimonyan2003/otus_go_home_work/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/pressly/goose"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config, err := internalConfig.NewConfig(configFile)

	if err != nil {
		os.Stderr.WriteString("failed to read config: " + err.Error() + "\n")
		os.Exit(1)
	}

	logg := logger.New(config.Logger.Level)

	var storage app.Storage
	switch config.Storage.Type {
	case "memory":
		storage = memorystorage.New()
	case "sql":
		dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			config.Database.Host, config.Database.Port, config.Database.User,
			config.Database.Password, config.Database.DBName, config.Database.SSLMode)
		if err := runMigrations(dsn); err != nil {
			logg.Error("failed to run migrations: " + err.Error())
			os.Exit(1)
		}
		var err error
		storage, err = sqlstorage.New(dsn)
		if err != nil {
			logg.Error("failed to connect to database: " + err.Error())
			os.Exit(1)
		}
	default:
		logg.Error("unsupported storage type: " + config.Storage.Type)
		os.Exit(1)
	}

	calendar := app.New(logg, storage)

	server := internalhttp.NewServer(logg, calendar)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}

func runMigrations(dsn string) error {
	db, err := goose.OpenDBWithDriver("postgres", dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := goose.Up(db, "./migrations"); err != nil {
		return err
	}
	return nil
}
