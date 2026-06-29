package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/timurdian/digital-wallet/internal/entity"
	"github.com/timurdian/digital-wallet/internal/repository"
	"github.com/timurdian/digital-wallet/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmailConflict      = errors.New("email already registered")
	ErrInvalidCredentials = errors.New("invalid email or password")
)

type UserService interface {
	Register(ctx context.Context, email, password string) (*entity.User, error)
	Login(ctx context.Context, email, password, jwtSecret string, expiryHours int) (string, error)
}

type userService struct {
	txManager  repository.TransactionManager
	userRepo   repository.UserRepository
	walletRepo repository.WalletRepository
}

func NewUserService(
	txManager repository.TransactionManager,
	userRepo repository.UserRepository,
	walletRepo repository.WalletRepository,
) UserService {
	return &userService{
		txManager:  txManager,
		userRepo:   userRepo,
		walletRepo: walletRepo,
	}
}

func (s *userService) Register(ctx context.Context, email, password string) (*entity.User, error) {
	if email == "" || password == "" {
		return nil, errors.New("email and password cannot be empty")
	}

	// 1. Cek duplikasi email
	existing, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrEmailConflict
	}

	// 2. Hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	var user *entity.User

	// 3. Simpan user & buat wallet dalam transaksi database tunggal
	err = s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		user = &entity.User{
			Email:        email,
			PasswordHash: string(hashed),
		}

		if err := s.userRepo.Create(txCtx, user); err != nil {
			return err
		}

		// Generate WalletNumber: format "W-1000" + userID
		walletNumber := fmt.Sprintf("W-1000%d", user.ID)

		wallet := &entity.Wallet{
			UserID:       user.ID,
			WalletNumber: walletNumber,
			Balance:      0.00,
		}

		if err := s.walletRepo.Create(txCtx, wallet); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) Login(ctx context.Context, email, password, jwtSecret string, expiryHours int) (string, error) {
	if email == "" || password == "" {
		return "", errors.New("email and password cannot be empty")
	}

	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", ErrInvalidCredentials
	}

	token, err := utils.GenerateToken(user.ID, user.Email, jwtSecret, expiryHours)
	if err != nil {
		return "", err
	}

	return token, nil
}
