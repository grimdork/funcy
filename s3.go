package funcy

import (
	"bytes"
	"context"
	"io"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// GetFile from storage.
func (c *Client) GetFile(fn string) []byte {
	res, err := c.GetObject(context.Background(), &s3.GetObjectInput{
		Bucket: aws.String(os.Getenv("BUCKET")),
		Key:    aws.String(fn),
	})

	if err != nil {
		c.Write(err.Error() + "\n")
		return nil
	}

	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		c.Write(err.Error() + "\n")
		return nil
	}

	return data
}

// PutFile to storage.
func (c *Client) PutFile(fn string, data []byte) {
	_, err := c.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String(os.Getenv("BUCKET")),
		Key:    aws.String(fn),
		Body:   bytes.NewReader(data),
	})

	if err != nil {
		c.Write(err.Error() + "\n")
		return
	}
}
