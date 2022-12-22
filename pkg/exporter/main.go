package exporter

import "github.com/jmoiron/sqlx"

var _ TableExporter = (*Exporter)(nil)

type Exporter struct {
	dbConn *sqlx.Conn
}

type TableExporter interface {
	ExportDataFromPsqlToClickhouse(tableName string) error
}

func New(DbConn *sqlx.Conn) Exporter {
	return Exporter{dbConn: DbConn}
}

func (e *Exporter) ExportDataFromPsqlToClickhouse(tableName string) error {

	return nil
}
