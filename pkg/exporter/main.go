package exporter

import "github.com/jmoiron/sqlx"

type Exporter interface {
	Export(tableName string) error
}

type Export struct {
	psqlConn   *sqlx.DB
	cHouseConn *sqlx.DB
}

func New(psqlConn, cHouseConn *sqlx.DB) Exporter {
	return &Export{
		psqlConn:   psqlConn,
		cHouseConn: cHouseConn,
	}
}

func (e *Export) Export(tableName string) error {
	return nil
}
