package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"time"

	_ "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/abdukhashimov/exporter_psql_clickhouse/config"
	"github.com/abdukhashimov/exporter_psql_clickhouse/pkg/cron"
	"github.com/abdukhashimov/exporter_psql_clickhouse/pkg/exporter"
	"github.com/abdukhashimov/exporter_psql_clickhouse/pkg/logger"
	"github.com/abdukhashimov/exporter_psql_clickhouse/pkg/logger/factory"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func init() {
	cfg := config.Load()

	log, err := factory.Build(&cfg.Logging)
	if err != nil {
		panic(err)
	}

	logger.SetLogger(log)

	logger.Log.Info("set logger successfully...")

	db, err := sqlx.Connect("postgres", cfg.PsqlConfig.ConnString)
	if err != nil {
		panic(err)
	}

	conn, err := sqlx.Open("clickhouse", cfg.Clickhouse.ConnString)
	if err != nil {
		panic(err)
	}

	exporterObj := exporter.New(db, conn, cfg)

	cronJob := cron.New(exporterObj)
	cronJob.RunTableExporter(cfg.Exporter.ExportPerid, cfg.Exporter.TableName)
}

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(7))
	defer cancel()

	var wg sync.WaitGroup

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		logger.Log.Info("shutting down")

		logger.Log.Info("shutdown successfully called")

		wg.Done()
	}(&wg)

	go func() {
		wg.Wait()
		cancel()
	}()

	<-ctx.Done()
	os.Exit(0)
}
