package main

import (
	"fmt"

	_ "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/abdukhashimov/exporter_psql_clickhouse/config"
	"github.com/abdukhashimov/exporter_psql_clickhouse/pkg/exporter"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	cfg := config.Load()
	psqlConnString := cfg.MakePSQLConnString()
	db, err := sqlx.Connect("postgres", psqlConnString)
	if err != nil {
		panic(err)
	}

	conn, err := sqlx.Open("clickhouse", fmt.Sprintf("clickhouse://%s:%d?username=%s&password=%s", "localhost", 9000, "default", ""))
	if err != nil {
		panic(err)
	}

	exp := exporter.New(db, conn)
	err = exp.ExportDataFromPsqlToClickhouse("towns")
	if err != nil {
		fmt.Println(err.Error())
	}
}
