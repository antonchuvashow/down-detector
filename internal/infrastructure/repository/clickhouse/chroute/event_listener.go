package chroute

import (
	"context"
	"fmt"

	"github.com/ClickHouse/clickhouse-go/v2"
	"go.uber.org/zap"

	"detector/internal/route/application"
	route "detector/internal/route/domain"
)

type EventListener struct {
	eventCh <-chan routeapp.Event
	conn    clickhouse.Conn
	logger  *zap.Logger
}

func NewEventListener(conn clickhouse.Conn, logger *zap.Logger, eventCh <-chan routeapp.Event) *EventListener {
	return &EventListener{
		conn:    conn,
		logger:  logger,
		eventCh: eventCh,
	}
}

func (l *EventListener) Listen() {
	for event := range l.eventCh {
		var err error
		switch event.Type {
		case routeapp.EventTypeCreate:
			err = l.Add(event.Route)
		case routeapp.EventTypeUpdate:
			err = l.Update(event.Route)
		case routeapp.EventTypeDelete:
			err = l.Delete(event.RouteID)
		}
		if err != nil {
			l.logger.Error("error adding route", zap.Error(err))
		}
	}
}

func (l *EventListener) Add(rt *route.Route) error {
	const query = `
        INSERT INTO routes
            (id, url) 
        VALUES
            (?, ?)`

	// TODO: context is not passed!
	err := l.conn.Exec(context.Background(), query,
		rt.ID,
		rt.URL.String(),
	)
	if err != nil {
		l.logger.Error("failed to save route",
			zap.Error(err),
			zap.String("routeID", string(rt.ID)),
			zap.String("routeURL", rt.URL.String()),
		)
		return fmt.Errorf("save report: %w", err)
	}

	l.logger.Debug("route saved",
		zap.String("routeID", string(rt.ID)),
		zap.String("routeURL", rt.URL.String()),
	)

	return nil
}

func (l *EventListener) Update(rt *route.Route) error {
	const query = `UPDATE routes set url = ? WHERE id = ?`

	// TODO: context is not passed!
	err := l.conn.Exec(context.Background(), query,
		rt.URL.String(),
		rt.ID,
	)
	if err != nil {
		l.logger.Error("failed to update route",
			zap.Error(err),
			zap.String("routeID", string(rt.ID)),
			zap.String("routeURL", rt.URL.String()),
		)
		return fmt.Errorf("save report: %w", err)
	}

	l.logger.Debug("route updated",
		zap.String("routeID", string(rt.ID)),
		zap.String("routeURL", rt.URL.String()),
	)

	return nil
}

func (l *EventListener) Delete(rtID route.ID) error {
	const query = `DELETE FROM routes WHERE id = ?`

	// TODO: context is not passed!
	err := l.conn.Exec(context.Background(), query,
		rtID,
	)
	if err != nil {
		l.logger.Error("failed to delete route",
			zap.Error(err),
			zap.String("routeID", string(rtID)),
		)
		return fmt.Errorf("save report: %w", err)
	}

	l.logger.Debug("route deleted",
		zap.String("routeID", string(rtID)),
	)

	return nil
}
