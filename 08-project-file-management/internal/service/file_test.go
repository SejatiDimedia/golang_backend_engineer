package service

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/url"
	"testing"
	"time"

	"github.com/timurdian/file-management/internal/entity"
)

// Mock FileRepository
type mockFileRepository struct {
	files  map[uint]*entity.File
	nextID uint

	onCreate func(file *entity.File) error
	onUpdate func(file *entity.File) error
	onDelete func(id uint) error
}

func (m *mockFileRepository) Create(ctx context.Context, file *entity.File) error {
	m.nextID++
	file.ID = m.nextID
	m.files[file.ID] = file
	if m.onCreate != nil {
		return m.onCreate(file)
	}
	return nil
}

func (m *mockFileRepository) Update(ctx context.Context, file *entity.File) error {
	m.files[file.ID] = file
	if m.onUpdate != nil {
		return m.onUpdate(file)
	}
	return nil
}

func (m *mockFileRepository) GetByID(ctx context.Context, id uint) (*entity.File, error) {
	file, exists := m.files[id]
	if !exists {
		return nil, nil
	}
	return file, nil
}

func (m *mockFileRepository) GetByUserID(ctx context.Context, userID uint) ([]entity.File, error) {
	var res []entity.File
	for _, f := range m.files {
		if f.UserID == userID && f.Status == "SUCCESS" {
			res = append(res, *f)
		}
	}
	return res, nil
}

func (m *mockFileRepository) Delete(ctx context.Context, id uint) error {
	delete(m.files, id)
	if m.onDelete != nil {
		return m.onDelete(id)
	}
	return nil
}

// Mock StorageClient
type mockStorageClient struct {
	uploads  map[string][]byte
	onUpload func(objectKey string, data []byte) error
}

func (m *mockStorageClient) UploadFile(ctx context.Context, bucketName, objectName string, reader io.Reader, objectSize int64, contentType string) error {
	data, _ := io.ReadAll(reader)
	m.uploads[objectName] = data
	if m.onUpload != nil {
		return m.onUpload(objectName, data)
	}
	return nil
}

func (m *mockStorageClient) GetPresignedURL(ctx context.Context, bucketName, objectName string, expiry time.Duration) (*url.URL, error) {
	return url.Parse("https://minio.mock/presigned/" + objectName)
}

func (m *mockStorageClient) GetFileStream(ctx context.Context, bucketName, objectName string) (io.ReadCloser, error) {
	data, exists := m.uploads[objectName]
	if !exists {
		return nil, errors.New("object not found")
	}
	return io.NopCloser(bytes.NewReader(data)), nil
}

func (m *mockStorageClient) DeleteFile(ctx context.Context, bucketName, objectName string) error {
	delete(m.uploads, objectName)
	return nil
}

func (m *mockStorageClient) Ping(ctx context.Context) error {
	return nil
}

func TestFileService_Upload_Success(t *testing.T) {
	fileRepo := &mockFileRepository{files: make(map[uint]*entity.File)}
	storage := &mockStorageClient{uploads: make(map[string][]byte)}
	svc := NewFileService(fileRepo, storage, "test-bucket", 15)

	content := []byte("fake image content")
	fileMeta, err := svc.Upload(
		context.Background(),
		1,
		"avatar.png",
		int64(len(content)),
		"image/png",
		bytes.NewReader(content),
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if fileMeta.Status != "SUCCESS" {
		t.Errorf("expected status SUCCESS, got %s", fileMeta.Status)
	}
	if len(storage.uploads) != 1 {
		t.Errorf("expected 1 physical file upload, got %d", len(storage.uploads))
	}
}

func TestFileService_Upload_ValidationLimits(t *testing.T) {
	fileRepo := &mockFileRepository{files: make(map[uint]*entity.File)}
	storage := &mockStorageClient{uploads: make(map[string][]byte)}
	svc := NewFileService(fileRepo, storage, "test-bucket", 15)

	// Test 1: File size limit (10MB + 1 byte)
	_, err := svc.Upload(
		context.Background(),
		1,
		"huge.pdf",
		10*1024*1024+1,
		"application/pdf",
		bytes.NewReader([]byte("")),
	)
	if !errors.Is(err, ErrFileTooLarge) {
		t.Errorf("expected ErrFileTooLarge, got %v", err)
	}

	// Test 2: Invalid content type (exe)
	_, err = svc.Upload(
		context.Background(),
		1,
		"malware.exe",
		100,
		"application/x-msdownload",
		bytes.NewReader([]byte("")),
	)
	if !errors.Is(err, ErrInvalidFileType) {
		t.Errorf("expected ErrInvalidFileType, got %v", err)
	}
}

func TestFileService_Upload_CompensatingRollback(t *testing.T) {
	fileRepo := &mockFileRepository{files: make(map[uint]*entity.File)}
	storage := &mockStorageClient{
		uploads: make(map[string][]byte),
		onUpload: func(objectKey string, data []byte) error {
			return errors.New("MinIO disk full")
		},
	}
	svc := NewFileService(fileRepo, storage, "test-bucket", 15)

	content := []byte("my doc content")
	_, err := svc.Upload(
		context.Background(),
		1,
		"resume.pdf",
		int64(len(content)),
		"application/pdf",
		bytes.NewReader(content),
	)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	// Verifikasi record PENDING dihapus dari database (compensating write)
	if len(fileRepo.files) != 0 {
		t.Errorf("expected metadata row to be rolled back/deleted, but found %d files in DB", len(fileRepo.files))
	}
}
