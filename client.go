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
	// Username is used to verify the token.
	Username string
	// Token is used to check if authenticated.
	Token string
}

// Write string to HTTP.
func (cl *Client) Write(s string) {
	cl.W.Write([]byte(s))
}

// NewClient creates a client.
func NewClient(w http.ResponseWriter, r *http.Request) *Client {
	s3cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(os.Getenv("REGION")),
	)
	if err != nil {
		return nil
	}

	cl := &Client{}
	cl.Client = s3.NewFromConfig(s3cfg)
	cl.W = w
	cl.R = r
	cl.Username = r.URL.Query().Get("username")
	cl.Token = r.URL.Query().Get("token")

	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE"))
	if err != nil {
		ll.Msg("Error connecting to database: %s", err.Error())
		return nil
	}

	cl.Conn = conn
	return cl
}

// SetCookie in HTTP response.
func (cl *Client) SetCookie(name, value string) {
	http.SetCookie(cl.W, &http.Cookie{
		Name:     name,
		Value:    value,
		MaxAge:   3600,
		Path:     "/",
		Domain:   os.Getenv("DOMAIN"),
		SameSite: http.SameSiteLaxMode,
		HttpOnly: true,
		Secure:   true,
	})
}

// ClearCookie in HTTP response.
func (cl *Client) ClearCookie(name string) {
	http.SetCookie(cl.W, &http.Cookie{
		Name:     name,
		Value:    "",
		MaxAge:   -1,
		Path:     "/",
		Domain:   os.Getenv("DOMAIN"),
		SameSite: http.SameSiteLaxMode,
		HttpOnly: true,
		Secure:   true,
	})
}

// Flush the HTTP response.
func (cl *Client) Flush() {
	if flusher, ok := cl.W.(http.Flusher); ok {
		flusher.Flush()
	}
}
