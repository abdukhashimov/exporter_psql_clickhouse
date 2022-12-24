package zap

import (
	"github.com/abdukhashimov/exporter_psql_clickhouse/pkg/logger"
	"github.com/abdukhashimov/exporter_psql_clickhouse/pkg/logger/options"
)

// Factory is the receiver for zap factory
type Factory struct{}

// Build zap logger
func (_ *Factory) Build(cfg *options.Logging) (logger.Logger, error) {
	l, err := RegisterLog(cfg)
	if err != nil {
		return nil, err
	}

	return l, nil
}
