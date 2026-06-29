package repository

import (
	"context"
	"errors"
	"time"

	"github.com/timurdian/prompt-management/internal/entity"
	"gorm.io/gorm"
)

type PromptRepository interface {
	CreateWorkspace(ctx context.Context, ws *entity.Workspace, ownerID uint) error
	GetWorkspaceByID(ctx context.Context, id uint) (*entity.Workspace, error)
	IsWorkspaceMember(ctx context.Context, wsID, userID uint) (bool, string, error)
	GetUserWorkspaces(ctx context.Context, userID uint) ([]entity.Workspace, error)

	CreatePrompt(ctx context.Context, prompt *entity.Prompt) error
	GetPromptByID(ctx context.Context, id uint) (*entity.Prompt, error)
	GetPromptsByWorkspace(ctx context.Context, wsID uint) ([]entity.Prompt, error)

	CreatePromptVersion(ctx context.Context, pv *entity.PromptVersion) error
	GetPromptVersion(ctx context.Context, promptID uint, versionNum int) (*entity.PromptVersion, error)
	GetActivePromptVersion(ctx context.Context, promptID uint) (*entity.PromptVersion, error)
	ActivatePromptVersion(ctx context.Context, promptID uint, versionNum int) error

	CreateApiKey(ctx context.Context, key *entity.ApiKey) error
	GetApiKeyByHash(ctx context.Context, hash string) (*entity.ApiKey, error)
	GetApiKeysByWorkspace(ctx context.Context, wsID uint) ([]entity.ApiKey, error)
	RevokeApiKey(ctx context.Context, keyID uint) error

	LogAnalytics(ctx context.Context, log *entity.AnalyticsLog) error
	GetWorkspaceAnalytics(ctx context.Context, wsID uint) ([]entity.AnalyticsLog, error)
}

type gormPromptRepository struct {
	db *gorm.DB
}

func NewPromptRepository(db *gorm.DB) PromptRepository {
	return &gormPromptRepository{db: db}
}

func (r *gormPromptRepository) CreateWorkspace(ctx context.Context, ws *entity.Workspace, ownerID uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(ws).Error; err != nil {
			return err
		}

		member := &entity.WorkspaceMember{
			WorkspaceID: ws.ID,
			UserID:      ownerID,
			Role:        "admin",
			CreatedAt:   time.Now(),
		}

		return tx.Create(member).Error
	})
}

func (r *gormPromptRepository) GetWorkspaceByID(ctx context.Context, id uint) (*entity.Workspace, error) {
	var ws entity.Workspace
	err := r.db.WithContext(ctx).First(&ws, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &ws, nil
}

func (r *gormPromptRepository) IsWorkspaceMember(ctx context.Context, wsID, userID uint) (bool, string, error) {
	var member entity.WorkspaceMember
	err := r.db.WithContext(ctx).Where("workspace_id = ? AND user_id = ?", wsID, userID).First(&member).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, "", nil
		}
		return false, "", err
	}
	return true, member.Role, nil
}

func (r *gormPromptRepository) GetUserWorkspaces(ctx context.Context, userID uint) ([]entity.Workspace, error) {
	var workspaces []entity.Workspace
	err := r.db.WithContext(ctx).
		Joins("JOIN workspace_members ON workspace_members.workspace_id = workspaces.id").
		Where("workspace_members.user_id = ?", userID).
		Find(&workspaces).Error
	return workspaces, err
}

func (r *gormPromptRepository) CreatePrompt(ctx context.Context, prompt *entity.Prompt) error {
	return r.db.WithContext(ctx).Create(prompt).Error
}

func (r *gormPromptRepository) GetPromptByID(ctx context.Context, id uint) (*entity.Prompt, error) {
	var prompt entity.Prompt
	err := r.db.WithContext(ctx).Preload("Versions").First(&prompt, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &prompt, nil
}

func (r *gormPromptRepository) GetPromptsByWorkspace(ctx context.Context, wsID uint) ([]entity.Prompt, error) {
	var prompts []entity.Prompt
	err := r.db.WithContext(ctx).Where("workspace_id = ?", wsID).Find(&prompts).Error
	return prompts, err
}

func (r *gormPromptRepository) CreatePromptVersion(ctx context.Context, pv *entity.PromptVersion) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Dapatkan nomor versi terakhir
		var lastNum int
		err := tx.Model(&entity.PromptVersion{}).
			Where("prompt_id = ?", pv.PromptID).
			Select("COALESCE(MAX(version_number), 0)").
			Scan(&lastNum).Error
		if err != nil {
			return err
		}

		pv.VersionNumber = lastNum + 1
		pv.CreatedAt = time.Now()
		return tx.Create(pv).Error
	})
}

func (r *gormPromptRepository) GetPromptVersion(ctx context.Context, promptID uint, versionNum int) (*entity.PromptVersion, error) {
	var pv entity.PromptVersion
	err := r.db.WithContext(ctx).Where("prompt_id = ? AND version_number = ?", promptID, versionNum).First(&pv).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &pv, nil
}

func (r *gormPromptRepository) GetActivePromptVersion(ctx context.Context, promptID uint) (*entity.PromptVersion, error) {
	var pv entity.PromptVersion
	err := r.db.WithContext(ctx).Where("prompt_id = ? AND status = 'ACTIVE'", promptID).First(&pv).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &pv, nil
}

func (r *gormPromptRepository) ActivatePromptVersion(ctx context.Context, promptID uint, versionNum int) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Nonaktifkan semua versi aktif sebelumnya
		err := tx.Model(&entity.PromptVersion{}).
			Where("prompt_id = ? AND status = 'ACTIVE'", promptID).
			Update("status", "DRAFT").Error
		if err != nil {
			return err
		}

		// Aktifkan versi terpilih
		res := tx.Model(&entity.PromptVersion{}).
			Where("prompt_id = ? AND version_number = ?", promptID, versionNum).
			Update("status", "ACTIVE")
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return errors.New("prompt version not found")
		}
		return nil
	})
}

func (r *gormPromptRepository) CreateApiKey(ctx context.Context, key *entity.ApiKey) error {
	return r.db.WithContext(ctx).Create(key).Error
}

func (r *gormPromptRepository) GetApiKeyByHash(ctx context.Context, hash string) (*entity.ApiKey, error) {
	var key entity.ApiKey
	err := r.db.WithContext(ctx).Where("key_hash = ? AND expires_at > ?", hash, time.Now()).First(&key).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &key, nil
}

func (r *gormPromptRepository) GetApiKeysByWorkspace(ctx context.Context, wsID uint) ([]entity.ApiKey, error) {
	var keys []entity.ApiKey
	err := r.db.WithContext(ctx).Where("workspace_id = ?", wsID).Find(&keys).Error
	return keys, err
}

func (r *gormPromptRepository) RevokeApiKey(ctx context.Context, keyID uint) error {
	return r.db.WithContext(ctx).Delete(&entity.ApiKey{}, keyID).Error
}

func (r *gormPromptRepository) LogAnalytics(ctx context.Context, log *entity.AnalyticsLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *gormPromptRepository) GetWorkspaceAnalytics(ctx context.Context, wsID uint) ([]entity.AnalyticsLog, error) {
	var logs []entity.AnalyticsLog
	err := r.db.WithContext(ctx).
		Joins("JOIN prompts ON prompts.id = analytics_logs.prompt_id").
		Where("prompts.workspace_id = ?", wsID).
		Order("analytics_logs.created_at DESC").
		Find(&logs).Error
	return logs, err
}
