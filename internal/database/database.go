package database

import (
	"context"
	"github.com/jackc/pgx/v4"
)

type Database interface {
	QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row
	Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error)
}
