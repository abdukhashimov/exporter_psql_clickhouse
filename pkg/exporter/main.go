package exporter

import (
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

	return nil
}
