package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEnsureRSAKeys(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "certs-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	privPath := filepath.Join(tempDir, "private.key")
	pubPath := filepath.Join(tempDir, "public.key")

	// 1. Jalankan pertama kali -> harus membuat file baru
	err = EnsureRSAKeys(privPath, pubPath)
	if err != nil {
		t.Fatalf("EnsureRSAKeys failed: %v", err)
	}

	if _, err := os.Stat(privPath); os.IsNotExist(err) {
		t.Errorf("expected private key file to exist")
	}
	if _, err := os.Stat(pubPath); os.IsNotExist(err) {
		t.Errorf("expected public key file to exist")
	}

	// 2. Jalankan kedua kali -> harus lewati pembuatan kunci baru (tidak menimpa kunci yang sudah ada)
	privInfo, _ := os.Stat(privPath)
	privModTime := privInfo.ModTime()

	err = EnsureRSAKeys(privPath, pubPath)
	if err != nil {
		t.Fatalf("EnsureRSAKeys second call failed: %v", err)
	}

	privInfo2, _ := os.Stat(privPath)
	if !privInfo2.ModTime().Equal(privModTime) {
		t.Errorf("expected private key to not be overwritten")
	}
}

func TestTokenManager_GenerateAndValidate(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "jwt-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	privPath := filepath.Join(tempDir, "private.key")
	pubPath := filepath.Join(tempDir, "public.key")

	_ = EnsureRSAKeys(privPath, pubPath)

	tokenMgr, err := NewTokenManager(privPath, pubPath, 15)
	if err != nil {
		t.Fatalf("failed to init TokenManager: %v", err)
	}

	perms := []string{"wallet:read", "wallet:write"}
	tokenStr, err := tokenMgr.GenerateAccessToken(101, "user@test.com", "admin", perms)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	claims, err := tokenMgr.ValidateAccessToken(tokenStr)
	if err != nil {
		t.Fatalf("failed to validate token: %v", err)
	}

	if claims.UserID != 101 {
		t.Errorf("expected user ID 101, got %d", claims.UserID)
	}
	if claims.Role != "admin" {
		t.Errorf("expected role 'admin', got '%s'", claims.Role)
	}
	if len(claims.Permissions) != 2 || claims.Permissions[0] != "wallet:read" {
		t.Errorf("unexpected claims permissions: %v", claims.Permissions)
	}
}
