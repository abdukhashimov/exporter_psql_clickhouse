package logrus

import (
	"github.com/abdukhashimov/exporter_psql_clickhouse/pkg/logger"
	"github.com/abdukhashimov/exporter_psql_clickhouse/pkg/logger/options"
)

// Factory is the receiver for logrus factory
type Factory struct{}

// Build logrus logger
func (_ *Factory) Build(cfg *options.Logging) (logger.Logger, error) {
	l, err := RegisterLogrusLog(cfg)
	if err != nil {
		return nil, err
	}

	return l, nil
}
