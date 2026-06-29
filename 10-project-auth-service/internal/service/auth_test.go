package service

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/timurdian/auth-service/internal/entity"
	"github.com/timurdian/auth-service/internal/repository"
	"github.com/timurdian/auth-service/internal/utils"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) (*gorm.DB, string, string) {
	// Setup SQLite in-memory
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite: %v", err)
	}

	// Auto-Migrate schema
	err = db.AutoMigrate(
		&entity.User{},
		&entity.Role{},
		&entity.Permission{},
		&entity.RefreshToken{},
		&entity.VerificationToken{},
		&entity.ResetToken{},
	)
	if err != nil {
		t.Fatalf("failed to auto migrate sqlite: %v", err)
	}

	// Setup certs temp
	tempDir, err := os.MkdirTemp("", "auth-certs")
	if err != nil {
		t.Fatalf("failed to create certs dir: %v", err)
	}

	privPath := filepath.Join(tempDir, "private.key")
	pubPath := filepath.Join(tempDir, "public.key")
	_ = utils.EnsureRSAKeys(privPath, pubPath)

	return db, privPath, pubPath
}

func cleanTestCerts(dir string) {
	_ = os.RemoveAll(filepath.Dir(dir))
}

func TestAuthService_RegisterAndLogin(t *testing.T) {
	db, priv, pub := setupTestDB(t)
	defer cleanTestCerts(priv)

	repo := repository.NewAuthRepository(db)
	txMgr := repository.NewTransactionManager(db)
	tokenMgr, _ := utils.NewTokenManager(priv, pub, 15)

	svc := NewAuthService(repo, txMgr, tokenMgr, nil, 7, 15)
	ctx := context.Background()

	// 1. Register User Baru
	user, err := svc.Register(ctx, "test@email.com", "password123")
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	if user.IsVerified {
		t.Errorf("expected user IsVerified to be false upon registration")
	}

	// Cek default role 'customer'
	roleName, _, err := repo.GetUserRolesAndPermissions(ctx, user.ID)
	if err != nil {
		t.Fatalf("failed to fetch user roles: %v", err)
	}
	if roleName != "customer" {
		t.Errorf("expected default role 'customer', got '%s'", roleName)
	}

	// 2. Cobalah Login sebelum email terverifikasi -> harus gagal (forbidden)
	_, _, err = svc.Login(ctx, "test@email.com", "password123")
	if err != ErrEmailNotVerified {
		t.Errorf("expected ErrEmailNotVerified error, got %v", err)
	}

	// 3. Verifikasi Email
	var vt entity.VerificationToken
	db.Where("user_id = ?", user.ID).First(&vt)

	err = svc.VerifyEmail(ctx, vt.Token)
	if err != nil {
		t.Fatalf("VerifyEmail failed: %v", err)
	}

	// 4. Login setelah email terverifikasi -> harus sukses
	access, refresh, err := svc.Login(ctx, "test@email.com", "password123")
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}

	if access == "" || refresh == "" {
		t.Errorf("expected access and refresh tokens to be returned")
	}
}

func TestAuthService_RefreshTokenRotation_ReplayAttack(t *testing.T) {
	db, priv, pub := setupTestDB(t)
	defer cleanTestCerts(priv)

	repo := repository.NewAuthRepository(db)
	txMgr := repository.NewTransactionManager(db)
	tokenMgr, _ := utils.NewTokenManager(priv, pub, 15)

	svc := NewAuthService(repo, txMgr, tokenMgr, nil, 7, 15)
	ctx := context.Background()

	// Register & Verify & Login
	user, _ := svc.Register(ctx, "user@rotate.com", "password123")
	_ = repo.UpdateUserVerification(ctx, user.ID, true)
	_, refresh1, _ := svc.Login(ctx, "user@rotate.com", "password123")

	// Tambahkan refresh token lain yang aktif (misal dari device kedua)
	rt2 := &entity.RefreshToken{
		Token:     "device2_token",
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(24 * time.Hour),
		IsRevoked: false,
	}
	_ = repo.CreateRefreshToken(ctx, rt2)

	// 1. Jalankan Rotasi Pertama -> Sukses
	time.Sleep(10 * time.Millisecond) // jeda waktu singkat
	access2, refresh2, err := svc.Refresh(ctx, refresh1)
	if err != nil {
		t.Fatalf("first refresh failed: %v", err)
	}
	if access2 == "" || refresh2 == "" {
		t.Fatalf("tokens empty")
	}

	// Verifikasi token 1 revoked di DB
	var oldRT entity.RefreshToken
	db.Where("token = ?", refresh1).First(&oldRT)
	if !oldRT.IsRevoked {
		t.Errorf("expected original refresh token to be revoked")
	}

	// 2. REPLAY ATTACK: Kirim kembali token 1 yang sudah di-revoke
	_, _, err = svc.Refresh(ctx, refresh1)
	if err != ErrInvalidOrExpiredToken {
		t.Errorf("expected replay attack error, got %v", err)
	}

	// Verifikasi logout massal (seluruh token user bersangkutan, termasuk device 2 & refresh 2, dicabut)
	var activeTokensCount int64
	db.Model(&entity.RefreshToken{}).Where("user_id = ? AND is_revoked = false", user.ID).Count(&activeTokensCount)
	if activeTokensCount != 0 {
		t.Errorf("expected 0 active tokens left after replay attack, got %d", activeTokensCount)
	}
}

func TestAuthService_RBAC_Query(t *testing.T) {
	db, priv, pub := setupTestDB(t)
	defer cleanTestCerts(priv)

	repo := repository.NewAuthRepository(db)
	txMgr := repository.NewTransactionManager(db)
	tokenMgr, _ := utils.NewTokenManager(priv, pub, 15)

	svc := NewAuthService(repo, txMgr, tokenMgr, nil, 7, 15)
	ctx := context.Background()

	// Buat dynamic permissions & roles
	adminRole, _ := svc.CreateRole(ctx, "admin", "Administrator")
	readPerm, _ := svc.CreatePermission(ctx, "wallet:read", "Read wallet metadata")
	writePerm, _ := svc.CreatePermission(ctx, "wallet:write", "Modify wallet balance")

	_ = svc.AssignPermissionToRole(ctx, adminRole.ID, readPerm.ID)
	_ = svc.AssignPermissionToRole(ctx, adminRole.ID, writePerm.ID)

	// Register user
	user, _ := svc.Register(ctx, "admin@platform.com", "adminpass")
	_ = svc.AssignRoleToUser(ctx, user.ID, adminRole.ID)

	// Verify JOIN queries GORM
	roleName, permissions, err := repo.GetUserRolesAndPermissions(ctx, user.ID)
	if err != nil {
		t.Fatalf("failed to query roles/permissions: %v", err)
	}

	// User memiliki roles: 'customer' (default saat register) dan 'admin' (yang ditambahkan)
	// Kita periksa apakah permissions 'wallet:read' dan 'wallet:write' berhasil terpetakan
	hasRead := false
	hasWrite := false
	for _, p := range permissions {
		if p == "wallet:read" {
			hasRead = true
		}
		if p == "wallet:write" {
			hasWrite = true
		}
	}

	if !hasRead || !hasWrite {
		t.Errorf("expected user to possess wallet permissions, got: %v (role: %s)", permissions, roleName)
	}
}
