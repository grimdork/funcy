package funcy

import (
	"bytes"
	"context"
	"io"
	"os"

	ll "github.com/grimdork/loglines"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// GetFile from storage.
func (cl *Client) GetFile(fn string) []byte {
	res, err := cl.GetObject(context.Background(), &s3.GetObjectInput{
		Bucket: aws.String(os.Getenv("BUCKET")),
		Key:    aws.String(fn),
	})

	if err != nil {
		ll.Err("Error loading '%s': %s", fn, err.Error())
		return nil
	}

	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		ll.Err("Error loading '%s': %s", fn, err.Error())
		return nil
	}

	return data
}

// PutFile to storage.
func (cl *Client) PutFile(fn string, data []byte) {
	_, err := cl.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String(os.Getenv("BUCKET")),
		Key:    aws.String(fn),
		Body:   bytes.NewReader(data),
	})

	if err != nil {
		ll.Err("Error saving '%s': %s", fn, err.Error())
		return
	}
}
