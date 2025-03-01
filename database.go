package funcy

import (
	"context"

	"github.com/jackc/pgx/v5"
)

// Query is a wrapper around the pgx.Query function.
func (cl *Client) Query(sql string, args ...interface{}) (pgx.Rows, error) {
	return cl.Conn.Query(context.Background(), sql, args...)
}

// QueryRow is a wrapper around the pgx.QueryRow function.
func (cl *Client) QueryRow(sql string, args ...interface{}) pgx.Row {
	return cl.Conn.QueryRow(context.Background(), sql, args...)
}
