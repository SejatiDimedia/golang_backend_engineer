package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/timurdian/file-management/internal/entity"
	"github.com/timurdian/file-management/internal/repository"
	"github.com/timurdian/file-management/internal/utils"
)

var (
	ErrFileTooLarge     = errors.New("file size exceeds the maximum limit of 10MB")
	ErrInvalidFileType  = errors.New("invalid file type, only JPEG, PNG, and PDF are allowed")
	ErrFileNotFound     = errors.New("file metadata not found")
	ErrPermissionDenied = errors.New("you do not have permission to access this file")
)

type FileService interface {
	Upload(ctx context.Context, userID uint, filename string, size int64, contentType string, reader io.Reader) (*entity.File, error)
	GetList(ctx context.Context, userID uint) ([]entity.File, error)
	GetDownloadURL(ctx context.Context, userID uint, fileID uint) (string, error)
	GetStream(ctx context.Context, userID uint, fileID uint) (io.ReadCloser, *entity.File, error)
	Delete(ctx context.Context, userID uint, fileID uint) error
}

type fileService struct {
	fileRepo      repository.FileRepository
	storage       utils.StorageClient
	bucketName    string
	expiryMinutes int
}

func NewFileService(fileRepo repository.FileRepository, storage utils.StorageClient, bucketName string, expiryMinutes int) FileService {
	return &fileService{
		fileRepo:      fileRepo,
		storage:       storage,
		bucketName:    bucketName,
		expiryMinutes: expiryMinutes,
	}
}

func (s *fileService) Upload(ctx context.Context, userID uint, filename string, size int64, contentType string, reader io.Reader) (*entity.File, error) {
	// 1. Validasi Ukuran Berkas (Max 10MB)
	if size > 10*1024*1024 {
		return nil, ErrFileTooLarge
	}

	// 2. Validasi Tipe Konten Berkas (JPEG, PNG, PDF)
	if contentType != "image/jpeg" && contentType != "image/png" && contentType != "application/pdf" {
		return nil, ErrInvalidFileType
	}

	// 3. Buat Object Key MinIO Terisolasi Per Pengguna
	objectKey := fmt.Sprintf("user_%d/%d_%s", userID, time.Now().UnixNano(), filename)

	// 4. Catat Metadata dengan status PENDING ke PostgreSQL
	fileMeta := &entity.File{
		UserID:      userID,
		FileName:    filename,
		FileSize:    size,
		ContentType: contentType,
		ObjectKey:   objectKey,
		Status:      "PENDING",
	}

	if err := s.fileRepo.Create(ctx, fileMeta); err != nil {
		return nil, err
	}

	// 5. Unggah Berkas Fisik ke MinIO
	err := s.storage.UploadFile(ctx, s.bucketName, objectKey, reader, size, contentType)
	if err != nil {
		// Compensating Rollback: Hapus baris metadata jika upload fisik ke storage gagal
		_ = s.fileRepo.Delete(ctx, fileMeta.ID)
		return nil, fmt.Errorf("failed to upload physical file: %w", err)
	}

	// 6. Update status metadata menjadi SUCCESS
	fileMeta.Status = "SUCCESS"
	if err := s.fileRepo.Update(ctx, fileMeta); err != nil {
		return nil, err
	}

	return fileMeta, nil
}

func (s *fileService) GetList(ctx context.Context, userID uint) ([]entity.File, error) {
	return s.fileRepo.GetByUserID(ctx, userID)
}

func (s *fileService) GetDownloadURL(ctx context.Context, userID uint, fileID uint) (string, error) {
	file, err := s.fileRepo.GetByID(ctx, fileID)
	if err != nil {
		return "", err
	}
	if file == nil || file.Status != "SUCCESS" {
		return "", ErrFileNotFound
	}

	// Cek Kepemilikan Berkas
	if file.UserID != userID {
		return "", ErrPermissionDenied
	}

	// Buat S3 Presigned URL
	expiry := time.Duration(s.expiryMinutes) * time.Minute
	presignedURL, err := s.storage.GetPresignedURL(ctx, s.bucketName, file.ObjectKey, expiry)
	if err != nil {
		return "", err
	}

	return presignedURL.String(), nil
}

func (s *fileService) GetStream(ctx context.Context, userID uint, fileID uint) (io.ReadCloser, *entity.File, error) {
	file, err := s.fileRepo.GetByID(ctx, fileID)
	if err != nil {
		return nil, nil, err
	}
	if file == nil || file.Status != "SUCCESS" {
		return nil, nil, ErrFileNotFound
	}

	// Cek Kepemilikan Berkas
	if file.UserID != userID {
		return nil, nil, ErrPermissionDenied
	}

	// Ambil direct stream dari MinIO
	stream, err := s.storage.GetFileStream(ctx, s.bucketName, file.ObjectKey)
	if err != nil {
		return nil, nil, err
	}

	return stream, file, nil
}

func (s *fileService) Delete(ctx context.Context, userID uint, fileID uint) error {
	file, err := s.fileRepo.GetByID(ctx, fileID)
	if err != nil {
		return err
	}
	if file == nil {
		return ErrFileNotFound
	}

	// Cek Kepemilikan Berkas
	if file.UserID != userID {
		return ErrPermissionDenied
	}

	// 1. Hapus metadata di database relasional
	if err := s.fileRepo.Delete(ctx, file.ID); err != nil {
		return err
	}

	// 2. Hapus berkas fisik di MinIO
	// Catatan: Jika berkas fisik tidak ada di storage (misal terlanjur hilang), error diabaikan agar record tetap terhapus bersih.
	_ = s.storage.DeleteFile(ctx, s.bucketName, file.ObjectKey)

	return nil
}
