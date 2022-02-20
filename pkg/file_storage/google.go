package file_storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/storage"
)

type GoogleCloudStorage struct {
	cl         *storage.Client
	projectID  string
	bucketName string
	uploadPath string
}

func NewGoogleCloudClient(projectID, bucketName, uploadPath string) (*GoogleCloudStorage, error) {
	client, err := storage.NewClient(context.Background())
	if err != nil {
		return nil, fmt.Errorf("NewGoogleCloudClient: %v", err)
	}
	return &GoogleCloudStorage{
		cl:         client,
		projectID:  projectID,
		bucketName: bucketName,
		uploadPath: uploadPath,
	}, nil
}

func (gcp GoogleCloudStorage) Write(ctx context.Context, filename string, data []byte) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	buf := bytes.NewBuffer(data)
	wc := gcp.cl.Bucket(gcp.bucketName).Object(gcp.uploadPath + filename).NewWriter(ctx)
	wc.ChunkSize = 1024 * 1000

	if _, err := io.Copy(wc, buf); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}
	return nil
}

func (gcp GoogleCloudStorage) Read(ctx context.Context, src string) ([]byte, error) {
	rc, err := gcp.cl.Bucket(gcp.bucketName).Object(gcp.uploadPath + src).NewReader(ctx)
	if err != nil {
		return nil, fmt.Errorf("Read: %v", err)
	}
	defer rc.Close()
	slurp, err := io.ReadAll(rc)
	if err != nil {
		return nil, fmt.Errorf("io.ReadAll: %v", err)
	}
	return slurp, nil
}
