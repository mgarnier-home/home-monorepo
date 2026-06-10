package s3

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// Config represents S3 configuration parameters.
type Config struct {
	Endpoint        string // Custom endpoint URL (e.g. SeaweedFS, MinIO)
	Region          string // AWS Region
	AccessKeyID     string // Access Key
	SecretAccessKey string // Secret Key
	Bucket          string // S3 Bucket name
}

// Client is a wrapper around the MinIO S3 client.
type Client struct {
	minioClient *minio.Client
	bucket      string
}

// ObjectInfo represents either a file or a virtual folder in the S3 bucket.
type ObjectInfo struct {
	Key          string    `json:"key"`          // Full S3 Key or prefix
	Name         string    `json:"name"`         // Base name of the file or folder
	Size         int64     `json:"size"`         // File size (0 for folders)
	LastModified time.Time `json:"lastModified"` // Last modified time (zero value for folders)
}

// NewClient initializes a new S3 Client using the MinIO Go SDK.
func NewClient(ctx context.Context, cfg Config) (*Client, error) {
	// Parse the Endpoint to extract host and secure flag
	endpoint := cfg.Endpoint
	useSSL := false

	if strings.HasPrefix(endpoint, "https://") {
		useSSL = true
		endpoint = strings.TrimPrefix(endpoint, "https://")
	} else if strings.HasPrefix(endpoint, "http://") {
		useSSL = false
		endpoint = strings.TrimPrefix(endpoint, "http://")
	}

	// Dynamic region defaults
	if cfg.Region == "" {
		cfg.Region = "us-east-1"
	}

	creds := credentials.NewStaticV4(cfg.AccessKeyID, cfg.SecretAccessKey, "")

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  creds,
		Secure: useSSL,
		Region: cfg.Region,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create MinIO client: %w", err)
	}

	return &Client{
		minioClient: minioClient,
		bucket:      cfg.Bucket,
	}, nil
}

func (c *Client) ListObjects(ctx context.Context, directory string) ([]ObjectInfo, error) {
	if c.bucket == "" {
		return nil, fmt.Errorf("bucket name cannot be empty")
	}

	opts := minio.ListObjectsOptions{
		Prefix:    directory,
		Recursive: true,
	}

	var results []ObjectInfo

	objectCh := c.minioClient.ListObjects(ctx, c.bucket, opts)
	for object := range objectCh {
		if object.Err != nil {
			return nil, fmt.Errorf("failed listing objects for directory %q: %w", directory, object.Err)
		}

		key := object.Key
		name := strings.Split(key, "/")[len(strings.Split(key, "/"))-1]
		isDir := strings.HasSuffix(key, "/")

		if isDir {
			continue
		} else {
			results = append(results, ObjectInfo{
				Key:          object.Key,
				Name:         name,
				Size:         object.Size,
				LastModified: object.LastModified,
			})
		}
	}

	return results, nil
}

// Download retrieves an object from S3 and writes it to the provided writer.
func (c *Client) Download(ctx context.Context, key string, writer io.Writer) error {
	if c.bucket == "" {
		return fmt.Errorf("bucket name cannot be empty")
	}

	object, err := c.minioClient.GetObject(ctx, c.bucket, key, minio.GetObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to get object %q preparation: %w", key, err)
	}
	defer object.Close()

	// Verify object exists and is accessible
	_, err = object.Stat()
	if err != nil {
		return fmt.Errorf("failed to retrieve object info for key %q: %w", key, err)
	}

	_, err = io.Copy(writer, object)
	if err != nil {
		return fmt.Errorf("failed to copy object data: %w", err)
	}

	return nil
}

// DownloadToFile retrieves an object from S3 and writes it to a local file.
// It automatically creates any parent directories if they don't exist.
func (c *Client) DownloadToFile(ctx context.Context, key string, localPath string) error {
	// Create all parent directories in the local path
	dir := filepath.Dir(localPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %q: %w", dir, err)
	}

	// Create/truncates the local file
	file, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("failed to create file %q: %w", localPath, err)
	}
	defer file.Close()

	return c.Download(ctx, key, file)
}
