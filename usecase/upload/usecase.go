package upload

import (
	"api_blog/infrastructure/storage"
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

const (
	maxSize       = 1 << 20 // 1 MB
	presignExpiry = 5 * time.Minute
)

var allowedExt = map[string]string{
	".jpg":  "image/jpeg",
	".jpeg": "image/jpeg",
	".png":  "image/png",
	".webp": "image/webp",
}

type UploadUsecase struct {
	storage *storage.R2Storage
}

func NewUploadUsecase(s *storage.R2Storage) *UploadUsecase {
	return &UploadUsecase{storage: s}
}

func (u *UploadUsecase) PresignUpload(ctx context.Context, in PresignInput) (*PresignOutput, error) {
	if u.storage == nil {
		return nil, errors.New("storage is not configured")
	}

	ext := strings.ToLower(filepath.Ext(in.Filename))
	contentType, ok := allowedExt[ext]
	if !ok {
		return nil, errors.New("unsupported file format, allowed: jpg, jpeg, png, webp")
	}

	if in.Size <= 0 || in.Size > maxSize {
		return nil, errors.New("file size must be greater than 0 and at most 1MB")
	}

	name, err := randomName(ext)
	if err != nil {
		return nil, err
	}

	key := u.storage.ObjectKey(name)

	result, err := u.storage.PresignPut(ctx, key, contentType, in.Size, presignExpiry)
	if err != nil {
		return nil, err
	}

	return &PresignOutput{
		UploadUrl: result.URL,
		Method:    result.Method,
		Headers:   result.Headers,
		Key:       key,
		ImagePath: name,
		PublicUrl: u.storage.PublicURL(key),
		ExpiresIn: int(presignExpiry.Seconds()),
	}, nil
}

func randomName(ext string) (string, error) {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("failed to generate file name: %w", err)
	}
	return hex.EncodeToString(buf) + ext, nil
}
