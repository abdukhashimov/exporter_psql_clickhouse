package exporter

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/abdukhashimov/exporter_psql_clickhouse/config"
	"github.com/abdukhashimov/exporter_psql_clickhouse/pkg/logger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/jmoiron/sqlx"
)

const (
	transferRowCount        = 100000
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
	selectListofIds        = "select code from %s limit $1 offset $2;"
	transferDataQuery      = "insert into %s (code, article, name, department) select code, article, name, department from postgresql('$1', '$2', '$3', '$4', '$5') LIMIT $6 OFFSET $7;"
	updateManyTransactions = "update %s set soft_delete = true where code in (?)"
)

func (e *Export) Export(tableName string) error {
	var (
		rowCount            int
		successfullRowCount int
	)

	logger.Log.Info("exporter started")
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(7)*time.Second)
	defer cancel()

	countRow := e.psqlConn.QueryRow(addTableName(countTransactions, tableName))
	err := countRow.Scan(&rowCount)
	if err != nil {
		return err
	}

	rowCountCeil := int(math.Ceil(float64(rowCount) / transferRowCount))

	for row := 0; row < rowCountCeil; row += transferRowCount {
		var (
			ids = []string{}
			id  string
		)

		cHTx, err := e.cHouseConn.BeginTx(ctx, nil)
		if err != nil {
			return err
		}

		tx, err := e.psqlConn.BeginTx(ctx, nil)
		if err != nil {
			return err
		}

		rows, err := tx.Query(addTableName(selectListofIds, tableName), transferRowCount, row)
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

		_, err = cHTx.ExecContext(
			ctx,
			addTableName(transferDataQuery, tableName),
			e.cfg.Network.PsqlAddress,
			e.cfg.PsqlConfig.Database,
			tableName,
			e.cfg.PsqlConfig.User,
			e.cfg.PsqlConfig.Passwrod,
		)
		if err != nil {
			return err
		}

		qry, args, err := sqlx.In(addTableName(updateManyTransactions, tableName), ids)
		if err != nil {
			return err
		}

		if _, err = tx.Exec(qry, args); err != nil {
			err := tx.Rollback()
			if err != nil {
				cHTx.Rollback()
				return err
			}

			err = cHTx.Rollback()
			if err != nil {
				return err
			}
		}

		logger.Log.Infof("successfully transferred from [%d - %d)", row, row+transferRowCount)
		successfullRowCount += row

		message := tgbotapi.NewMessage(e.cfg.Exporter.TelegramBotChannelID, fmt.Sprintf("%d/%d", successfullRowCount, rowCount))
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
