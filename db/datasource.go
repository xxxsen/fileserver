package db

import (
	"context"
	"database/sql"
)

type IDataSource interface {
	IQueryer
	IExecutor
}

type IQueryer interface {
	Query(sql string, args ...interface{}) (*sql.Rows, error)
	QueryContext(ctx context.Context, sql string, args ...interface{}) (*sql.Rows, error)
	QueryRow(sql string, args ...interface{}) *sql.Row
	QueryRowContext(ctx context.Context, sql string, args ...interface{}) *sql.Row
}

type IExecutor interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}
