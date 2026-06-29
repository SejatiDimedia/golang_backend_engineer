package repository

import (
	"context"
	"errors"

	"github.com/timurdian/digital-wallet/internal/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type WalletRepository interface {
	Create(ctx context.Context, wallet *entity.Wallet) error
	GetByUserID(ctx context.Context, userID uint) (*entity.Wallet, error)
	GetByWalletNumber(ctx context.Context, number string) (*entity.Wallet, error)
	GetByWalletNumberForUpdate(ctx context.Context, number string) (*entity.Wallet, error)
	Update(ctx context.Context, wallet *entity.Wallet) error
}

type walletRepository struct {
	db *gorm.DB
}

func NewWalletRepository(db *gorm.DB) WalletRepository {
	return &walletRepository{db: db}
}

func (r *walletRepository) Create(ctx context.Context, wallet *entity.Wallet) error {
	db := GetDBFromContext(ctx, r.db)
	return db.WithContext(ctx).Create(wallet).Error
}

func (r *walletRepository) GetByUserID(ctx context.Context, userID uint) (*entity.Wallet, error) {
	db := GetDBFromContext(ctx, r.db)
	var wallet entity.Wallet
	err := db.WithContext(ctx).Where("user_id = ?", userID).First(&wallet).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &wallet, nil
}

func (r *walletRepository) GetByWalletNumber(ctx context.Context, number string) (*entity.Wallet, error) {
	db := GetDBFromContext(ctx, r.db)
	var wallet entity.Wallet
	err := db.WithContext(ctx).Where("wallet_number = ?", number).First(&wallet).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &wallet, nil
}

func (r *walletRepository) GetByWalletNumberForUpdate(ctx context.Context, number string) (*entity.Wallet, error) {
	db := GetDBFromContext(ctx, r.db)
	var wallet entity.Wallet
	err := db.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("wallet_number = ?", number).
		First(&wallet).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &wallet, nil
}

func (r *walletRepository) Update(ctx context.Context, wallet *entity.Wallet) error {
	db := GetDBFromContext(ctx, r.db)
	return db.WithContext(ctx).Save(wallet).Error
}
