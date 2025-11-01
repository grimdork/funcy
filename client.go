package funcy

import (
	"context"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	ll "github.com/grimdork/loglines"
	"github.com/jackc/pgx/v5"
)

// Client is a wrapper around the S3 client and authentication tokens.
type Client struct {
	// Client for S3 operations.
	*s3.Client
	// W is the HTTP response writer.
	W http.ResponseWriter
	// R is the HTTP request.
	R *http.Request
	// Conn is the PostgreSQL connection.
	Conn *pgx.Conn
}

// Write string to HTTP.
func (cl *Client) Write(s string) {
	cl.W.Write([]byte(s))
}

// NewClient creates a client.
func NewClient(w http.ResponseWriter, r *http.Request) *Client {
	cl := &Client{
		W: w,
		R: r,
	}
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE"))
	if err != nil {
		ll.Msg("Error connecting to database: %s", err.Error())
		// We'll allow continuing without a database connection.
	}

	cl.Conn = conn
	s3cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(os.Getenv("REGION")),
	)
	if err != nil {
		return nil
	}

	cl.Client = s3.NewFromConfig(s3cfg)
	return cl
}
