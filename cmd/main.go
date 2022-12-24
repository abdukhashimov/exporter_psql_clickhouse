package main

import (
	"fmt"

	_ "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/abdukhashimov/exporter_psql_clickhouse/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	cfg := config.Load()

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
