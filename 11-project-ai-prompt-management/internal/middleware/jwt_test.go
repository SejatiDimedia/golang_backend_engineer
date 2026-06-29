package middleware

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func generateTestRSAKeys(t *testing.T) (string, *rsa.PrivateKey) {
	// Generate Private Key
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate private key: %v", err)
	}

	// Create temp dir
	tmpDir, err := ioutil.TempDir("", "jwt-test-keys")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	// Save Public Key to PEM
	pubASN1, err := x509.MarshalPKIXPublicKey(&privKey.PublicKey)
	if err != nil {
		t.Fatalf("failed to marshal public key: %v", err)
	}

	pubBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubASN1,
	}

	pubPath := filepath.Join(tmpDir, "public.key")
	pubFile, err := os.Create(pubPath)
	if err != nil {
		t.Fatalf("failed to create public.key file: %v", err)
	}
	defer pubFile.Close()

	if err := pem.Encode(pubFile, pubBlock); err != nil {
		t.Fatalf("failed to encode public key: %v", err)
	}

	return pubPath, privKey
}

func TestJWTMiddleware_OfflineVerification(t *testing.T) {
	gin.SetMode(gin.TestMode)

	pubPath, privKey := generateTestRSAKeys(t)
	defer os.RemoveAll(filepath.Dir(pubPath)) // cleanup

	mw, err := NewJWTMiddleware(pubPath)
	if err != nil {
		t.Fatalf("Failed to initialize JWTMiddleware: %v", err)
	}

	r := gin.New()
	r.Use(mw.AuthRequired())
	r.GET("/protected", func(c *gin.Context) {
		uid, _ := GetUserID(c)
		role, _ := GetUserRole(c)
		c.JSON(http.StatusOK, gin.H{"status": "success", "user_id": uid, "role": role})
	})

	// 1. Generate Valid RS256 JWT Token
	claims := &AccessTokenClaims{
		UserID:      99,
		Email:       "test@example.com",
		Role:        "admin",
		Permissions: []string{"read", "write"},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(10 * time.Minute)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	validTokenStr, err := token.SignedString(privKey)
	if err != nil {
		t.Fatalf("failed to sign token: %v", err)
	}

	// 2. Test Request Tanpa Token -> 401
	req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}

	// 3. Test Request Token Valid -> 200 & Claims Injection
	req, _ = http.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+validTokenStr)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d, body: %s", w.Code, w.Body.String())
	}

	// 4. Test Request Token Expired / Invalid Signing Method -> 401
	invalidClaims := &AccessTokenClaims{
		UserID: 99,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-10 * time.Minute)), // expired
		},
	}
	invalidToken := jwt.NewWithClaims(jwt.SigningMethodHS256, invalidClaims) // wrong signing method
	invalidTokenStr, _ := invalidToken.SignedString([]byte("badsecret"))

	req, _ = http.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+invalidTokenStr)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}
