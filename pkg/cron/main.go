package cron

import (
	"github.com/abdukhashimov/exporter_psql_clickhouse/pkg/exporter"
	"github.com/abdukhashimov/exporter_psql_clickhouse/pkg/logger"
	"github.com/robfig/cron/v3"
)

type CronJob struct {
	cronJob  *cron.Cron
	exporter exporter.Exporter
}

func New(exporter exporter.Exporter) *CronJob {
	return &CronJob{exporter: exporter, cronJob: cron.New()}
}

func (c *CronJob) RunTableExporter(cronPeriod string, tableName string) error {
	_, err := c.cronJob.AddFunc(cronPeriod, func() {
		err := c.exporter.Export(tableName)
		logger.Log.Error("failed to run the given function", err)
	})
	if err != nil {
		return err
	}

	return nil
}
