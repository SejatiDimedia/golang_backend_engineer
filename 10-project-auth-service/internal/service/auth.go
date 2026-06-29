package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/timurdian/auth-service/internal/entity"
	"github.com/timurdian/auth-service/internal/repository"
	"github.com/timurdian/auth-service/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserExists           = errors.New("user with this email already exists")
	ErrInvalidCredentials    = errors.New("invalid email or password")
	ErrEmailNotVerified     = errors.New("email is not verified yet")
	ErrInvalidOrExpiredToken = errors.New("invalid or expired token")
)

type AuthService interface {
	Register(ctx context.Context, email, password string) (*entity.User, error)
	Login(ctx context.Context, email, password string) (string, string, error)
	Refresh(ctx context.Context, refreshTokenStr string) (string, string, error)
	Logout(ctx context.Context, refreshTokenStr string) error

	VerifyEmail(ctx context.Context, token string) error
	ForgotPassword(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, token, newPassword string) error

	Introspect(ctx context.Context, tokenStr string) (bool, *utils.AccessTokenClaims, error)

	// RBAC dynamic setup (for setup script / admin endpoint)
	CreateRole(ctx context.Context, name, desc string) (*entity.Role, error)
	CreatePermission(ctx context.Context, name, desc string) (*entity.Permission, error)
	AssignRoleToUser(ctx context.Context, userID, roleID uint) error
	AssignPermissionToRole(ctx context.Context, roleID, permID uint) error
}

type authService struct {
	authRepo           repository.AuthRepository
	txMgr              repository.TransactionManager
	tokenMgr           utils.TokenManager
	rdb                *redis.Client
	refreshTokenDays   int
	accessTokenMinutes int
}

func NewAuthService(
	authRepo repository.AuthRepository,
	txMgr repository.TransactionManager,
	tokenMgr utils.TokenManager,
	rdb *redis.Client,
	refreshTokenDays int,
	accessTokenMinutes int,
) AuthService {
	return &authService{
		authRepo:           authRepo,
		txMgr:              txMgr,
		tokenMgr:           tokenMgr,
		rdb:                rdb,
		refreshTokenDays:   refreshTokenDays,
		accessTokenMinutes: accessTokenMinutes,
	}
}

func (s *authService) getRepo() repository.AuthRepository {
	return s.authRepo
}

func generateSecureToken() string {
	b := make([]byte, 32)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func (s *authService) Register(ctx context.Context, email, password string) (*entity.User, error) {
	var user *entity.User

	err := s.txMgr.WithTransaction(ctx, func(txCtx context.Context) error {
		repo := s.getRepo()
		existing, err := repo.GetUserByEmail(txCtx, email)
		if err != nil {
			return err
		}
		if existing != nil {
			return ErrUserExists
		}

		hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		user = &entity.User{
			Email:        email,
			PasswordHash: string(hashedPass),
			IsVerified:   false,
		}

		if err := repo.CreateUser(txCtx, user); err != nil {
			return err
		}

		// 1. Dapatkan / Buat default role 'customer'
		role, err := repo.GetRoleByName(txCtx, "customer")
		if err != nil {
			return err
		}
		if role == nil {
			role = &entity.Role{
				Name:        "customer",
				Description: "Default platform user role",
			}
			if err := repo.CreateRole(txCtx, role); err != nil {
				return err
			}
		}

		// Hubungkan role ke user
		if err := repo.AssignRoleToUser(txCtx, user.ID, role.ID); err != nil {
			return err
		}

		// 2. Generate email verification token
		tokenStr := generateSecureToken()
		vt := &entity.VerificationToken{
			UserID:    user.ID,
			Token:     tokenStr,
			ExpiresAt: time.Now().Add(24 * time.Hour), // Berlaku 24 jam
		}

		if err := repo.CreateVerificationToken(txCtx, vt); err != nil {
			return err
		}

		// Mock email sending
		log.Printf("[EMAIL OUT] To: %s | Verification link: http://localhost:8081/auth/verify-email?token=%s", email, tokenStr)
		return nil
	})

	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *authService) Login(ctx context.Context, email, password string) (string, string, error) {
	repo := s.getRepo()
	user, err := repo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", "", err
	}
	if user == nil {
		return "", "", ErrInvalidCredentials
	}

	// Verify Password Hash
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", "", ErrInvalidCredentials
	}

	// Verify Email Activation
	if !user.IsVerified {
		return "", "", ErrEmailNotVerified
	}

	// Ambil role dan permissions
	roleName, permissions, err := repo.GetUserRolesAndPermissions(ctx, user.ID)
	if err != nil {
		return "", "", err
	}

	// Generate Access Token (RS256 JWT)
	accessToken, err := s.tokenMgr.GenerateAccessToken(user.ID, user.Email, roleName, permissions)
	if err != nil {
		return "", "", err
	}

	// Generate Refresh Token
	refreshTokenStr := generateSecureToken()
	rt := &entity.RefreshToken{
		Token:     refreshTokenStr,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(time.Duration(s.refreshTokenDays) * 24 * time.Hour),
		IsRevoked: false,
	}

	if err := repo.CreateRefreshToken(ctx, rt); err != nil {
		return "", "", err
	}

	return accessToken, refreshTokenStr, nil
}

func (s *authService) Refresh(ctx context.Context, refreshTokenStr string) (string, string, error) {
	var accessToken, newRefreshTokenStr string
	var isReplayAttack bool
	var replayUserID uint

	// RTR pemutaran token dilakukan di dalam database transaction
	err := s.txMgr.WithTransaction(ctx, func(txCtx context.Context) error {
		repo := s.getRepo()

		newRefreshUUID := generateSecureToken()
		expiry := time.Now().Add(time.Duration(s.refreshTokenDays) * 24 * time.Hour)

		// Lakukan rotasi dengan lock baris di Postgres
		userID, err := repo.RotateRefreshToken(txCtx, refreshTokenStr, newRefreshUUID, expiry)
		if err != nil {
			if err == repository.ErrReplayAttackDetected {
				isReplayAttack = true
				replayUserID = userID
				return nil // Kembalikan nil agar pembatalan massal ter-commit di DB
			}
			return err
		}

		user, err := repo.GetUserByID(txCtx, userID)
		if err != nil {
			return err
		}

		// Ambil role & permissions
		roleName, permissions, err := repo.GetUserRolesAndPermissions(txCtx, user.ID)
		if err != nil {
			return err
		}

		// Generate access token baru
		accessToken, err = s.tokenMgr.GenerateAccessToken(user.ID, user.Email, roleName, permissions)
		if err != nil {
			return err
		}

		newRefreshTokenStr = newRefreshUUID
		return nil
	})

	if err != nil {
		return "", "", err
	}

	if isReplayAttack {
		log.Printf("[SECURITY] Replay attack detected for user ID %d! Revoked all sessions", replayUserID)
		return "", "", ErrInvalidOrExpiredToken
	}

	return accessToken, newRefreshTokenStr, nil
}

func (s *authService) Logout(ctx context.Context, refreshTokenStr string) error {
	return s.getRepo().RevokeToken(ctx, refreshTokenStr)
}

func (s *authService) VerifyEmail(ctx context.Context, token string) error {
	repo := s.getRepo()
	vt, err := repo.GetVerificationToken(ctx, token)
	if err != nil {
		return err
	}
	if vt == nil || vt.ExpiresAt.Before(time.Now()) {
		return ErrInvalidOrExpiredToken
	}

	err = s.txMgr.WithTransaction(ctx, func(txCtx context.Context) error {
		txRepo := s.getRepo()
		// 1. Set verified
		if err := txRepo.UpdateUserVerification(txCtx, vt.UserID, true); err != nil {
			return err
		}
		// 2. Hapus token verifikasi agar tidak bisa dipakai lagi
		return txRepo.DeleteVerificationToken(txCtx, token)
	})

	return err
}

func (s *authService) ForgotPassword(ctx context.Context, email string) error {
	repo := s.getRepo()
	user, err := repo.GetUserByEmail(ctx, email)
	if err != nil {
		return err
	}
	if user == nil {
		// return nil agar hacker tidak bisa brute-force mencari email terdaftar (security best practice)
		return nil
	}

	tokenStr := generateSecureToken()
	rt := &entity.ResetToken{
		UserID:    user.ID,
		Token:     tokenStr,
		ExpiresAt: time.Now().Add(1 * time.Hour), // Berlaku 1 jam
	}

	if err := repo.CreateResetToken(ctx, rt); err != nil {
		return err
	}

	// Mock email log
	log.Printf("[EMAIL OUT] To: %s | Password reset link: http://localhost:8081/auth/reset-password?token=%s", email, tokenStr)
	return nil
}

func (s *authService) ResetPassword(ctx context.Context, token, newPassword string) error {
	repo := s.getRepo()
	rt, err := repo.GetResetToken(ctx, token)
	if err != nil {
		return err
	}
	if rt == nil || rt.ExpiresAt.Before(time.Now()) {
		return ErrInvalidOrExpiredToken
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	err = s.txMgr.WithTransaction(ctx, func(txCtx context.Context) error {
		txRepo := s.getRepo()
		// 1. Update password
		if err := txRepo.UpdateUserPassword(txCtx, rt.UserID, string(hashedPass)); err != nil {
			return err
		}
		// 2. Hapus seluruh refresh token aktif demi keamanan (force logout other sessions setelah ganti password)
		_ = txRepo.RevokeAllUserTokens(txCtx, rt.UserID)
		// 3. Hapus reset token
		return txRepo.DeleteResetToken(txCtx, token)
	})

	return err
}

func (s *authService) Introspect(ctx context.Context, tokenStr string) (bool, *utils.AccessTokenClaims, error) {
	// Cek caching di Redis untuk menghindari overload database relasional
	cacheKey := fmt.Sprintf("jwt:introspect:%s", tokenStr)
	if s.rdb != nil {
		cachedRole, err := s.rdb.Get(ctx, cacheKey).Result()
		if err == nil && cachedRole == "revoked" {
			return false, nil, nil
		}
	}

	// Validasi JWT token secara offline menggunakan public key
	claims, err := s.tokenMgr.ValidateAccessToken(tokenStr)
	if err != nil {
		return false, nil, nil
	}

	// Simpan ke Redis cache untuk performa tinggi microservice introspect
	if s.rdb != nil {
		remainingTime := time.Until(claims.ExpiresAt.Time)
		if remainingTime > 0 {
			_ = s.rdb.Set(ctx, cacheKey, claims.Role, remainingTime).Err()
		}
	}

	return true, claims, nil
}

// Admin RBAC Management
func (s *authService) CreateRole(ctx context.Context, name, desc string) (*entity.Role, error) {
	role := &entity.Role{Name: name, Description: desc}
	err := s.getRepo().CreateRole(ctx, role)
	if err != nil {
		return nil, err
	}
	return role, nil
}

func (s *authService) CreatePermission(ctx context.Context, name, desc string) (*entity.Permission, error) {
	perm := &entity.Permission{Name: name, Description: desc}
	err := s.getRepo().CreatePermission(ctx, perm)
	if err != nil {
		return nil, err
	}
	return perm, nil
}

func (s *authService) AssignRoleToUser(ctx context.Context, userID, roleID uint) error {
	return s.getRepo().AssignRoleToUser(ctx, userID, roleID)
}
func (s *authService) AssignPermissionToRole(ctx context.Context, roleID, permID uint) error {
	return s.getRepo().AssignPermissionToRole(ctx, roleID, permID)
}
