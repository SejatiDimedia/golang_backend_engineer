package service

import (
	"context"
	"errors"

	"github.com/timurdian/booking-system/internal/entity"
	"github.com/timurdian/booking-system/internal/repository"
	"github.com/timurdian/booking-system/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmailConflict       = errors.New("email already registered")
	ErrInvalidCredentials  = errors.New("invalid email or password")
)

type UserService interface {
	Register(ctx context.Context, email, password, role string) (*entity.User, error)
	Login(ctx context.Context, email, password, jwtSecret string, expiryHours int) (string, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) Register(ctx context.Context, email, password, role string) (*entity.User, error) {
	if email == "" || password == "" {
		return nil, errors.New("email and password cannot be empty")
	}

	// 1. Cek duplikasi email
	existing, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrEmailConflict
	}

	// 2. Hash password menggunakan bcrypt
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	if role == "" {
		role = "customer"
	}

	user := &entity.User{
		Email:        email,
		PasswordHash: string(hashed),
		Role:         role,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) Login(ctx context.Context, email, password, jwtSecret string, expiryHours int) (string, error) {
	if email == "" || password == "" {
		return "", errors.New("email and password cannot be empty")
	}

	// 1. Cari user berdasarkan email
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", ErrInvalidCredentials
	}

	// 2. Verifikasi hash password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", ErrInvalidCredentials
	}

	// 3. Generate token JWT
	token, err := utils.GenerateToken(user.ID, user.Email, user.Role, jwtSecret, expiryHours)
	if err != nil {
		return "", err
	}

	return token, nil
}
