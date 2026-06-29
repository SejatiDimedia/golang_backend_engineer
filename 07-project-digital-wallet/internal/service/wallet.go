package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/timurdian/digital-wallet/internal/entity"
	"github.com/timurdian/digital-wallet/internal/repository"
	"github.com/timurdian/digital-wallet/internal/utils"
)

var (
	ErrInsufficientBalance = errors.New("insufficient wallet balance")
	ErrWalletNotFound       = errors.New("recipient wallet not found")
	ErrInvalidTransfer      = errors.New("cannot transfer to your own wallet")
	ErrInvalidAmount        = errors.New("amount must be greater than zero")
)

type WalletService interface {
	GetBalance(ctx context.Context, userID uint) (float64, error)
	TopUp(ctx context.Context, userID uint, amount float64, description string) (*entity.Transaction, error)
	Withdraw(ctx context.Context, userID uint, amount float64, description string) (*entity.Transaction, error)
	Transfer(ctx context.Context, userID uint, destWalletNumber string, amount float64, description string) (*entity.Transaction, error)
	GetTransactions(ctx context.Context, userID uint) ([]entity.Transaction, error)
}

type walletService struct {
	txManager   repository.TransactionManager
	walletRepo  repository.WalletRepository
	transRepo   repository.TransactionRepository
	rdb         *redis.Client
	lockManager utils.LockManager
}

func NewWalletService(
	txManager repository.TransactionManager,
	walletRepo repository.WalletRepository,
	transRepo repository.TransactionRepository,
	rdb *redis.Client,
	lockManager utils.LockManager,
) WalletService {
	return &walletService{
		txManager:   txManager,
		walletRepo:  walletRepo,
		transRepo:   transRepo,
		rdb:         rdb,
		lockManager: lockManager,
	}
}

func (s *walletService) GetBalance(ctx context.Context, userID uint) (float64, error) {
	wallet, err := s.walletRepo.GetByUserID(ctx, userID)
	if err != nil {
		return 0, err
	}
	if wallet == nil {
		return 0, ErrWalletNotFound
	}

	cacheKey := fmt.Sprintf("wallet:balance:%d", wallet.ID)

	// 1. Cek cache Redis
	if s.rdb != nil {
		val, err := s.rdb.Get(ctx, cacheKey).Result()
		if err == nil {
			if balance, err := strconv.ParseFloat(val, 64); err == nil {
				return balance, nil
			}
		}
	}

	// 2. Cache miss -> Ambil dari Database
	// Kita peroleh balance aktual
	balance := wallet.Balance

	// 3. Simpan ke Redis cache (TTL 10 Menit)
	if s.rdb != nil {
		s.rdb.Set(ctx, cacheKey, fmt.Sprintf("%.2f", balance), 10*time.Minute)
	}

	return balance, nil
}

func (s *walletService) TopUp(ctx context.Context, userID uint, amount float64, description string) (*entity.Transaction, error) {
	if amount <= 0 {
		return nil, ErrInvalidAmount
	}

	wallet, err := s.walletRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if wallet == nil {
		return nil, ErrWalletNotFound
	}

	var transaction *entity.Transaction

	err = s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		// Dapatkan data ter-update
		w, err := s.walletRepo.GetByUserID(txCtx, userID)
		if err != nil {
			return err
		}

		w.Balance += amount
		if err := s.walletRepo.Update(txCtx, w); err != nil {
			return err
		}

		transaction = &entity.Transaction{
			DestinationWalletID: &w.ID,
			Amount:              amount,
			Type:                "top-up",
			Description:         description,
		}

		return s.transRepo.Create(txCtx, transaction)
	})

	if err != nil {
		return nil, err
	}

	// Invalidate Cache
	if s.rdb != nil {
		cacheKey := fmt.Sprintf("wallet:balance:%d", wallet.ID)
		s.rdb.Del(ctx, cacheKey)
	}

	return transaction, nil
}

func (s *walletService) Withdraw(ctx context.Context, userID uint, amount float64, description string) (*entity.Transaction, error) {
	if amount <= 0 {
		return nil, ErrInvalidAmount
	}

	wallet, err := s.walletRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if wallet == nil {
		return nil, ErrWalletNotFound
	}

	var transaction *entity.Transaction

	err = s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		w, err := s.walletRepo.GetByUserID(txCtx, userID)
		if err != nil {
			return err
		}

		if w.Balance < amount {
			return ErrInsufficientBalance
		}

		w.Balance -= amount
		if err := s.walletRepo.Update(txCtx, w); err != nil {
			return err
		}

		transaction = &entity.Transaction{
			SourceWalletID: &w.ID,
			Amount:         amount,
			Type:           "withdraw",
			Description:    description,
		}

		return s.transRepo.Create(txCtx, transaction)
	})

	if err != nil {
		return nil, err
	}

	// Invalidate Cache
	if s.rdb != nil {
		cacheKey := fmt.Sprintf("wallet:balance:%d", wallet.ID)
		s.rdb.Del(ctx, cacheKey)
	}

	return transaction, nil
}

func (s *walletService) Transfer(ctx context.Context, userID uint, destWalletNumber string, amount float64, description string) (*entity.Transaction, error) {
	if amount <= 0 {
		return nil, ErrInvalidAmount
	}

	// 1. Dapatkan info wallet pengirim dan penerima
	senderWallet, err := s.walletRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if senderWallet == nil {
		return nil, ErrWalletNotFound
	}

	receiverWallet, err := s.walletRepo.GetByWalletNumber(ctx, destWalletNumber)
	if err != nil {
		return nil, err
	}
	if receiverWallet == nil {
		return nil, ErrWalletNotFound
	}

	if senderWallet.ID == receiverWallet.ID {
		return nil, ErrInvalidTransfer
	}

	// 2. Tentukan urutan locking terdistribusi untuk menghindari deadlock
	firstKey := fmt.Sprintf("lock:wallet:%d", senderWallet.ID)
	secondKey := fmt.Sprintf("lock:wallet:%d", receiverWallet.ID)
	firstID := senderWallet.ID
	secondID := receiverWallet.ID

	if senderWallet.ID > receiverWallet.ID {
		firstKey, secondKey = secondKey, firstKey
		firstID, secondID = secondID, firstID
	}

	// 3. Acquire Distributed Locks di Redis
	log.Printf("[Lock] Acquiring locks for wallets %d and %d", firstID, secondID)
	token1, err := s.lockManager.AcquireLock(ctx, firstKey, 5*time.Second, 5, 100*time.Millisecond)
	if err != nil {
		return nil, err
	}
	defer s.lockManager.ReleaseLock(ctx, firstKey, token1)

	token2, err := s.lockManager.AcquireLock(ctx, secondKey, 5*time.Second, 5, 100*time.Millisecond)
	if err != nil {
		return nil, err
	}
	defer s.lockManager.ReleaseLock(ctx, secondKey, token2)

	var transaction *entity.Transaction

	// 4. Jalankan transaksi database relasional atomic
	err = s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		// Dapatkan baris terupdate dari database
		sender, err := s.walletRepo.GetByUserID(txCtx, userID)
		if err != nil {
			return err
		}

		receiver, err := s.walletRepo.GetByWalletNumber(txCtx, destWalletNumber)
		if err != nil {
			return err
		}

		if sender.Balance < amount {
			return ErrInsufficientBalance
		}

		// Debit-Kredit Balanced Ledger Update
		sender.Balance -= amount
		receiver.Balance += amount

		if err := s.walletRepo.Update(txCtx, sender); err != nil {
			return err
		}
		if err := s.walletRepo.Update(txCtx, receiver); err != nil {
			return err
		}

		transaction = &entity.Transaction{
			SourceWalletID:      &sender.ID,
			DestinationWalletID: &receiver.ID,
			Amount:              amount,
			Type:                "transfer",
			Description:         description,
		}

		return s.transRepo.Create(txCtx, transaction)
	})

	if err != nil {
		return nil, err
	}

	// 5. Invalidate Caches
	if s.rdb != nil {
		s.rdb.Del(ctx, fmt.Sprintf("wallet:balance:%d", senderWallet.ID))
		s.rdb.Del(ctx, fmt.Sprintf("wallet:balance:%d", receiverWallet.ID))
	}

	return transaction, nil
}

func (s *walletService) GetTransactions(ctx context.Context, userID uint) ([]entity.Transaction, error) {
	wallet, err := s.walletRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if wallet == nil {
		return nil, ErrWalletNotFound
	}

	return s.transRepo.GetByWalletID(ctx, wallet.ID)
}
