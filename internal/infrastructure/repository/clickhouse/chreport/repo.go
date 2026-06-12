package chreport

import (
	"context"
	"fmt"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"go.uber.org/zap"

	"detector/internal/report/domain"
)

type Repository struct {
	conn   clickhouse.Conn
	logger *zap.Logger
}

func NewRepository(conn clickhouse.Conn, logger *zap.Logger) *Repository {
	return &Repository{
		conn:   conn,
		logger: logger.With(zap.String("component", "report_repository")),
	}
}

func (r *Repository) Save(report report.Report) error {
	const query = `
        INSERT INTO reports
            (time, route_id, success, error_types, latency_ms, source, latitude, longitude, ip, platform)
        VALUES
            (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	// TODO: context is not passed!
	err := r.conn.Exec(context.Background(), query,
		report.Time.UTC().Truncate(time.Millisecond),
		report.RouteID,
		report.Success,
		errorTypesToStrings(report.ErrorTypes),
		report.Summary.Latency.Milliseconds(),
		string(report.Descriptor.Source),
		report.Descriptor.Latitude,
		report.Descriptor.Longitude,
		report.Descriptor.IP.String(),
		report.Descriptor.Platform,
	)
	if err != nil {
		r.logger.Error("failed to save report",
			zap.Error(err),
			zap.String("source", string(report.Descriptor.Source)),
			zap.Time("time", report.Time),
		)
		return fmt.Errorf("save report: %w", err)
	}

	r.logger.Debug("report saved",
		zap.String("source", string(report.Descriptor.Source)),
		zap.Bool("success", report.Success),
		zap.Duration("latency", report.Summary.Latency),
	)

	return nil
}

func (r *Repository) SaveBatch(ctx context.Context, reports []report.Report) error {
	if len(reports) == 0 {
		return nil
	}

	batch, err := r.conn.PrepareBatch(ctx, `
        INSERT INTO reports
            (time, route_id, success, error_types, latency_ms, source, latitude, longitude, ip, platform)`)
	if err != nil {
		return fmt.Errorf("prepare batch: %w", err)
	}

	for _, report := range reports {
		if err = batch.Append(
			report.Time.UTC().Truncate(time.Millisecond),
			report.RouteID,
			report.Success,
			errorTypesToStrings(report.ErrorTypes),
			report.Summary.Latency.Milliseconds(),
			string(report.Descriptor.Source),
			report.Descriptor.Latitude,
			report.Descriptor.Longitude,
			report.Descriptor.IP.String(),
			report.Descriptor.Platform,
		); err != nil {
			return fmt.Errorf("append to batch: %w", err)
		}
	}

	if err = batch.Send(); err != nil {
		r.logger.Error("failed to send batch",
			zap.Error(err),
			zap.Int("batch_size", len(reports)),
		)
		return fmt.Errorf("send batch: %w", err)
	}

	r.logger.Debug("batch saved", zap.Int("batch_size", len(reports)))
	return nil
}

// --- helpers ---

func errorTypesToStrings(errorTypes map[report.ErrorType]struct{}) []string {
	if len(errorTypes) == 0 {
		return []string{}
	}

	result := make([]string, 0, len(errorTypes))
	for et := range errorTypes {
		result = append(result, string(et))
	}
	return result
}
