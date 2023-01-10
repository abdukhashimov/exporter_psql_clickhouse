package exporter

import (
	"fmt"
	"math"

	"github.com/abdukhashimov/exporter_psql_clickhouse/config"
	"github.com/abdukhashimov/exporter_psql_clickhouse/pkg/logger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/jmoiron/sqlx"
)

const (
	transferRowCount        = 40000
	contextDeadlineDuration = 7
	psqlUpdateCount         = 40000
)

type Exporter interface {
	Export(tableName string) error
}

type Export struct {
	psqlConn   *sqlx.DB
	cHouseConn *sqlx.DB
	cfg        *config.Config
	tbBot      *tgbotapi.BotAPI
}

func New(psqlConn, cHouseConn *sqlx.DB, cfg *config.Config, bot *tgbotapi.BotAPI) Exporter {
	return &Export{
		psqlConn:   psqlConn,
		cHouseConn: cHouseConn,
		cfg:        cfg,
		tbBot:      bot,
	}
}

var (
	countTransactions = "select count(1) from %s where deleted_at is null;"
	selectListofIds   = "select id from %s WHERE deleted_at is null limit $1"
	transferDataQuery = `insert into 
		%s (id, user_id, balls, level_id, step, deleted_at, updated_at, created_at)
		select id, user_id, balls, level_id, step, deleted_at, updated_at, created_at
		from postgresql('%s', %s, %s, %s, %s)
		WHERE deleted_at is null LIMIT $1`
	updateManyTransactions = "update %s set deleted_at = now() where id in (?);"
)

func (e *Export) Export(tableName string) error {
	var (
		rowCount int
	)

	transferDataQuery = fmt.Sprintf(
		transferDataQuery,
		tableName,
		e.cfg.Network.PsqlAddress,
		e.cfg.PsqlConfig.Database,
		tableName,
		e.cfg.PsqlConfig.User,
		e.cfg.PsqlConfig.Passwrod,
	)

	logger.Log.Info("exporter started")

	countRow := e.psqlConn.QueryRow(addTableName(countTransactions, tableName))
	err := countRow.Scan(&rowCount)
	if err != nil {
		return err
	}

	rowCountCeil := int(math.Ceil(float64(rowCount)/transferRowCount)) * transferRowCount

	for row := 0; row < rowCountCeil; row += transferRowCount {
		var (
			ids = []string{}
			id  string
		)

		rows, err := e.psqlConn.Query(addTableName(selectListofIds, tableName), transferRowCount)
		if err != nil {
			return err
		}

		for rows.Next() {
			err := rows.Scan(&id)
			if err != nil {
				return err
			}

			ids = append(ids, id)
		}

		_, err = e.cHouseConn.Exec(transferDataQuery, transferRowCount)
		if err != nil {
			return err
		}

		qry, args, err := sqlx.In(addTableName(updateManyTransactions, tableName), ids)
		if err != nil {
			return err
		}

		if _, err = e.psqlConn.Exec(e.psqlConn.Rebind(qry), args...); err != nil {
			return err
		}

		logger.Log.Infof("successfully transferred from [%d - %d)", row, len(ids)+row)
	}

	return nil
}

func chunkBy[T any](items []T, chunkSize int) (chunks [][]T) {
	for chunkSize < len(items) {
		items, chunks = items[chunkSize:], append(chunks, items[0:chunkSize:chunkSize])
	}
	return append(chunks, items)
}

func addTableName(query string, tableName string) string {
	return fmt.Sprintf(query, tableName)
}
