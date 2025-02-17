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
	*s3.Client
	w        http.ResponseWriter
	r        *http.Request
	username string
	token    string
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

	c := &Client{s3.NewFromConfig(s3cfg), w, r,
		r.URL.Query().Get("username"),
		r.URL.Query().Get("token"),
	}
	return c
}
