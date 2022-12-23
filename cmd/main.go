package main

import (
	"fmt"

	_ "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/abdukhashimov/integration/config"
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
	fmt.Println(db)
	err = ConnectDSN()
	if err != nil {
		fmt.Println(err.Error())
	}
}

func ConnectDSN() error {
	conn, err := sqlx.Open("clickhouse", fmt.Sprintf("clickhouse://%s:%d?username=%s&password=%s", "localhost", 9000, "default", ""))
	if err != nil {
		return err
	}

	rows, err := conn.Query("select * from towns")
	if err != nil {
		return err
	}

	for rows.Next() {
		fmt.Println(rows)
	}

	return conn.Ping()
}
