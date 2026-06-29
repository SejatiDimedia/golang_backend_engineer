package middleware

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/timurdian/prompt-management/internal/entity"
	"github.com/timurdian/prompt-management/internal/utils"
)

type mockPromptRepo struct {
	apiKeys map[string]*entity.ApiKey
}

func (m *mockPromptRepo) CreateWorkspace(ctx context.Context, ws *entity.Workspace, ownerID uint) error {
	return nil
}
func (m *mockPromptRepo) GetWorkspaceByID(ctx context.Context, id uint) (*entity.Workspace, error) {
	return nil, nil
}
func (m *mockPromptRepo) IsWorkspaceMember(ctx context.Context, wsID, userID uint) (bool, string, error) {
	return true, "admin", nil
}
func (m *mockPromptRepo) GetUserWorkspaces(ctx context.Context, userID uint) ([]entity.Workspace, error) {
	return nil, nil
}
func (m *mockPromptRepo) CreatePrompt(ctx context.Context, prompt *entity.Prompt) error {
	return nil
}
func (m *mockPromptRepo) GetPromptByID(ctx context.Context, id uint) (*entity.Prompt, error) {
	return nil, nil
}
func (m *mockPromptRepo) GetPromptsByWorkspace(ctx context.Context, wsID uint) ([]entity.Prompt, error) {
	return nil, nil
}
func (m *mockPromptRepo) CreatePromptVersion(ctx context.Context, pv *entity.PromptVersion) error {
	return nil
}
func (m *mockPromptRepo) GetPromptVersion(ctx context.Context, promptID uint, versionNum int) (*entity.PromptVersion, error) {
	return nil, nil
}
func (m *mockPromptRepo) GetActivePromptVersion(ctx context.Context, promptID uint) (*entity.PromptVersion, error) {
	return nil, nil
}
func (m *mockPromptRepo) ActivatePromptVersion(ctx context.Context, promptID uint, versionNum int) error {
	return nil
}
func (m *mockPromptRepo) CreateApiKey(ctx context.Context, key *entity.ApiKey) error {
	return nil
}
func (m *mockPromptRepo) GetApiKeyByHash(ctx context.Context, hash string) (*entity.ApiKey, error) {
	key, exists := m.apiKeys[hash]
	if !exists {
		return nil, nil
	}
	return key, nil
}
func (m *mockPromptRepo) GetApiKeysByWorkspace(ctx context.Context, wsID uint) ([]entity.ApiKey, error) {
	return nil, nil
}
func (m *mockPromptRepo) RevokeApiKey(ctx context.Context, keyID uint) error {
	return nil
}
func (m *mockPromptRepo) LogAnalytics(ctx context.Context, log *entity.AnalyticsLog) error {
	return nil
}
func (m *mockPromptRepo) GetWorkspaceAnalytics(ctx context.Context, wsID uint) ([]entity.AnalyticsLog, error) {
	return nil, nil
}

func TestAPIKeyMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("Failed to run miniredis: %v", err)
	}
	defer mr.Close()

	rdb := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	defer rdb.Close()

	rawKey := "prompt_live_mysecureapikeyhashvalue"
	hash := utils.HashAPIKey(rawKey)

	repo := &mockPromptRepo{
		apiKeys: map[string]*entity.ApiKey{
			hash: {
				ID:          1,
				WorkspaceID: 10,
				Name:        "Test Key",
				KeyHash:     hash,
				MaskedKey:   "prompt_live_xxxx...value",
				ExpiresAt:   time.Now().Add(1 * time.Hour),
			},
		},
	}

	mw := NewAPIKeyMiddleware(repo, rdb)

	r := gin.New()
	r.Use(mw.APIKeyRequired())
	r.GET("/test-auth", func(c *gin.Context) {
		h, _ := c.Get("api_key_hash")
		c.JSON(http.StatusOK, gin.H{"status": "ok", "hash": h})
	})

	// 1. Request Tanpa Header
	req, _ := http.NewRequest(http.MethodGet, "/test-auth", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", w.Code)
	}

	// 2. Request dengan prefix salah
	req, _ = http.NewRequest(http.MethodGet, "/test-auth", nil)
	req.Header.Set("Authorization", "Bearer invalid_prefix_key")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", w.Code)
	}

	// 3. Request Valid (Cache Miss -> DB Hit -> Save Cache)
	req, _ = http.NewRequest(http.MethodGet, "/test-auth", nil)
	req.Header.Set("Authorization", "Bearer "+rawKey)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d, body: %s", w.Code, w.Body.String())
	}

	// Cek apakah key tersimpan di Redis cache
	cacheKey := fmt.Sprintf("apikey:%s", hash)
	val, err := rdb.Get(context.Background(), cacheKey).Result()
	if err != nil || val != "10" {
		t.Errorf("expected apikey in Redis cache with value '10', got value %s, err: %v", val, err)
	}

	// 4. Request Valid Kedua (Cache Hit!)
	req, _ = http.NewRequest(http.MethodGet, "/test-auth", nil)
	req.Header.Set("Authorization", "Bearer "+rawKey)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}
