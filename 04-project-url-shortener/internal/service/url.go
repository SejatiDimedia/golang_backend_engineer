package service

import (
	"context"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"net/url"
	"time"

	"github.com/timurdian/url-shortener/internal/entity"
	"github.com/timurdian/url-shortener/internal/repository"
)

var (
	ErrURLNotFound   = errors.New("url not found")
	ErrURLExpired     = errors.New("url has expired")
	ErrInvalidURL    = errors.New("invalid url format")
	ErrAliasConflict = errors.New("custom alias already exists")
)

type URLService interface {
	Shorten(ctx context.Context, longURL, customAlias string, expiresAt *time.Time) (*entity.URL, error)
	GetAndRecordClick(ctx context.Context, code string) (*entity.URL, error)
	GetStats(ctx context.Context, code string) (*entity.URL, error)
}

type urlService struct {
	repo repository.URLRepository
}

func NewURLService(repo repository.URLRepository) URLService {
	return &urlService{repo: repo}
}

func (s *urlService) Shorten(ctx context.Context, longURL, customAlias string, expiresAt *time.Time) (*entity.URL, error) {
	// 1. Validasi URL target
	if err := validateURL(longURL); err != nil {
		return nil, ErrInvalidURL
	}

	var shortCode string
	if customAlias != "" {
		// Pastikan alias kustom belum digunakan
		existing, err := s.repo.GetByShortCode(ctx, customAlias)
		if err != nil {
			return nil, err
		}
		if existing != nil {
			return nil, ErrAliasConflict
		}
		shortCode = customAlias
	} else {
		// Generate short code menggunakan timestamp nanodetik (Base64 URL-safe)
		shortCode = generateShortCode()
	}

	urlObj := &entity.URL{
		LongURL:    longURL,
		ShortCode:  shortCode,
		ClickCount: 0,
		ExpiresAt:  expiresAt,
	}

	if err := s.repo.Create(ctx, urlObj); err != nil {
		return nil, err
	}

	return urlObj, nil
}

func (s *urlService) GetAndRecordClick(ctx context.Context, code string) (*entity.URL, error) {
	urlObj, err := s.repo.GetByShortCode(ctx, code)
	if err != nil {
		return nil, err
	}
	if urlObj == nil {
		return nil, ErrURLNotFound
	}

	// Cek apakah kedaluwarsa
	if urlObj.IsExpired() {
		return nil, ErrURLExpired
	}

	// Rekam klik (sinkron)
	if err := s.repo.IncrementClick(ctx, code); err != nil {
		return nil, err
	}

	return urlObj, nil
}

func (s *urlService) GetStats(ctx context.Context, code string) (*entity.URL, error) {
	urlObj, err := s.repo.GetByShortCode(ctx, code)
	if err != nil {
		return nil, err
	}
	if urlObj == nil {
		return nil, ErrURLNotFound
	}

	return urlObj, nil
}

// generateShortCode membuat kode pendek berbasis timestamp nanodetik di-encode ke Base64 URL-Safe
func generateShortCode() string {
	now := time.Now().UnixNano()
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(now))
	
	// Gunakan RawURLEncoding untuk membuang padding '=' dan karakter non-URL-safe (+, /)
	encoded := base64.RawURLEncoding.EncodeToString(buf)
	
	// Untuk estetika short code, kita hilangkan padding byte di awal jika nilainya 0 (karena timestamp positif)
	// Namun agar tetap aman dan unik, string utuh (sekitar 11 karakter) sangat disarankan.
	return encoded
}

// validateURL mengecek validitas format target URL
func validateURL(targetURL string) error {
	u, err := url.ParseRequestURI(targetURL)
	if err != nil {
		return err
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return errors.New("invalid scheme")
	}
	return nil
}
