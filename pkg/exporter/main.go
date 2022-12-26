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
	transferRowCount        = 10000
	contextDeadlineDuration = 7
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
	countTransactions      = "select count(1) from %s;"
	selectListofIds        = "select code from %s WHERE soft_delete=false limit $1 offset $2"
	transferDataQuery      = "insert into towns (code, article, name, department) select code, article, name, department from postgresql('postgres-container:5432', 'sample', 'towns', 'postgres', 'postgres') LIMIT $1 OFFSET $2"
	updateManyTransactions = "update %s set soft_delete = true where code in (?);"
)

func (e *Export) Export(tableName string) error {
	var (
		rowCount int
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

		rows, err := e.psqlConn.Query(addTableName(selectListofIds, tableName), transferRowCount, row)
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

		_, err = e.cHouseConn.Exec(
			transferDataQuery,
			transferRowCount,
			row,
		)
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

		logger.Log.Infof("successfully transferred from [%d - %d)", row, row+transferRowCount)

		message := tgbotapi.NewMessage(e.cfg.Exporter.TelegramBotChannelID, fmt.Sprintf("%d/%d", row, rowCount))
		_, err = e.tbBot.Send(message)
		if err != nil {
			logger.Log.Error("failed to publish the message to telegram chat", err)
		}
	}

	return nil
}

func addTableName(query string, tableName string) string {
	return fmt.Sprintf(query, tableName)
}
