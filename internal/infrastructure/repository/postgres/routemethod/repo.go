package routemethod

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"detector/internal/inspection/application"
	inspectiondto "detector/internal/inspection/application/dto"
	"detector/internal/inspection/domain/inspector"
	routedomain "detector/internal/route/domain"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgresRouteMethodRepository struct {
	db      *sql.DB
	timeout time.Duration
}

func NewPostgresRouteMethodRepository(db *sql.DB) *PostgresRouteMethodRepository {
	return &PostgresRouteMethodRepository{
		db:      db,
		timeout: 5 * time.Second,
	}
}

func (r *PostgresRouteMethodRepository) SaveRouteMethod(routeMethod inspectiondto.RouteMethod) error {
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
	`, string(routeMethod.RouteID), string(routeMethod.FactoryKey), []byte(routeMethod.SerializedConfig))

	if err != nil {
		return fmt.Errorf("save route method: %w", err)
	}

	return nil
}

func (r *PostgresRouteMethodRepository) GetRouteMethod(routeID routedomain.RouteID) (*inspectiondto.RouteMethod, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	row := r.db.QueryRowContext(ctx, `
		SELECT route_id, factory_key, serialized_config
		FROM route_methods
		WHERE route_id = $1
	`, string(routeID))

	routeMethod, err := scanRouteMethod(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, &application.ErrInspectorNotFound{RouteID: routeID}
	}
	if err != nil {
		return nil, err
	}

	return &routeMethod, nil
}

type routeMethodScanner interface {
	Scan(dest ...any) error
}

func scanRouteMethod(scanner routeMethodScanner) (inspectiondto.RouteMethod, error) {
	var (
		routeID          string
		factoryKey       string
		serializedConfig []byte
	)

	if err := scanner.Scan(&routeID, &factoryKey, &serializedConfig); err != nil {
		return inspectiondto.RouteMethod{}, err
	}

	return inspectiondto.RouteMethod{
		RouteID:          routedomain.RouteID(routeID),
		FactoryKey:       inspector.FactoryKey(factoryKey),
		SerializedConfig: inspector.SerializedConfig(serializedConfig),
	}, nil
}
