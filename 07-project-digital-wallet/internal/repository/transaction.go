package repository

import (
	"context"

	"github.com/timurdian/digital-wallet/internal/entity"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	Create(ctx context.Context, transaction *entity.Transaction) error
	GetByWalletID(ctx context.Context, walletID uint) ([]entity.Transaction, error)
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) Create(ctx context.Context, transaction *entity.Transaction) error {
	db := GetDBFromContext(ctx, r.db)
	return db.WithContext(ctx).Create(transaction).Error
}

func (r *transactionRepository) GetByWalletID(ctx context.Context, walletID uint) ([]entity.Transaction, error) {
	db := GetDBFromContext(ctx, r.db)
	var transactions []entity.Transaction
	
	// Cari transaksi yang melibatkan walletID baik sebagai source (debit) maupun destination (kredit)
	err := db.WithContext(ctx).
		Where("source_wallet_id = ? OR destination_wallet_id = ?", walletID, walletID).
		Order("id DESC").
		Find(&transactions).
		Error
	
	return transactions, err
}
