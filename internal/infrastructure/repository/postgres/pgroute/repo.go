package pgroute

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"time"

	"detector/internal/route/application"
	"detector/internal/route/domain"

	"github.com/jackc/pgx/v5/pgconn"

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

func (r *Repository) GetAllRoutes() ([]route.Route, error) {
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

	routes := make([]route.Route, 0)

	for rows.Next() {
		scannedRoute, err := scanRoute(rows)
		if err != nil {
			return nil, err
		}

		routes = append(routes, scannedRoute)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return routes, nil
}

func (r *Repository) Get(id route.ID) (route.Route, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	row := r.db.QueryRowContext(ctx, `
		SELECT id, url
		FROM routes
		WHERE id = $1
	`, string(id))

	scannedRoute, err := scanRoute(row)
	if errors.Is(err, sql.ErrNoRows) {
		return route.Route{}, routeapp.ErrNotFound{ID: id}
	}
	if err != nil {
		return route.Route{}, err
	}

	return scannedRoute, nil
}

func (r *Repository) Add(route route.Route) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	_, err := r.db.ExecContext(ctx, `
		INSERT INTO routes 
		    (id, url)
		VALUES 
		    ($1, $2)
	`, string(route.ID), route.URL.String())

	if isUniqueViolation(err) {
		return routeapp.ErrAlreadyExists{ID: route.ID}
	}
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) Update(route route.Route) error {
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
		return routeapp.ErrNotFound{ID: route.ID}
	}

	return nil
}

func (r *Repository) Delete(id route.ID) error {
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
		return routeapp.ErrNotFound{ID: id}
	}

	return nil
}

func (r *Repository) Search(path url.URL) ([]route.Route, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	rows, err := r.db.QueryContext(ctx, `
		SELECT id, url FROM routes
		WHERE url = $1`, path.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	routes := make([]route.Route, 0)
	for rows.Next() {
		scannedRoute, err := scanRoute(rows)
		if err != nil {
			return nil, err
		}
		routes = append(routes, scannedRoute)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return routes, nil
}

type routeScanner interface {
	Scan(dest ...any) error
}

func scanRoute(scanner routeScanner) (route.Route, error) {
	var (
		id     string
		rawURL string
	)

	if err := scanner.Scan(&id, &rawURL); err != nil {
		return route.Route{}, err
	}

	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return route.Route{}, fmt.Errorf("parse route url %q: %w", rawURL, err)
	}

	return route.Route{
		ID:  route.ID(id),
		URL: *parsedURL,
	}, nil
}

func isUniqueViolation(err error) bool {
	if err == nil {
		return false
	}

	if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
		return pgErr.Code == "23505"
	}

	return false
}
