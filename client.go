package funcy

import (
	"context"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gorilla/sessions"
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
	// Store is the session store.
	Store *sessions.CookieStore
	// Session for this client.
	Session *sessions.Session
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
		W:     w,
		R:     r,
		Store: sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY"))),
	}
	cl.Store.Options = &sessions.Options{
		Path:     "/",
		Domain:   os.Getenv("DOMAIN"),
		MaxAge:   86400,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	session, err := cl.Store.Get(cl.R, "grimdork-session")
	if err != nil {
		ll.Msg("Error getting session: %s", err.Error())
		return nil
	}
	cl.Session = session
	err = cl.Save()
	if err != nil {
		ll.Msg("Error saving session: %s", err.Error())
		return nil
	}

	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE"))
	if err != nil {
		ll.Msg("Error connecting to database: %s", err.Error())
		return nil
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
