package exporter

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

var _ TableExporter = (*Exporter)(nil)

type Exporter struct {
	dbConn         *sqlx.DB
	clickHouseConn *sqlx.DB
}

type TableExporter interface {
	ExportDataFromPsqlToClickhouse(tableName string) error
}

func New(dbConn *sqlx.DB, clickHouseConn *sqlx.DB) Exporter {
	return Exporter{dbConn: dbConn, clickHouseConn: clickHouseConn}
}

var (
	countQuery = "select count(1) from towns"
)

func (e *Exporter) ExportDataFromPsqlToClickhouse(tableName string) error {
	var (
		count    int
		pageSize = 100000
	)

	row := e.dbConn.QueryRow(countQuery)
	err := row.Scan(&count)
	if err != nil {
		return err
	}

	for offset := 0; offset < count; offset += pageSize {
		res, err := e.clickHouseConn.Exec("insert into sample (code, article, name, department) select code, article, name, department from postgresql('postgres-container:5432', 'sample', 'towns', 'postgres', 'postgres') LIMIT $1 OFFSET $2", pageSize, offset)
		if err != nil {
			return err
		}
		fmt.Println(res.RowsAffected())
	}

	return nil
}
