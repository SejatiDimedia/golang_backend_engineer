package service

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/timurdian/prompt-management/internal/entity"
	"github.com/timurdian/prompt-management/internal/repository"
	"github.com/timurdian/prompt-management/internal/utils"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDBAndRedis(t *testing.T) (*gorm.DB, *redis.Client, func()) {
	// Setup SQLite in-memory
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open SQLite test db: %v", err)
	}

	// Auto-Migrate
	err = db.AutoMigrate(
		&entity.Workspace{},
		&entity.WorkspaceMember{},
		&entity.Prompt{},
		&entity.PromptVersion{},
		&entity.ApiKey{},
		&entity.AnalyticsLog{},
	)
	if err != nil {
		t.Fatalf("failed to run GORM auto-migrations: %v", err)
	}

	// Setup Miniredis
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}

	rdb := redis.NewClient(&redis.Options{Addr: mr.Addr()})

	cleanup := func() {
		sqlDB, _ := db.DB()
		_ = sqlDB.Close()
		_ = rdb.Close()
		mr.Close()
	}

	return db, rdb, cleanup
}

func TestPromptService_Workflow(t *testing.T) {
	db, rdb, cleanup := setupTestDBAndRedis(t)
	defer cleanup()

	repo := repository.NewPromptRepository(db)
	analytics := NewAnalyticsService(repo)

	// Start background analytics worker daemon
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	analytics.StartWorker(ctx)

	svc := NewPromptService(repo, rdb, analytics)

	userID := uint(1)
	otherUserID := uint(2)

	// 1. Create Workspace Tim
	ws, err := svc.CreateWorkspace(ctx, "AI Research", userID)
	if err != nil {
		t.Fatalf("CreateWorkspace failed: %v", err)
	}
	if ws.Name != "AI Research" {
		t.Errorf("expected workspace name 'AI Research', got %q", ws.Name)
	}

	// 2. Create Prompt Template
	prompt, err := svc.CreatePrompt(ctx, ws.ID, userID, "Summarize Text", "Summarizes long text into bullet points")
	if err != nil {
		t.Fatalf("CreatePrompt failed: %v", err)
	}

	// Uji validasi akses (User lain tidak boleh menulis prompt ke workspace)
	_, err = svc.CreatePrompt(ctx, ws.ID, otherUserID, "Hacker Prompt", "Unauthorized")
	if err == nil {
		t.Error("expected error when other user tries to access workspace, got nil")
	}

	// 3. Create Prompt Versions (v1, v2)
	v1, err := svc.CreatePromptVersion(ctx, prompt.ID, userID, "Summarize the following in {{count}} points: {{text}}")
	if err != nil {
		t.Fatalf("CreatePromptVersion v1 failed: %v", err)
	}
	if v1.VersionNumber != 1 {
		t.Errorf("expected version number 1, got %d", v1.VersionNumber)
	}

	v2, err := svc.CreatePromptVersion(ctx, prompt.ID, userID, "Create a short TL;DR summary: {{text}}")
	if err != nil {
		t.Fatalf("CreatePromptVersion v2 failed: %v", err)
	}
	if v2.VersionNumber != 2 {
		t.Errorf("expected version number 2, got %d", v2.VersionNumber)
	}

	// 4. Coba compile prompt sebelum diaktivasi -> Error NoActiveVersion
	_, _, err = svc.CompilePrompt(ctx, "some_hash", prompt.ID, nil)
	if err != ErrInvalidApiKey {
		t.Errorf("expected ErrInvalidApiKey before creating API Key, got %v", err)
	}

	// 5. Generate API Key
	rawKey, apiKey, err := svc.CreateApiKey(ctx, ws.ID, userID, "Production Server Key")
	if err != nil {
		t.Fatalf("CreateApiKey failed: %v", err)
	}
	if apiKey.MaskedKey == "" {
		t.Error("expected masked key representation to be present")
	}

	hash := utils.HashAPIKey(rawKey)

	// Uji compile setelah API Key ada tetapi belum ada versi ACTIVE -> ErrNoActiveVersion
	_, _, err = svc.CompilePrompt(ctx, hash, prompt.ID, nil)
	if err != ErrNoActiveVersion {
		t.Errorf("expected ErrNoActiveVersion, got %v", err)
	}

	// 6. Aktifkan Versi 1
	err = svc.ActivatePromptVersion(ctx, prompt.ID, userID, 1)
	if err != nil {
		t.Fatalf("ActivatePromptVersion failed: %v", err)
	}

	// 7. Lakukan Kompilasi Prompt Versi 1
	vars := map[string]string{
		"count": "3",
		"text":  "Golang is an open-source programming language developed by Google.",
	}
	compiledText, tokens, err := svc.CompilePrompt(ctx, hash, prompt.ID, vars)
	if err != nil {
		t.Fatalf("CompilePrompt failed: %v", err)
	}

	expectedCompiled := "Summarize the following in 3 points: Golang is an open-source programming language developed by Google."
	if compiledText != expectedCompiled {
		t.Errorf("expected compile output %q, got %q", expectedCompiled, compiledText)
	}
	if tokens == 0 {
		t.Error("expected estimated tokens count to be greater than 0")
	}

	// 8. Aktifkan Versi 2
	err = svc.ActivatePromptVersion(ctx, prompt.ID, userID, 2)
	if err != nil {
		t.Fatalf("Activate v2 failed: %v", err)
	}

	// Lakukan Kompilasi Prompt Versi 2 (teks prompt template berubah!)
	compiledText2, _, err := svc.CompilePrompt(ctx, hash, prompt.ID, map[string]string{
		"text": "Antigravity AI is pairs coding.",
	})
	if err != nil {
		t.Fatalf("CompilePrompt v2 failed: %v", err)
	}

	expectedCompiled2 := "Create a short TL;DR summary: Antigravity AI is pairs coding."
	if compiledText2 != expectedCompiled2 {
		t.Errorf("expected compile output v2 %q, got %q", expectedCompiled2, compiledText2)
	}

	// Tunggu sebentar agar asinkronous analytics logger menulis ke DB
	time.Sleep(100 * time.Millisecond)

	// 9. Verifikasi analytics log tersimpan di DB
	analyticsLogs, err := svc.GetWorkspaceAnalytics(ctx, ws.ID, userID)
	if err != nil {
		t.Fatalf("GetWorkspaceAnalytics failed: %v", err)
	}
	if len(analyticsLogs) < 2 {
		t.Errorf("expected at least 2 analytics log entry, got %d", len(analyticsLogs))
	}
}
