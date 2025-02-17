package funcy

import (
	"context"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// Client is a wrapper around the S3 client and authentication tokens.
type Client struct {
	// Client for S3 operations.
	*s3.Client
	// W is the HTTP response writer.
	W http.ResponseWriter
	// R is the HTTP request.
	R *http.Request
	// Username is used to verify the token.
	Username string
	// Token is used to check if authenticated.
	Token string
}

// Write string to HTTP.
func (c *Client) Write(s string) {
	c.W.Write([]byte(s))
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

	c := &Client{}
	c.Client = s3.NewFromConfig(s3cfg)
	c.W = w
	c.R = r
	c.Username = r.URL.Query().Get("username")
	c.Token = r.URL.Query().Get("token")
	return c
}
