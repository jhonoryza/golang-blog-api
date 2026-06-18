package storage

import (
	"context"
	"errors"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type R2Storage struct {
	client     *s3.Client
	presign    *s3.PresignClient
	bucket     string
	publicURL  string
	pathPrefix string
}

type PresignResult struct {
	URL     string
	Method  string
	Headers map[string]string
}

// PublicBaseURL returns R2_PUBLIC_URL without a trailing slash.
func PublicBaseURL() string {
	return strings.TrimRight(os.Getenv("R2_PUBLIC_URL"), "/")
}

// PathPrefix returns R2_PATH_PREFIX without surrounding slashes.
func PathPrefix() string {
	return strings.Trim(os.Getenv("R2_PATH_PREFIX"), "/")
}

// ImageURL builds the full public URL for a stored image path. The path is what
// is persisted in the database (the portion after the configured prefix).
func ImageURL(path string) string {
	if path == "" {
		return ""
	}
	base := PublicBaseURL()
	prefix := PathPrefix()
	path = strings.TrimLeft(path, "/")
	if prefix != "" {
		return base + "/" + prefix + "/" + path
	}
	return base + "/" + path
}

func NewR2Storage() (*R2Storage, error) {
	accessKey := os.Getenv("R2_ACCESS_KEY_ID")
	secretKey := os.Getenv("R2_SECRET_ACCESS_KEY")
	bucket := os.Getenv("R2_BUCKET")
	endpoint := os.Getenv("R2_ENDPOINT")
	region := os.Getenv("R2_REGION")

	if accessKey == "" || secretKey == "" || bucket == "" || endpoint == "" {
		return nil, errors.New("R2 storage is not configured (missing R2_ACCESS_KEY_ID/R2_SECRET_ACCESS_KEY/R2_BUCKET/R2_ENDPOINT)")
	}
	if region == "" {
		region = "auto"
	}

	client := s3.New(s3.Options{
		Region:       region,
		Credentials:  credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""),
		BaseEndpoint: aws.String(endpoint),
		UsePathStyle: true,
	})

	return &R2Storage{
		client:     client,
		presign:    s3.NewPresignClient(client),
		bucket:     bucket,
		publicURL:  PublicBaseURL(),
		pathPrefix: PathPrefix(),
	}, nil
}

// ObjectKey returns the full R2 object key for a stored file name, applying the
// configured path prefix.
func (s *R2Storage) ObjectKey(name string) string {
	if s.pathPrefix == "" {
		return name
	}
	return s.pathPrefix + "/" + name
}

func (s *R2Storage) PresignPut(ctx context.Context, key, contentType string, size int64, expires time.Duration) (*PresignResult, error) {
	req, err := s.presign.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(s.bucket),
		Key:           aws.String(key),
		ContentType:   aws.String(contentType),
		ContentLength: aws.Int64(size),
	}, s3.WithPresignExpires(expires))
	if err != nil {
		return nil, err
	}

	headers := map[string]string{}
	for name := range req.SignedHeader {
		if strings.EqualFold(name, "host") {
			continue
		}
		headers[name] = req.SignedHeader.Get(name)
	}

	return &PresignResult{
		URL:     req.URL,
		Method:  req.Method,
		Headers: headers,
	}, nil
}

func (s *R2Storage) PublicURL(key string) string {
	return s.publicURL + "/" + strings.TrimLeft(key, "/")
}
