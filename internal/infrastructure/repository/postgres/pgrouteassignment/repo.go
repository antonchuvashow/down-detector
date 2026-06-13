package pgrouteassignment

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"detector/internal/inspector/application"
	"detector/internal/inspector/application/dto"
	"detector/internal/inspector/domain"
	"detector/internal/route/domain"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Repository struct {
	db      *sql.DB
	timeout time.Duration
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db:      db,
		timeout: 5 * time.Second,
	}
}

func (r *Repository) Save(routeAssignment inspectordto.RouteAssignment) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	_, err := r.db.ExecContext(ctx, `
		INSERT INTO route_methods 
		    (route_id, factory_key, serialized_config)
		VALUES 
		    ($1, $2, $3)
		ON CONFLICT (route_id) DO UPDATE
		SET
		    factory_key = EXCLUDED.factory_key,
		    serialized_config = EXCLUDED.serialized_config,
		    updated_at = now()
	`, string(routeAssignment.RouteID), string(routeAssignment.FactoryKey), []byte(routeAssignment.SerializedConfig))

	if err != nil {
		return fmt.Errorf("save route method: %w", err)
	}

	return nil
}

func (r *Repository) Get(routeID route.ID) (*inspectordto.RouteAssignment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	row := r.db.QueryRowContext(ctx, `
		SELECT route_id, factory_key, serialized_config
		FROM route_methods
		WHERE route_id = $1
	`, string(routeID))

	routeAssignment, err := scanRouteAssignment(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, &inspectorapp.ErrNotFound{RouteID: routeID}
	}
	if err != nil {
		return nil, err
	}

	return &routeAssignment, nil
}

func (r *Repository) Delete(routeID route.ID) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	row := r.db.QueryRowContext(ctx, `SELECT exists(SELECT 1 FROM route_methods WHERE route_id = $1)`, string(routeID))
	var exists bool
	err := row.Scan(&exists)
	if err != nil {
		return err
	}

	if !exists {
		return &inspectorapp.ErrNotFound{RouteID: routeID}
	}

	_, err = r.db.ExecContext(ctx, `DELETE FROM route_methods WHERE route_id = $1`, string(routeID))
	if err != nil {
		return err
	}

	return nil
}

type routeAssignmentScanner interface {
	Scan(dest ...any) error
}

func scanRouteAssignment(scanner routeAssignmentScanner) (inspectordto.RouteAssignment, error) {
	var (
		routeID          string
		factoryKey       string
		serializedConfig []byte
	)

	if err := scanner.Scan(&routeID, &factoryKey, &serializedConfig); err != nil {
		return inspectordto.RouteAssignment{}, err
	}

	return inspectordto.RouteAssignment{
		RouteID:          route.ID(routeID),
		FactoryKey:       inspector.FactoryKey(factoryKey),
		SerializedConfig: serializedConfig,
	}, nil
}
