package route

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"time"

	routeapplication "detector/internal/route/application"
	routedomain "detector/internal/route/domain"

	"github.com/jackc/pgx/v5/pgconn"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Repository struct {
	db      *sql.DB
	timeout time.Duration
}

func NewPostgresRouteRepository(db *sql.DB) *Repository {
	return &Repository{
		db:      db,
		timeout: 5 * time.Second,
	}
}

func (r *Repository) GetAllRoutes() ([]routedomain.Route, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	rows, err := r.db.QueryContext(ctx, `
		SELECT id, url
		FROM routes
		ORDER BY id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	routes := make([]routedomain.Route, 0)

	for rows.Next() {
		route, err := scanRoute(rows)
		if err != nil {
			return nil, err
		}

		routes = append(routes, route)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return routes, nil
}

func (r *Repository) Get(id routedomain.RouteID) (routedomain.Route, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	row := r.db.QueryRowContext(ctx, `
		SELECT id, url
		FROM routes
		WHERE id = $1
	`, string(id))

	route, err := scanRoute(row)
	if errors.Is(err, sql.ErrNoRows) {
		return routedomain.Route{}, routeapplication.ErrRouteNotFound{ID: id}
	}
	if err != nil {
		return routedomain.Route{}, err
	}

	return route, nil
}

func (r *Repository) Add(route routedomain.Route) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	_, err := r.db.ExecContext(ctx, `
		INSERT INTO routes 
		    (id, url)
		VALUES 
		    ($1, $2)
	`, string(route.ID), route.URL.String())

	if isUniqueViolation(err) {
		return routeapplication.ErrRouteAlreadyExists{ID: route.ID}
	}
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) Update(route routedomain.Route) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	result, err := r.db.ExecContext(ctx, `
		UPDATE routes
		SET 
		    url = $2,
		    updated_at = now()
		WHERE id = $1
	`, string(route.ID), route.URL.String())

	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return routeapplication.ErrRouteNotFound{ID: route.ID}
	}

	return nil
}

func (r *Repository) Delete(id routedomain.RouteID) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	result, err := r.db.ExecContext(ctx, `
		DELETE FROM routes
		WHERE id = $1
	`, string(id))

	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return routeapplication.ErrRouteNotFound{ID: id}
	}

	return nil
}

func (r *Repository) Search(path url.URL) ([]routedomain.Route, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	rows, err := r.db.QueryContext(ctx, `
		SELECT id, url FROM routes
		WHERE url = $1`, path.String())
	if err != nil {
		return nil, err
	}

	routes := make([]routedomain.Route, 0)
	for rows.Next() {
		route, err := scanRoute(rows)
		if err != nil {
			return nil, err
		}
		routes = append(routes, route)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return routes, nil
}

type routeScanner interface {
	Scan(dest ...any) error
}

func scanRoute(scanner routeScanner) (routedomain.Route, error) {
	var (
		id     string
		rawURL string
	)

	if err := scanner.Scan(&id, &rawURL); err != nil {
		return routedomain.Route{}, err
	}

	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return routedomain.Route{}, fmt.Errorf("parse route url %q: %w", rawURL, err)
	}

	return routedomain.Route{
		ID:  routedomain.RouteID(id),
		URL: *parsedURL,
	}, nil
}

func isUniqueViolation(err error) bool {
	if err == nil {
		return false
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}

	return false
}
