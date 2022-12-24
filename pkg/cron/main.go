package cron

import (
	"github.com/abdukhashimov/exporter_psql_clickhouse/pkg/exporter"
	"github.com/robfig/cron/v3"
)

type CronJob struct {
	cronJob  *cron.Cron
	exporter exporter.Exporter
}

func New(exporter exporter.Exporter) *CronJob {
	return &CronJob{exporter: exporter, cronJob: cron.New()}
}

func (c *CronJob) RunTableExporter(cronPeriod string) error {
	_, err := c.cronJob.AddFunc(cronPeriod, c.runPsqlTableExporterToClickHouse)
	if err != nil {
		return err
	}

	return nil
}

func (c *CronJob) runPsqlTableExporterToClickHouse() {
	// TODO: execute exporter
}
