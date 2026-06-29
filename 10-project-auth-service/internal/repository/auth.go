package repository

import (
	"context"
	"errors"
	"time"

	"github.com/timurdian/auth-service/internal/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrReplayAttackDetected = errors.New("refresh token replay attack detected")
	ErrTokenExpired         = errors.New("refresh token expired")
	ErrTokenRevoked         = errors.New("refresh token was revoked")
	ErrTokenNotFound        = errors.New("refresh token not found")
)

type AuthRepository interface {
	CreateUser(ctx context.Context, user *entity.User) error
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	GetUserByID(ctx context.Context, id uint) (*entity.User, error)
	UpdateUserVerification(ctx context.Context, userID uint, isVerified bool) error
	UpdateUserPassword(ctx context.Context, userID uint, hashedPass string) error

	CreateRefreshToken(ctx context.Context, rt *entity.RefreshToken) error
	GetRefreshToken(ctx context.Context, token string) (*entity.RefreshToken, error)
	RotateRefreshToken(ctx context.Context, oldTokenStr, newTokenStr string, expiry time.Time) (uint, error)
	RevokeAllUserTokens(ctx context.Context, userID uint) error
	RevokeToken(ctx context.Context, token string) error

	// RBAC Management
	CreateRole(ctx context.Context, role *entity.Role) error
	GetRoleByName(ctx context.Context, name string) (*entity.Role, error)
	CreatePermission(ctx context.Context, perm *entity.Permission) error
	GetPermissionByName(ctx context.Context, name string) (*entity.Permission, error)
	AssignRoleToUser(ctx context.Context, userID, roleID uint) error
	AssignPermissionToRole(ctx context.Context, roleID, permID uint) error
	GetUserRolesAndPermissions(ctx context.Context, userID uint) (string, []string, error)

	// Verification & Reset Password Tokens
	CreateVerificationToken(ctx context.Context, vt *entity.VerificationToken) error
	GetVerificationToken(ctx context.Context, token string) (*entity.VerificationToken, error)
	DeleteVerificationToken(ctx context.Context, token string) error

	CreateResetToken(ctx context.Context, rt *entity.ResetToken) error
	GetResetToken(ctx context.Context, token string) (*entity.ResetToken, error)
	DeleteResetToken(ctx context.Context, token string) error
}

type gormAuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &gormAuthRepository{db: db}
}

func (r *gormAuthRepository) CreateUser(ctx context.Context, user *entity.User) error {
	return GetDB(ctx, r.db).Create(user).Error
}

func (r *gormAuthRepository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	err := GetDB(ctx, r.db).Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *gormAuthRepository) GetUserByID(ctx context.Context, id uint) (*entity.User, error) {
	var user entity.User
	err := GetDB(ctx, r.db).First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *gormAuthRepository) UpdateUserVerification(ctx context.Context, userID uint, isVerified bool) error {
	return GetDB(ctx, r.db).Model(&entity.User{}).Where("id = ?", userID).Update("is_verified", isVerified).Error
}

func (r *gormAuthRepository) UpdateUserPassword(ctx context.Context, userID uint, hashedPass string) error {
	return GetDB(ctx, r.db).Model(&entity.User{}).Where("id = ?", userID).Update("password_hash", hashedPass).Error
}

func (r *gormAuthRepository) CreateRefreshToken(ctx context.Context, rt *entity.RefreshToken) error {
	return GetDB(ctx, r.db).Create(rt).Error
}

func (r *gormAuthRepository) GetRefreshToken(ctx context.Context, token string) (*entity.RefreshToken, error) {
	var rt entity.RefreshToken
	err := GetDB(ctx, r.db).Where("token = ?", token).First(&rt).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &rt, nil
}

func (r *gormAuthRepository) RevokeAllUserTokens(ctx context.Context, userID uint) error {
	return GetDB(ctx, r.db).Model(&entity.RefreshToken{}).Where("user_id = ? AND is_revoked = false", userID).Update("is_revoked", true).Error
}

func (r *gormAuthRepository) RevokeToken(ctx context.Context, token string) error {
	return GetDB(ctx, r.db).Model(&entity.RefreshToken{}).Where("token = ?", token).Update("is_revoked", true).Error
}

func (r *gormAuthRepository) RotateRefreshToken(ctx context.Context, oldTokenStr, newTokenStr string, expiry time.Time) (uint, error) {
	var userID uint
	// Kita jalankan pemutaran token di dalam transaksi database transaksional (jika belum ada transaksi terluar)
	txDB := GetDB(ctx, r.db)

	var oldToken entity.RefreshToken
	// Query GORM dengan SELECT ... FOR UPDATE untuk mengunci baris data
	err := txDB.Clauses(clause.Locking{Strength: "UPDATE"}).Where("token = ?", oldTokenStr).First(&oldToken).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, ErrTokenNotFound
		}
		return 0, err
	}

	// 1. Deteksi Replay Attack: Jika token lama sudah dicabut / di-revoke sebelumnya
	if oldToken.IsRevoked {
		// Batalkan secara massal seluruh refresh token aktif milik user tersebut demi keamanan
		_ = txDB.Model(&entity.RefreshToken{}).Where("user_id = ?", oldToken.UserID).Update("is_revoked", true)
		return oldToken.UserID, ErrReplayAttackDetected
	}

	// 2. Cek Expiry
	if oldToken.ExpiresAt.Before(time.Now()) {
		return 0, ErrTokenExpired
	}

	// 3. Lakukan rotasi: cabut token lama
	oldToken.IsRevoked = true
	if err := txDB.Save(&oldToken).Error; err != nil {
		return 0, err
	}

	// 4. Buat token refresh baru
	newToken := entity.RefreshToken{
		Token:       newTokenStr,
		UserID:      oldToken.UserID,
		ExpiresAt:   expiry,
		ParentToken: oldTokenStr,
	}
	if err := txDB.Create(&newToken).Error; err != nil {
		return 0, err
	}

	userID = oldToken.UserID
	return userID, nil
}

// RBAC Management
func (r *gormAuthRepository) CreateRole(ctx context.Context, role *entity.Role) error {
	return GetDB(ctx, r.db).Create(role).Error
}

func (r *gormAuthRepository) GetRoleByName(ctx context.Context, name string) (*entity.Role, error) {
	var role entity.Role
	err := GetDB(ctx, r.db).Where("name = ?", name).First(&role).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &role, nil
}

func (r *gormAuthRepository) CreatePermission(ctx context.Context, perm *entity.Permission) error {
	return GetDB(ctx, r.db).Create(perm).Error
}

func (r *gormAuthRepository) GetPermissionByName(ctx context.Context, name string) (*entity.Permission, error) {
	var perm entity.Permission
	err := GetDB(ctx, r.db).Where("name = ?", name).First(&perm).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &perm, nil
}

func (r *gormAuthRepository) AssignRoleToUser(ctx context.Context, userID, roleID uint) error {
	// Menambahkan baris ke join table user_roles
	return GetDB(ctx, r.db).Exec("INSERT INTO user_roles (user_id, role_id) VALUES (?, ?) ON CONFLICT DO NOTHING", userID, roleID).Error
}

func (r *gormAuthRepository) AssignPermissionToRole(ctx context.Context, roleID, permID uint) error {
	// Menambahkan baris ke join table role_permissions
	return GetDB(ctx, r.db).Exec("INSERT INTO role_permissions (role_id, permission_id) VALUES (?, ?) ON CONFLICT DO NOTHING", roleID, permID).Error
}

func (r *gormAuthRepository) GetUserRolesAndPermissions(ctx context.Context, userID uint) (string, []string, error) {
	var user entity.User
	// Load relasi many-to-many Roles dan sub-relasi many-to-many Permissions
	err := GetDB(ctx, r.db).Preload("Roles.Permissions").First(&user, userID).Error
	if err != nil {
		return "", nil, err
	}

	roleName := ""
	permissions := []string{}
	// Jika user memiliki roles, ambil role pertama (single role design untuk kemudahan, atau gabungkan)
	if len(user.Roles) > 0 {
		roleName = user.Roles[0].Name
		for _, r := range user.Roles {
			for _, p := range r.Permissions {
				permissions = append(permissions, p.Name)
			}
		}
	}

	// Bersihkan duplikasi permissions
	uniquePerms := []string{}
	seen := make(map[string]bool)
	for _, p := range permissions {
		if !seen[p] {
			seen[p] = true
			uniquePerms = append(uniquePerms, p)
		}
	}

	return roleName, uniquePerms, nil
}

// Verification & Reset Tokens
func (r *gormAuthRepository) CreateVerificationToken(ctx context.Context, vt *entity.VerificationToken) error {
	return GetDB(ctx, r.db).Create(vt).Error
}

func (r *gormAuthRepository) GetVerificationToken(ctx context.Context, token string) (*entity.VerificationToken, error) {
	var vt entity.VerificationToken
	err := GetDB(ctx, r.db).Where("token = ?", token).First(&vt).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &vt, nil
}

func (r *gormAuthRepository) DeleteVerificationToken(ctx context.Context, token string) error {
	return GetDB(ctx, r.db).Where("token = ?", token).Delete(&entity.VerificationToken{}).Error
}

func (r *gormAuthRepository) CreateResetToken(ctx context.Context, rt *entity.ResetToken) error {
	return GetDB(ctx, r.db).Create(rt).Error
}

func (r *gormAuthRepository) GetResetToken(ctx context.Context, token string) (*entity.ResetToken, error) {
	var rt entity.ResetToken
	err := GetDB(ctx, r.db).Where("token = ?", token).First(&rt).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &rt, nil
}

func (r *gormAuthRepository) DeleteResetToken(ctx context.Context, token string) error {
	return GetDB(ctx, r.db).Where("token = ?", token).Delete(&entity.ResetToken{}).Error
}
