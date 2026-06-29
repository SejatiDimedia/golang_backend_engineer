package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/timurdian/url-shortener/internal/entity"
	"github.com/timurdian/url-shortener/internal/service"
)

// mockURLService mengimplementasikan service.URLService untuk pengujian HTTP Handler
type mockURLService struct {
	shortenFunc           func(ctx context.Context, longURL, customAlias string, expiresAt *time.Time) (*entity.URL, error)
	getAndRecordClickFunc func(ctx context.Context, code string) (*entity.URL, error)
	getStatsFunc          func(ctx context.Context, code string) (*entity.URL, error)
}

func (m *mockURLService) Shorten(ctx context.Context, longURL, customAlias string, expiresAt *time.Time) (*entity.URL, error) {
	return m.shortenFunc(ctx, longURL, customAlias, expiresAt)
}

func (m *mockURLService) GetAndRecordClick(ctx context.Context, code string) (*entity.URL, error) {
	return m.getAndRecordClickFunc(ctx, code)
}

func (m *mockURLService) GetStats(ctx context.Context, code string) (*entity.URL, error) {
	return m.getStatsFunc(ctx, code)
}

func TestShorten_HandlerSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := &mockURLService{
		shortenFunc: func(ctx context.Context, longURL, customAlias string, expiresAt *time.Time) (*entity.URL, error) {
			return &entity.URL{
				ID:         1,
				LongURL:    longURL,
				ShortCode:  "abcd123",
				ExpiresAt:  expiresAt,
				CreatedAt:  time.Now(),
				ClickCount: 0,
			}, nil
		},
	}

	h := NewURLHandler(mockSvc)
	r := gin.Default()
	r.POST("/shorten", h.Shorten)

	payload := ShortenRequest{
		LongURL: "https://example.com",
	}
	body, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/shorten", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201 Created, got %d", w.Code)
	}

	var resp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)

	if resp["short_code"] != "abcd123" {
		t.Errorf("expected short_code 'abcd123', got '%v'", resp["short_code"])
	}
}

func TestShorten_HandlerValidationError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	h := NewURLHandler(&mockURLService{})
	r := gin.Default()
	r.POST("/shorten", h.Shorten)

	// URL tidak valid karena format string acak
	payload := ShortenRequest{
		LongURL: "not-a-valid-url",
	}
	body, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/shorten", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 Bad Request, got %d", w.Code)
	}
}

func TestRedirect_HandlerSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := &mockURLService{
		getAndRecordClickFunc: func(ctx context.Context, code string) (*entity.URL, error) {
			return &entity.URL{
				LongURL:   "https://example.com/target",
				ShortCode: code,
			}, nil
		},
	}

	h := NewURLHandler(mockSvc)
	r := gin.Default()
	r.GET("/r/:short_code", h.Redirect)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/r/my-code", nil)

	r.ServeHTTP(w, req)

	if w.Code != http.StatusFound {
		t.Errorf("expected status 302 Found, got %d", w.Code)
	}

	loc := w.Header().Get("Location")
	if loc != "https://example.com/target" {
		t.Errorf("expected redirect location 'https://example.com/target', got '%s'", loc)
	}
}

func TestRedirect_HandlerNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := &mockURLService{
		getAndRecordClickFunc: func(ctx context.Context, code string) (*entity.URL, error) {
			return nil, service.ErrURLNotFound
		},
	}

	h := NewURLHandler(mockSvc)
	r := gin.Default()
	r.GET("/r/:short_code", h.Redirect)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/r/unknown", nil)

	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404 Not Found, got %d", w.Code)
	}
}
