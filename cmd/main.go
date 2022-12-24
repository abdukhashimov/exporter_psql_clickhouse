package main

import (
	"fmt"
	"log"

	"github.com/abdukhashimov/exporter_psql_clickhouse/config"
	"github.com/abdukhashimov/exporter_psql_clickhouse/pkg/logger"
	"github.com/abdukhashimov/exporter_psql_clickhouse/pkg/logger/factory"
	"github.com/jmoiron/sqlx"
	"github.com/sevlyar/go-daemon"

	_ "github.com/ClickHouse/clickhouse-go/v2"
	_ "github.com/lib/pq"
)

func main() {
	cntxt := &daemon.Context{
		PidFileName: "sample.pid",
		PidFilePerm: 0644,
		LogFileName: "logs/sample.log",
		LogFilePerm: 0640,
		WorkDir:     "./",
		Umask:       027,
		Args:        []string{"[go-daemon sample]"},
	}

	d, err := cntxt.Reborn()
	if err != nil {
		log.Fatal("Unable to run: ", err)
	}
	if d != nil {
		return
	}

	defer cntxt.Release() //nolint:all

	log.Print("- - - - - - - - - - - - - - -")
	log.Print("daemon started")

	serve()
}

func serve() {
	cfg := config.Load()

	log, err := factory.Build(&cfg.Logging)
	if err != nil {
		panic(err)
	}

	logger.SetLogger(log)

	db, err := sqlx.Connect("postgres", cfg.PsqlConfig.ConnString)
	if err != nil {
		panic(err)
	}

	conn, err := sqlx.Open("clickhouse", cfg.Clickhouse.ConnString)
	if err != nil {
		panic(err)
	}

	fmt.Println(db, conn)
}
