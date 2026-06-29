package repository

import (
	"context"

	"gorm.io/gorm"
)

type txKey struct{}

// TransactionManager mengontrol inisiasi transaksi database relasional
type TransactionManager interface {
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type gormTxManager struct {
	db *gorm.DB
}

// NewGormTxManager membuat instance baru dari TransactionManager berbasis GORM
func NewGormTxManager(db *gorm.DB) TransactionManager {
	return &gormTxManager{db: db}
}

// WithTransaction membungkus fungsi callback di dalam transaksi database GORM
func (m *gormTxManager) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return m.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Pasang db transaction (GORM) ke context
		txCtx := context.WithValue(ctx, txKey{}, tx)
		return fn(txCtx)
	})
}

// GetDBFromContext mengambil instance database transaction jika ada, jika tidak, mengembalikan fallbackDB default
func GetDBFromContext(ctx context.Context, fallbackDB *gorm.DB) *gorm.DB {
	if tx, ok := ctx.Value(txKey{}).(*gorm.DB); ok {
		return tx
	}
	return fallbackDB
}
