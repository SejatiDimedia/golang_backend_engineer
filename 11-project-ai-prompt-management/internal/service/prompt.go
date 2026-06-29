package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/timurdian/prompt-management/internal/entity"
	"github.com/timurdian/prompt-management/internal/repository"
	"github.com/timurdian/prompt-management/internal/utils"
)

var (
	ErrUnauthorizedAccess = errors.New("unauthorized access to workspace")
	ErrPromptNotFound     = errors.New("prompt not found")
	ErrVersionNotFound    = errors.New("prompt version not found")
	ErrNoActiveVersion    = errors.New("no active version found for this prompt")
	ErrInvalidApiKey      = errors.New("invalid or expired API key")
)

type PromptService interface {
	CreateWorkspace(ctx context.Context, name string, ownerID uint) (*entity.Workspace, error)
	GetUserWorkspaces(ctx context.Context, userID uint) ([]entity.Workspace, error)

	CreatePrompt(ctx context.Context, wsID, userID uint, name, desc string) (*entity.Prompt, error)
	GetPrompt(ctx context.Context, promptID, userID uint) (*entity.Prompt, error)
	GetWorkspacePrompts(ctx context.Context, wsID, userID uint) ([]entity.Prompt, error)

	CreatePromptVersion(ctx context.Context, promptID, userID uint, text string) (*entity.PromptVersion, error)
	ActivatePromptVersion(ctx context.Context, promptID, userID uint, versionNum int) error

	CreateApiKey(ctx context.Context, wsID, userID uint, keyName string) (string, *entity.ApiKey, error)
	GetWorkspaceApiKeys(ctx context.Context, wsID, userID uint) ([]entity.ApiKey, error)
	RevokeApiKey(ctx context.Context, wsID, userID, keyID uint) error

	CompilePrompt(ctx context.Context, apiHash string, promptID uint, vars map[string]string) (string, int, error)
	GetWorkspaceAnalytics(ctx context.Context, wsID, userID uint) ([]entity.AnalyticsLog, error)
}

type promptService struct {
	repo      repository.PromptRepository
	rdb       *redis.Client
	analytics AnalyticsService
}

func NewPromptService(repo repository.PromptRepository, rdb *redis.Client, analytics AnalyticsService) PromptService {
	return &promptService{repo: repo, rdb: rdb, analytics: analytics}
}

func (s *promptService) validateAccess(ctx context.Context, wsID, userID uint) (bool, string, error) {
	isMember, role, err := s.repo.IsWorkspaceMember(ctx, wsID, userID)
	if err != nil {
		return false, "", err
	}
	if !isMember {
		return false, "", ErrUnauthorizedAccess
	}
	return true, role, nil
}

func (s *promptService) CreateWorkspace(ctx context.Context, name string, ownerID uint) (*entity.Workspace, error) {
	ws := &entity.Workspace{
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := s.repo.CreateWorkspace(ctx, ws, ownerID)
	if err != nil {
		return nil, err
	}
	return ws, nil
}

func (s *promptService) GetUserWorkspaces(ctx context.Context, userID uint) ([]entity.Workspace, error) {
	return s.repo.GetUserWorkspaces(ctx, userID)
}

func (s *promptService) CreatePrompt(ctx context.Context, wsID, userID uint, name, desc string) (*entity.Prompt, error) {
	if _, _, err := s.validateAccess(ctx, wsID, userID); err != nil {
		return nil, err
	}

	prompt := &entity.Prompt{
		WorkspaceID: wsID,
		Name:        name,
		Description: desc,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := s.repo.CreatePrompt(ctx, prompt)
	if err != nil {
		return nil, err
	}
	return prompt, nil
}

func (s *promptService) GetPrompt(ctx context.Context, promptID, userID uint) (*entity.Prompt, error) {
	prompt, err := s.repo.GetPromptByID(ctx, promptID)
	if err != nil {
		return nil, err
	}
	if prompt == nil {
		return nil, ErrPromptNotFound
	}

	if _, _, err := s.validateAccess(ctx, prompt.WorkspaceID, userID); err != nil {
		return nil, err
	}

	return prompt, nil
}

func (s *promptService) GetWorkspacePrompts(ctx context.Context, wsID, userID uint) ([]entity.Prompt, error) {
	if _, _, err := s.validateAccess(ctx, wsID, userID); err != nil {
		return nil, err
	}
	return s.repo.GetPromptsByWorkspace(ctx, wsID)
}

func (s *promptService) CreatePromptVersion(ctx context.Context, promptID, userID uint, text string) (*entity.PromptVersion, error) {
	prompt, err := s.repo.GetPromptByID(ctx, promptID)
	if err != nil {
		return nil, err
	}
	if prompt == nil {
		return nil, ErrPromptNotFound
	}

	if _, _, err := s.validateAccess(ctx, prompt.WorkspaceID, userID); err != nil {
		return nil, err
	}

	pv := &entity.PromptVersion{
		PromptID:   promptID,
		PromptText: text,
		Status:     "DRAFT",
	}

	err = s.repo.CreatePromptVersion(ctx, pv)
	if err != nil {
		return nil, err
	}
	return pv, nil
}

func (s *promptService) ActivatePromptVersion(ctx context.Context, promptID, userID uint, versionNum int) error {
	prompt, err := s.repo.GetPromptByID(ctx, promptID)
	if err != nil {
		return err
	}
	if prompt == nil {
		return ErrPromptNotFound
	}

	if _, _, err := s.validateAccess(ctx, prompt.WorkspaceID, userID); err != nil {
		return err
	}

	return s.repo.ActivatePromptVersion(ctx, promptID, versionNum)
}

func (s *promptService) CreateApiKey(ctx context.Context, wsID, userID uint, keyName string) (string, *entity.ApiKey, error) {
	if _, role, err := s.validateAccess(ctx, wsID, userID); err != nil {
		return "", nil, err
	} else if role != "admin" {
		return "", nil, errors.New("only admins can manage API keys")
	}

	rawKey, hash, masked, err := utils.GenerateAPIKey()
	if err != nil {
		return "", nil, err
	}

	apiKey := &entity.ApiKey{
		WorkspaceID: wsID,
		Name:        keyName,
		KeyHash:     hash,
		MaskedKey:   masked,
		CreatedAt:   time.Now(),
		ExpiresAt:   time.Now().Add(365 * 24 * time.Hour), // 1 Year TTL
	}

	err = s.repo.CreateApiKey(ctx, apiKey)
	if err != nil {
		return "", nil, err
	}

	// Cache API Key di Redis
	cacheKey := fmt.Sprintf("apikey:%s", hash)
	_ = s.rdb.Set(ctx, cacheKey, fmt.Sprintf("%d", wsID), 1*time.Hour).Err()

	return rawKey, apiKey, nil
}

func (s *promptService) GetWorkspaceApiKeys(ctx context.Context, wsID, userID uint) ([]entity.ApiKey, error) {
	if _, _, err := s.validateAccess(ctx, wsID, userID); err != nil {
		return nil, err
	}
	return s.repo.GetApiKeysByWorkspace(ctx, wsID)
}

func (s *promptService) RevokeApiKey(ctx context.Context, wsID, userID, keyID uint) error {
	if _, role, err := s.validateAccess(ctx, wsID, userID); err != nil {
		return err
	} else if role != "admin" {
		return errors.New("only admins can manage API keys")
	}

	// Cari key untuk dihapus dari Redis
	var keys []entity.ApiKey
	var targetKey *entity.ApiKey
	keys, err := s.repo.GetApiKeysByWorkspace(ctx, wsID)
	if err == nil {
		for _, k := range keys {
			if k.ID == keyID {
				targetKey = &k
				break
			}
		}
	}

	err = s.repo.RevokeApiKey(ctx, keyID)
	if err != nil {
		return err
	}

	if targetKey != nil {
		cacheKey := fmt.Sprintf("apikey:%s", targetKey.KeyHash)
		_ = s.rdb.Del(ctx, cacheKey).Err()
	}

	return nil
}

func (s *promptService) CompilePrompt(ctx context.Context, apiHash string, promptID uint, vars map[string]string) (string, int, error) {
	startTime := time.Now()

	// 1. Validasi API Key
	apiKey, err := s.repo.GetApiKeyByHash(ctx, apiHash)
	if err != nil {
		return "", 0, err
	}
	if apiKey == nil {
		return "", 0, ErrInvalidApiKey
	}

	// 2. Load Prompt
	prompt, err := s.repo.GetPromptByID(ctx, promptID)
	if err != nil {
		return "", 0, err
	}
	if prompt == nil {
		return "", 0, ErrPromptNotFound
	}

	// Pastikan Prompt berada di Workspace pemilik API Key
	if prompt.WorkspaceID != apiKey.WorkspaceID {
		return "", 0, ErrUnauthorizedAccess
	}

	// 3. Load Active Version
	activeVersion, err := s.repo.GetActivePromptVersion(ctx, promptID)
	if err != nil {
		return "", 0, err
	}
	if activeVersion == nil {
		return "", 0, ErrNoActiveVersion
	}

	// 4. Compile menggunakan Compiler Engine
	compiledText, wordCount := utils.CompilePrompt(activeVersion.PromptText, vars)

	// 5. Asynchronous Logging Analytics
	latency := time.Since(startTime).Milliseconds()
	logEntry := &entity.AnalyticsLog{
		ApiKeyID:      apiKey.ID,
		PromptID:      promptID,
		LatencyMs:     latency,
		TokenEstimate: wordCount,
		ResponseCode:  200,
		CreatedAt:     time.Now(),
	}

	s.analytics.LogAsync(logEntry)

	return compiledText, wordCount, nil
}

func (s *promptService) GetWorkspaceAnalytics(ctx context.Context, wsID, userID uint) ([]entity.AnalyticsLog, error) {
	if _, _, err := s.validateAccess(ctx, wsID, userID); err != nil {
		return nil, err
	}
	return s.repo.GetWorkspaceAnalytics(ctx, wsID)
}
