package service

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/timurdian/digital-wallet/internal/entity"
)

type mockTxManager struct{}

func (m *mockTxManager) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

type mockWalletRepository struct {
	wallets       map[uint]*entity.Wallet
	walletsByNum  map[string]*entity.Wallet
	mu            sync.Mutex
}

func (m *mockWalletRepository) Create(ctx context.Context, wallet *entity.Wallet) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.wallets[wallet.UserID] = wallet
	m.walletsByNum[wallet.WalletNumber] = wallet
	return nil
}

func (m *mockWalletRepository) GetByUserID(ctx context.Context, userID uint) (*entity.Wallet, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	w, exists := m.wallets[userID]
	if !exists {
		return nil, nil
	}
	return w, nil
}

func (m *mockWalletRepository) GetByWalletNumber(ctx context.Context, number string) (*entity.Wallet, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	w, exists := m.walletsByNum[number]
	if !exists {
		return nil, nil
	}
	return w, nil
}

func (m *mockWalletRepository) GetByWalletNumberForUpdate(ctx context.Context, number string) (*entity.Wallet, error) {
	return m.GetByWalletNumber(ctx, number)
}

func (m *mockWalletRepository) Update(ctx context.Context, wallet *entity.Wallet) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.wallets[wallet.UserID] = wallet
	m.walletsByNum[wallet.WalletNumber] = wallet
	return nil
}

type mockTransactionRepository struct {
	transactions []entity.Transaction
	mu           sync.Mutex
}

func (m *mockTransactionRepository) Create(ctx context.Context, transaction *entity.Transaction) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	transaction.ID = uint(len(m.transactions) + 1)
	m.transactions = append(m.transactions, *transaction)
	return nil
}

func (m *mockTransactionRepository) GetByWalletID(ctx context.Context, walletID uint) ([]entity.Transaction, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var result []entity.Transaction
	for _, t := range m.transactions {
		if (t.SourceWalletID != nil && *t.SourceWalletID == walletID) || (t.DestinationWalletID != nil && *t.DestinationWalletID == walletID) {
			result = append(result, t)
		}
	}
	return result, nil
}

type mockLocalLockManager struct {
	locks map[string]string
	mu    sync.Mutex
}

func (m *mockLocalLockManager) AcquireLock(ctx context.Context, key string, ttl time.Duration, retries int, backoff time.Duration) (string, error) {
	for i := 0; i <= retries; i++ {
		m.mu.Lock()
		_, exists := m.locks[key]
		if !exists {
			token := "token_" + key
			m.locks[key] = token
			m.mu.Unlock()
			return token, nil
		}
		m.mu.Unlock()
		
		if i < retries {
			time.Sleep(backoff)
		}
	}
	return "", errors.New("lock acquire failed")
}

func (m *mockLocalLockManager) ReleaseLock(ctx context.Context, key string, token string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	val, exists := m.locks[key]
	if exists && val == token {
		delete(m.locks, key)
		return nil
	}
	return errors.New("release lock failed")
}

func TestWalletService_TopUpAndWithdraw(t *testing.T) {
	txM := &mockTxManager{}
	walletRepo := &mockWalletRepository{
		wallets:      make(map[uint]*entity.Wallet),
		walletsByNum: make(map[string]*entity.Wallet),
	}
	transRepo := &mockTransactionRepository{}
	lockM := &mockLocalLockManager{locks: make(map[string]string)}

	// Setup Wallet user 1
	walletRepo.Create(context.Background(), &entity.Wallet{ID: 1, UserID: 1, WalletNumber: "W-1001", Balance: 100.0})

	svc := NewWalletService(txM, walletRepo, transRepo, nil, lockM)

	// Test Topup +50
	tx, err := svc.TopUp(context.Background(), 1, 50.0, "topup bonus")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if tx.Type != "top-up" || tx.Amount != 50.0 {
		t.Errorf("invalid transaction log: %+v", tx)
	}

	w, _ := walletRepo.GetByUserID(context.Background(), 1)
	if w.Balance != 150.0 {
		t.Errorf("expected balance 150, got %.2f", w.Balance)
	}

	// Test Withdraw -30
	tx2, err := svc.Withdraw(context.Background(), 1, 30.0, "cashout")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if tx2.Type != "withdraw" || tx2.Amount != 30.0 {
		t.Errorf("invalid transaction log: %+v", tx2)
	}

	w, _ = walletRepo.GetByUserID(context.Background(), 1)
	if w.Balance != 120.0 {
		t.Errorf("expected balance 120, got %.2f", w.Balance)
	}

	// Test Withdraw over balance
	_, err = svc.Withdraw(context.Background(), 1, 200.0, "cashout large")
	if !errors.Is(err, ErrInsufficientBalance) {
		t.Errorf("expected ErrInsufficientBalance, got %v", err)
	}
}

func TestWalletService_ConcurrentTransferSafety(t *testing.T) {
	txM := &mockTxManager{}
	walletRepo := &mockWalletRepository{
		wallets:      make(map[uint]*entity.Wallet),
		walletsByNum: make(map[string]*entity.Wallet),
	}
	transRepo := &mockTransactionRepository{}
	lockM := &mockLocalLockManager{locks: make(map[string]string)}

	// User 1 (Sender): balance Rp 100.000
	walletRepo.Create(context.Background(), &entity.Wallet{ID: 1, UserID: 1, WalletNumber: "W-1001", Balance: 100000.0})
	// User 2 (Receiver): balance Rp 0
	walletRepo.Create(context.Background(), &entity.Wallet{ID: 2, UserID: 2, WalletNumber: "W-1002", Balance: 0.0})

	svc := NewWalletService(txM, walletRepo, transRepo, nil, lockM)

	// Simulasi transfer konkuren: 10 transfer paralel masing-masing Rp 10.000
	// Nilai transfer total = Rp 100.000 (saldo habis pas)
	// Kita jalankan Goroutines secara paralel
	var wg sync.WaitGroup
	workers := 10
	transferAmount := 10000.0

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := svc.Transfer(context.Background(), 1, "W-1002", transferAmount, "transfer thr")
			if err != nil {
				// Log error testing jika transfer gagal
				t.Logf("Concurrent transfer error: %v", err)
			}
		}()
	}

	wg.Wait()

	// Cek saldo akhir
	w1, _ := walletRepo.GetByUserID(context.Background(), 1)
	w2, _ := walletRepo.GetByUserID(context.Background(), 2)

	if w1.Balance != 0.0 {
		t.Errorf("expected sender balance to be exactly 0, got %.2f", w1.Balance)
	}

	if w2.Balance != 100000.0 {
		t.Errorf("expected receiver balance to be exactly 100000, got %.2f", w2.Balance)
	}

	// Cek record log ledger transaksi
	txs1, _ := transRepo.GetByWalletID(context.Background(), 1)
	if len(txs1) != 10 {
		t.Errorf("expected exactly 10 ledger transactions, got %d", len(txs1))
	}
}
