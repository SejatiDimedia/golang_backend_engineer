package service

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/timurdian/url-shortener/internal/entity"
)

// mockURLRepository mengimplementasikan repository.URLRepository untuk kebutuhan pengujian unit
type mockURLRepository struct {
	urls map[string]*entity.URL
}

func newMockURLRepository() *mockURLRepository {
	return &mockURLRepository{
		urls: make(map[string]*entity.URL),
	}
}

func (m *mockURLRepository) Create(ctx context.Context, url *entity.URL) error {
	m.urls[url.ShortCode] = url
	return nil
}

func (m *mockURLRepository) GetByShortCode(ctx context.Context, code string) (*entity.URL, error) {
	u, exists := m.urls[code]
	if !exists {
		return nil, nil
	}
	return u, nil
}

func (m *mockURLRepository) IncrementClick(ctx context.Context, code string) error {
	if u, exists := m.urls[code]; exists {
		u.ClickCount++
		return nil
	}
	return errors.New("url not found")
}

func TestShorten_Success(t *testing.T) {
	repo := newMockURLRepository()
	svc := NewURLService(repo)

	longURL := "https://google.com"
	urlObj, err := svc.Shorten(context.Background(), longURL, "", nil)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if urlObj.LongURL != longURL {
		t.Errorf("expected LongURL %s, got %s", longURL, urlObj.LongURL)
	}

	if urlObj.ShortCode == "" {
		t.Error("expected generated short code, got empty string")
	}

	// Memastikan format short code aman di URL (Base64 URL-Safe)
	if strings.Contains(urlObj.ShortCode, "+") || strings.Contains(urlObj.ShortCode, "/") || strings.Contains(urlObj.ShortCode, "=") {
		t.Errorf("generated short code contains non-URL-safe characters: %s", urlObj.ShortCode)
	}
}

func TestShorten_InvalidURL(t *testing.T) {
	repo := newMockURLRepository()
	svc := NewURLService(repo)

	invalidURLs := []string{
		"google.com",         // skema hilang
		"ftp://google.com",   // skema salah
		"invalid-url-string",
	}

	for _, badURL := range invalidURLs {
		_, err := svc.Shorten(context.Background(), badURL, "", nil)
		if err == nil {
			t.Errorf("expected error for URL '%s', but got none", badURL)
		}
	}
}

func TestShorten_CustomAlias(t *testing.T) {
	repo := newMockURLRepository()
	svc := NewURLService(repo)

	alias := "g-search"
	longURL := "https://google.com"

	urlObj, err := svc.Shorten(context.Background(), longURL, alias, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if urlObj.ShortCode != alias {
		t.Errorf("expected short code to be custom alias '%s', got '%s'", alias, urlObj.ShortCode)
	}

	// Uji alias konflik
	_, err = svc.Shorten(context.Background(), "https://yahoo.com", alias, nil)
	if err == nil {
		t.Error("expected error due to alias conflict, got none")
	}
}

func TestGetAndRecordClick_Expired(t *testing.T) {
	repo := newMockURLRepository()
	svc := NewURLService(repo)

	// URL kedaluwarsa 1 jam yang lalu
	pastTime := time.Now().Add(-1 * time.Hour)
	urlObj, _ := svc.Shorten(context.Background(), "https://google.com", "expired-link", &pastTime)

	// Pastikan terdaftar
	if urlObj == nil {
		t.Fatal("failed to set up test expired link")
	}

	_, err := svc.GetAndRecordClick(context.Background(), "expired-link")
	if err == nil {
		t.Error("expected error for expired URL redirect, got none")
	}
}

func TestGetAndRecordClick_IncrementClickCount(t *testing.T) {
	repo := newMockURLRepository()
	svc := NewURLService(repo)

	urlObj, _ := svc.Shorten(context.Background(), "https://google.com", "my-link", nil)
	if urlObj.ClickCount != 0 {
		t.Errorf("expected click count to start at 0, got %d", urlObj.ClickCount)
	}

	// Redirect 1
	res, err := svc.GetAndRecordClick(context.Background(), "my-link")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if res.ClickCount != 1 {
		t.Errorf("expected click count to be 1, got %d", res.ClickCount)
	}

	// Redirect 2
	res, err = svc.GetAndRecordClick(context.Background(), "my-link")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if res.ClickCount != 2 {
		t.Errorf("expected click count to be 2, got %d", res.ClickCount)
	}
}
