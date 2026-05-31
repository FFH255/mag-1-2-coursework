package postgres

import (
	"context"
	"errors"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Config interface {
	GetConnString() string
}

type Database struct {
	pool *pgxpool.Pool
}

// MustConnect initializes and connects to the database using a connection pool.
func (d *Database) MustConnect(ctx context.Context, connectionString string) {
	var err error
	d.pool, err = pgxpool.New(ctx, connectionString)
	if err != nil {
		panic("failed to create pgxpool: " + err.Error())
	}
}

// MustClose closes the connection pool.
func (d *Database) MustClose() {
	if d.pool == nil {
		return
	}
	d.pool.Close()
}

// QueryRow executes a query that returns a single row.
func (d *Database) QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row {
	return d.pool.QueryRow(ctx, query, args...)
}

// Query executes a query that returns multiple rows.
func (d *Database) Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error) {
	return d.pool.Query(ctx, query, args...)
}

// Exec executes a query without returning rows.
func (d *Database) Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error) {
	return d.pool.Exec(ctx, query, args...)
}

// MustMigrate runs database migrations
func (d *Database) MustMigrate(sourceURL string, connString string) {
	m, err := migrate.New(sourceURL, connString)
	if err != nil {
		panic("failed to create migrate instance: " + err.Error())
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		panic("migration failed: " + err.Error())
	}
}

// New creates a new Database instance.
func New() *Database {
	return &Database{}
}
